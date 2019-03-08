package fabricSetup

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/tracechain/fabric-service/memberChaincode"
)

func (fs *FabricSetup) QueryMember(memberId string) (string, error) {
	if len(memberId) <= 0 {
		return "", fmt.Errorf("input memeberId missed")
	}
	resp, err := fs.chCli.Query(channel.Request{
		ChaincodeID: fs.ChaincodeID,
		Fcn:         "queryMember",
		Args:        [][]byte{[]byte(memberId)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to query :%v", err)
	}
	member := memberChaincode.Member{}
	err = json.Unmarshal([]byte(resp.Payload), &member)
	fmt.Printf(string(resp.Payload))
	return string(resp.Payload), nil
}
