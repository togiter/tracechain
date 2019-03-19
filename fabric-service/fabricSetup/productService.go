package fabricSetup

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/tracechain/fabric-service/product"
	"time"
)

//发布产品

func (sdkInfo *FabricSetup) IssueProduct(name string, number string, millPrice string, price string, color string, owner string, productor string) (string, error) {
	if len(number) <= 0 || len(name) <= 0 || len(millPrice) <= 0 || len(price) <= 0 || len(color) <= 0 || len(owner) <= 0 || len(productor) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	product := product.ProductC{
		ObjectType: "productC",
		Name:       name,
		Number:     number,
		MillPrice:  millPrice,
		Price:      price,
		Color:      color,
		Owner:      owner,
		Productor:  productor,
	}

	productBytes, err := json.Marshal(product)
	if err != nil || productBytes == nil {
		return "", fmt.Errorf("product marshal failed!")
	}

	eventID := "eventOfIssueProduct!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in IssueProduct")
	reg, noti, err := sdkInfo.chEvent.RegisterChaincodeEvent(sdkInfo.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer sdkInfo.chEvent.Unregister(reg)
	fmt.Println("productBytes",productBytes)
	var params [][]byte 
	params = append(params,[]byte("IssueProduct"))
	params = append(params,[]byte(number))
	params = append(params,productBytes)
	resp, err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:  sdkInfo.ChaincodeID,
		Fcn:          "invoke",
		Args:         params,
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke issueProduct:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)\n", eventID)
	}
	return string(resp.TransactionID), nil
}

//产品转移
func (sdkInfo *FabricSetup) TransferProduct(nOwner string, number string, price string) (string, error) {
	if len(nOwner) <= 0 || len(number) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	eventID := "eventOfTransfer!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in tranfer")
	reg, noti, err := sdkInfo.chEvent.RegisterChaincodeEvent(sdkInfo.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer sdkInfo.chEvent.Unregister(reg)

	resp, err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:  sdkInfo.ChaincodeID,
		Fcn:          "invoke",
		Args:         [][]byte{[]byte("TransferProduct"), []byte(nOwner), []byte(number), []byte(price)},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke TransferProduct:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)", eventID)
	}
	return string(resp.TransactionID), nil
}

//修改价格
func (sdkInfo *FabricSetup) AlterProductPrice(owner string, number string, price string) (string, error) {
	if len(owner) <= 0 || len(number) <= 0 || len(price) <= 0 {
		return "", fmt.Errorf("args input error!")
	}
	// var args []string
	// args = append(args,"AlterProductPrice")
	// args = append(args,owner)
	// args = append(args,number)
	// args = append(args,price)
	eventID := "eventOfAlterPrice"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in alter price")
	reg, noti, err := sdkInfo.chEvent.RegisterChaincodeEvent(sdkInfo.ChaincodeID, eventID)
	if err != nil {
		return "", fmt.Errorf("event register error!")
	}
	defer sdkInfo.chEvent.Unregister(reg)

	resp, err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:  sdkInfo.ChaincodeID,
		Fcn:          "invoke",
		Args:         [][]byte{[]byte("AlterProductPrice"), []byte(owner), []byte(number), []byte(price)},
		TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("failed to invoke AlterProductPrice:%v", err)
	}

	//等待事件回调
	select {
	case ccEvent := <-noti:
		fmt.Printf("Received CC event:%v\n", ccEvent)
	case <-time.After(time.Second * 20):
		fmt.Printf("did NOT receive CC event for eventId(%s)", eventID)
	}
	return string(resp.TransactionID), nil
}

//查询指定范围的产品
func (sdkInfo *FabricSetup) QueryProductRange(startKey string, endKey string) (string, error) {
	if len(startKey) <= 0 || len(endKey) <= 0 {
		return "", fmt.Errorf("参数有误！")
	}
	resp, err := sdkInfo.chCli.Query(channel.Request{
		ChaincodeID: sdkInfo.ChaincodeID,
		Fcn:         "invoke",
		Args:        [][]byte{[]byte("QueryProductRange"), []byte(startKey), []byte(endKey)},
	})
	if err != nil {
		return "", fmt.Errorf("failed to range query:%v", err)
	}
	fmt.Println("range query:", resp.Payload)
	//products := []Product{}

	return string(resp.Payload), nil
}

func (sdkInfo *FabricSetup) QueryProductNo(productNo string) (string, error) {
	if len(productNo) <= 0 {
		return "", fmt.Errorf("参数有误！")
	}
	var args []string
	args = append(args, "invoke")
	args = append(args, "QueryProductNo")
	args = append(args, productNo)
	resp, err := sdkInfo.chCli.Query(channel.Request{
		ChaincodeID: sdkInfo.ChaincodeID,
		Fcn:         args[0],
		Args:        [][]byte{[]byte(args[1]), []byte(args[2])},
	})
	if err != nil {
		return "", fmt.Errorf("failed to query:%v", err)
	}

	fmt.Println("productNo query resp:", resp.Payload)
	return string(resp.Payload), nil
}
