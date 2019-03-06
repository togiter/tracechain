package fabricSetup

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"strings"
)

type FabricSetup struct {
	ConnectionProfile string
	OrgID             string
	OrdererID         string
	ChannelID         string
	ChaincodeID       string
	ChaincodeVersion  string
	ChannelConfig     string
	ChaincodeGoPath   string
	ChaincodePath     string
	OrgAdmin          string
	OrgName           string
	OrgPeer0          string
	UserName          string
	initialized       bool
	chCli             *channel.Client
	rmCli             *resmgmt.Client //资源管理客户端(相当于管理员admin)，用以创建或更新通道
	sdk               *fabsdk.FabricSDK
	chEvent           *event.Client
}

//由配置文件connetionProfile配置sdk，初始化client，chain，event hub
func (setup *FabricSetup) Initialize() error {

	if setup.initialized {
		return errors.New("sdk已经初始化...")
	}
	sdk, err := fabsdk.New(config.FromFile(setup.ConnectionProfile))
	if err != nil {
		fmt.Println("sdk初始化失败:%v", err)
		return errors.WithMessage(err, "sdk初始化失败")
	}
	setup.sdk = sdk
	fmt.Println("sdk成功初始化...")

	//返回资源管理上下文 contextApi.ClientProvider
	rmCtx := setup.sdk.Context(fabsdk.WithUser(setup.OrgAdmin), fabsdk.WithOrg(setup.OrgName))
	if rmCtx == nil {
		fmt.Println("根据指定的组织名称orgName与管理员OrgAdmin创建资源管理客户端Context失败")
		return errors.New("根据指定的组织名称orgName与管理员OrgAdmin创建资源管理客户端Context失败")
	}
	//创建资源管理器管理通道的创建或更新
	rmCli, err := resmgmt.New(rmCtx)
	if err != nil {
		fmt.Println("资源管理器创建失败:%v", err)
		return errors.WithMessage(err, "资源管理器创建失败！")
	}
	setup.rmCli = rmCli
	fmt.Println("资源管理器resource management创建成功！")

	//msp client允许我们从用户的identify检索用户信息
	mspCli, err := mspclient.New(sdk.Context(), mspclient.WithOrg(setup.OrgName))
	if err != nil {
		fmt.Println("failed to create msp client")
		return errors.WithMessage(err, "failed to create msp client")
	}

	//获取用户身份admin
	adminID, err := mspCli.GetSigningIdentity(setup.OrgAdmin)
	if err != nil {
		fmt.Println("failed to query channel")
		return errors.WithMessage(err, "failed to query channel")
	}
	channelHasInstalled := false
	//查询已经存在的channel
	channelResp, err := setup.rmCli.QueryChannels(resmgmt.WithTargetEndpoints(setup.OrgPeer0))
	if channelResp != nil {
		for _, channel := range channelResp.Channels {
			if strings.EqualFold(setup.ChannelID, channel.ChannelId) {
				channelHasInstalled = true
			}
		}
	}
	fmt.Println("channelHasInstalled:", channelHasInstalled)
	if !channelHasInstalled {
		//创建通道
		chReq := resmgmt.SaveChannelRequest{ChannelID: setup.ChannelID, ChannelConfigPath: setup.ChannelConfig, SigningIdentities: []msp.SigningIdentity{adminID}}
		txID, err := setup.rmCli.SaveChannel(chReq, resmgmt.WithOrdererEndpoint(setup.OrdererID))
		if err != nil || txID.TransactionID == "" {
			fmt.Println("faild to save channel:%v", err)
			return errors.WithMessage(err, "faild to save channel")
		}
		fmt.Println("Channel created")
		//加入通道
		if err = setup.rmCli.JoinChannel(setup.ChannelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithOrdererEndpoint(setup.OrdererID)); err != nil {
			fmt.Println("faild to join channel:%v", err)
			return errors.WithMessage(err, "failed to join channel")
		}
		fmt.Println("channel joined successed")
	} else {
		fmt.Println("channel already exist")
	}
	fmt.Println("Initialzation Successful")
	setup.initialized = true
	return nil
}

//安装并实例化链码
func (setup *FabricSetup) InstallAndInstantiateCC() error {
	//创建发送给peers的chaincode package
	ccPkg, err := packager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return errors.WithMessage(err, "failed to create chaincode package")
	}
	fmt.Println("chaincode package created")
	ccHasInstalled := false
	//查询已安装的链码
	ccInstalledRes, err := setup.rmCli.QueryInstalledChaincodes(resmgmt.WithTargetEndpoints(setup.OrgPeer0))
	if err != nil {
		return errors.WithMessage(err, "failed to Query Installed chaincode")
	}
	if ccInstalledRes != nil {
		for _, cc := range ccInstalledRes.Chaincodes {
			if strings.EqualFold(cc.Name, setup.ChaincodeID) {
				ccHasInstalled = true
			}
		}
	}
	fmt.Println("ccHasInstall", ccHasInstalled)
	if !ccHasInstalled {
		//安装链码(智能合约)到org peers
		installCCReq := resmgmt.InstallCCRequest{Name: setup.ChaincodeID, Path: setup.ChaincodePath, Version: setup.ChaincodeVersion, Package: ccPkg}
		_, err = setup.rmCli.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
		if err != nil {
			return errors.WithMessage(err, "failed to install chaincode")
		}
		fmt.Println("chaincode install success")
	} else {
		fmt.Println("Chaincode already exist")
	}
	ccHasInstantiated := false
	//查询已经实例化的链码
	ccInstantiatedResp, err := setup.rmCli.QueryInstantiatedChaincodes(setup.ChannelID, resmgmt.WithTargetEndpoints(setup.OrgPeer0))
	if ccInstantiatedResp.Chaincodes != nil && len(ccInstantiatedResp.Chaincodes) > 0 {
		for _, chaincodeIno := range ccInstantiatedResp.Chaincodes {
			fmt.Println(chaincodeIno)
			if strings.EqualFold(chaincodeIno.Name, setup.ChaincodeID) {
				ccHasInstantiated = true
			}
		}
	}
	// could not get chConfig cache reference:read configuration for channel peers failed

	// Set up chaincode policy
	// ccPolicy := cauthdsl.SignedByAnyMember([]string{"fbi.citizens.com"})
	if !ccHasInstantiated {
		//msp名称，非域名
		ccPolicy := cauthdsl.SignedByAnyMember([]string{"Manufacturer"})
		req := resmgmt.InstantiateCCRequest{Name: setup.ChaincodeID, Path: setup.ChaincodeGoPath, Version: setup.ChaincodeVersion, Policy: ccPolicy}

		resp, err := setup.rmCli.InstantiateCC(setup.ChannelID, req)

		if err != nil || resp.TransactionID == "" {
			return errors.WithMessage(err, "failed to instantiate the chaincode")
		}
		fmt.Println("Chaincode instantiated successed")
	} else {
		fmt.Println("chaincode has instantiated")
	}
	//channel Context用于查询和执行事务交易
	chCtx := setup.sdk.ChannelContext(setup.ChannelID, fabsdk.WithUser(setup.UserName))
	setup.chCli, err = channel.New(chCtx)
	if err != nil {
		return errors.WithMessage(err, "failed to create new Channel client")
	}
	fmt.Println("channel client created")
	//访问通道事件
	setup.chEvent, err = event.New(chCtx)
	if err != nil {
		return errors.WithMessage(err, "failed to create new event client")
	}
	fmt.Println("Event client created")
	fmt.Println("Chaincode Installation and Instantiation successful!")
	return nil
}

func (setup *FabricSetup) CloseSDK() {
	setup.sdk.Close()
}
