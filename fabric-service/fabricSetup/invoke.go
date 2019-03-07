package fabricSetup

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/tracechain/fabric-service/memberChaincode"
	"time"
)

func (sp *FabricSetup) AddMember() (string, error) {
	// Prepare arguments
	var args []string
	args = append(args, "535636789302345673")
	member := memberChaincode.Member{
		DataType:      "citizens",
		Id:            "535636789302345673",
		Sex:           "男",
		Name:          "张三",
		BirthLocation: memberChaincode.Location{Province: "海南", City: "三亚市", Detail: "天涯海角"},
		LiveLocation:  memberChaincode.Location{Province: "北京", Town: "朝阳区", Detail: "大悦城"},
		MotherId:      "535636789302345671",
		FatherId:      "535636789302345672",
		Childs:        []string{"535636789302345674", "535636789302345675"},
	}
	memberBytes, _ := json.Marshal(member)
	eventId := "eventInvoke"
	//添加描述
	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	req, notifier, err := sp.chEvent.RegisterChaincodeEvent(sp.ChaincodeID, eventId)
	if err != nil {
		return "", err
	}
	defer sp.chEvent.Unregister(req)

	//创建并提交提案proposal
	resp, err := sp.chCli.Execute(channel.Request{
		ChaincodeID:  sp.ChaincodeID,
		Fcn:          "AddMember",
		Args:         [][]byte{[]byte(args[0]), memberBytes},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to send propose:%v", err)
	}

	//等待结果返回
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Receved CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("did Not receive CC event for eventId(%s)", eventId)
	}

	return string(resp.TransactionID), nil
}
