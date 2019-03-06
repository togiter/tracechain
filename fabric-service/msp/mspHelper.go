package msp

import (
	"fmt"

	"github.com/cloudflare/cfssl/log"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type MspHandlers struct {
	sdk      *fabsdk.FabricSDK
	orgName  string
	orgAdmin string
}

func (mh *MspHandlers) mockClientProvider() context.ClientProvider {
	log.SetLogger(nil)
	sdk := mh.sdk
	logging.SetLevel("fabsdk/fab", logging.ERROR)
	return sdk.Context()
}

func (mh *MspHandlers) EnrollUser(username string) string {
	if username == nil || len(username) <= 0 {
		fmt.Println("parameter username error")
		return ""
	}
	ctx := mockClientProvider()
	//create msp client
	mspCli, err := mspclient.New(ctx)
	if err != nil {
		fmt.Println("failed to create msp client")
		return ""
	}
	enrollmentSecret, err := mspCli.Register(&RegistrationRequest{Name: username})
	if err != nil {
		fmt.Println("Register return error %s\n", err)
		return ""
	}
	err = mspCli.Enroll(username, mspCli.WithSecret(enrollmentSecret))
	if err != nil {
		fmt.Println("failed to enroll user:%s\n", err)
		return ""
	}
	fmt.Println("enroll user is completed")
	return enrollmentSecret
}
