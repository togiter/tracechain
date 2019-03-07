package main

import (
	"fmt"
	"os"

	"github.com/tracechain/fabric-service/fabricSetup"
	webServer "github.com/tracechain/web-service"
	"github.com/tracechain/web-service/controller"
)

func main() {
	fSetup := fabricSetup.FabricSetup{
		OrdererID:         "orderer.tracechain.com",
		ChannelID:         "tracechain",
		ChannelConfig:     os.Getenv("GOPATH") + "/src/github.com//tracechain/network/artifacts/tracechain.tx",
		ChaincodeID:       "tracechaincode",
		ChaincodeVersion:  "v0",
		ChaincodeGoPath:   os.Getenv("GOPATH"),
		ChaincodePath:     "github.com/tracechain/fabric-service/chaincode",
		OrgAdmin:          "Admin",
		OrgName:           "ManufacturerMSP",
		OrgPeer0:          "peer0.manufacturer.tracechain.com",
		ConnectionProfile: "./fabric-service/fabricSetup/connectionprofile.yaml",
		UserName:          "User1",
	}
	err := fSetup.Initialize()
	if err != nil {
		fmt.Errorf("Fabric SDK初始化失败:%v", err)
		return
	}

	// Close SDK
	defer fSetup.CloseSDK()

	// Install and instantiate the chaincode
	err = fSetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Printf("Unable to install and instantiate the chaincode: %v\n", err)
		return
	}

	// Launch the web application listening
	app := controller.Application{
		Fabric: &fSetup,
	}
	webServer.WebStart(&app)
	// fSetup.AddMember(nil)
	// fSetup.QueryMember("123456")
}
