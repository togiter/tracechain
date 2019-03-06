package main

import (
	"fmt"
	"genealogy/fabric-service/fabricSetup"
	"genealogy/web-service"
	"genealogy/web-service/controller"
	"os"
)

func main() {
	fSetup := fabricSetup.FabricSetup{
		OrdererID: "orderer.genealogy.com",

		ChannelID:     "genealogychannel",
		ChannelConfig: os.Getenv("GOPATH") + "/src/genealogy/networkConfig/channel-artifacts/channel.tx",

		ChaincodeID:       "genealoyChaincode",
		ChaincodeVersion:  "v0",
		ChaincodeGoPath:   os.Getenv("GOPATH"),
		ChaincodePath:     "genealogy/fabric-service/chaincode",
		OrgAdmin:          "admin",
		OrgName:           "Org1",
		ConnectionProfile: "fabric-service/fabricSetup/connectionprofile.yaml",

		UserName: "user1",
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
