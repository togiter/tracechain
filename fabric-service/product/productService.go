package productservice

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/tracechain/fabric-service/fabricSetup"
	"github.com/tracechain/fabric-service/product/productChaincode"
)

//发布产品

func IssueProduct(sdkInfo *fabricSetup.FabricSetup,name string,millPrice string,price string,color string,owner string,productor string) []byte,error {
	if len(number) <= 0 || len(name) <= 0 || len(millPrice) <= 0 || len(price)<= 0 || len(color) <= 0 || len(owner) <= 0 || len(productor) <= 0 {
		return "",fmt.Errorf("args input error!")
	}
	product := productChaincode.Product{
		"product",
		name,
		number,
		millPrice,
		price,
		color,
		owner,
		productor
	}

	productBytes,err := json.Marshal(product)
	if err != nil || productBytes == nil{
		return "",fmt.Errorf("product marshal failed!")
	}

	eventID := "eventOfIssueProduct!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in IssueProduct")
	reg,noti,err := sdkInfo.eventID.RegisterChaincodeEvent(sdkInfo.ChaincodeID,eventID)
	if err != nil {
		return nil, fmt.Errorf("event register error!")
	}
	defer sdkInfo.eventID.Unregister(reg)

	resp,err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:sdkInfo.ChaincodeID,
		Fcn:"invoke",
		Args:[][]byte{[]byte("IssueProduct"),productBytes},
		TransientMap:transientDataMap})
		if err != nil {
			return nil, fmt.Errorf("failed to invoke IssueProduct:%v",err)
		}
		
		//等待事件回调
		select {
		case ccEvent := <-noti:
			fmt.Printf("Received CC event:%v\n",ccEvent)
		case <-time.After(time.Second*20):
			return nil, fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		}
		return string(resp.TransactionID),nil
}
//产品转移
func TransferProduct(sdkInfo *fabricSetup.FabricSetup,nOwner string,number string,price string) []byte,error {
	if len(nOwner) <= 0 || len(number) <= 0 {
		return,nil,fmt.Errorf("args input error!")
	}
	if price == nil{
		price = ""
	}
	
	eventID := "eventOfTransfer!"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in tranfer")
	reg,noti,err := sdkInfo.eventID.RegisterChaincodeEvent(sdkInfo.ChaincodeID,eventID)
	if err != nil {
		return nil, fmt.Errorf("event register error!")
	}
	defer sdkInfo.eventID.Unregister(reg)

	resp,err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:sdkInfo.ChaincodeID,
		Fcn:"invoke",
		Args:[][]byte{[]byte("TransferProduct"),[]byte(nOwner),[]byte(number),[]byte(price)},
		TransientMap:transientDataMap})
		if err != nil {
			return nil, fmt.Errorf("failed to invoke TransferProduct:%v",err)
		}
		
		//等待事件回调
		select {
		case ccEvent := <-noti:
			fmt.Printf("Received CC event:%v\n",ccEvent)
		case <-time.After(time.Second*20):
			return nil, fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		}
		return string(resp.TransactionID),nil
}

//修改价格
func AlterProductPrice(sdkInfo *fabricSetup.FabricSetup,owner string,number string,price string) string, error{
	if len(owner) <= 0 || len(number) <= 0 || len(price) <= 0{
		return nil,fmt.Errorf("args input error!")
	}
	// var args []string
	// args = append(args,"AlterProductPrice")
	// args = append(args,owner)
	// args = append(args,number)
	// args = append(args,price)
	eventID := "eventOfAlterPrice"
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in alter price")
	reg,noti,err := sdkInfo.eventID.RegisterChaincodeEvent(sdkInfo.ChaincodeID,eventID)
	if err != nil {
		return nil, fmt.Errorf("event register error!")
	}
	defer sdkInfo.eventID.Unregister(reg)

	resp,err := sdkInfo.chCli.Execute(channel.Request{
		ChaincodeID:sdkInfo.ChaincodeID,
		Fcn:"invoke",
		Args:[][]byte{[]byte("AlterProductPrice"),[]byte(owner),[]byte(number),[]byte(price)},
		TransientMap:transientDataMap})
		if err != nil {
			return nil, fmt.Errorf("failed to invoke AlterProductPrice:%v",err)
		}
		
		//等待事件回调
		select {
		case ccEvent := <-noti:
			fmt.Printf("Received CC event:%v\n",ccEvent)
		case <-time.After(time.Second*20):
			return nil, fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
		}
		return string(resp.TransactionID),nil
}

//查询指定范围的产品
func QueryProductRange(sdkInfo *fabricSetup.FabricSetup,startKey string ,endKey string)[]byte,error {
	if len(startKey) <= 0 || len(endKey).len <= 0{
		return nil, fmt.Errorf("参数有误！")
	}
	resp,err := sdkInfo.chCli.Query(channel.Request{
		ChaincodeID:sdkInfo.ChaincodeID,
		Fcn:"invoke",
		Args:[][]byte{[]byte("QueryProductRange"),[]byte(startKey),[]byte(endKey)},
	})
	if err != nil {
		return nil,fmt.Println("failed to range query:%v",err)
	}
	fmt.Println("range query:",resp.Payload)
	//products := []Product{}

	return []byte(resp.Payload),nil
}

func QueryProductNo(sdkInfo *fabricSetup.FabricSetup,productNo string) productChaincode.Product,error{
	if len(productNo) <= 0{
		return nil, fmt.Errorf("参数有误！")
	}
	var args []string
	args = append(args,"invoke")
	args = append(args,"QueryProductNo")
	args = append(args,productNo)
	resp,err := sdkInfo.chCli.Query(channel.Request{
		ChaincodeID:sdkInfo.ChaincodeID,
		Fcn:args[0],
		Args:[][]byte{[]byte(args[1]),[]byte(args[2])},
	})
	if err != nil {
		return nil,fmt.Println("failed to query:%v",err)
	}
	product := Product{}
	err = json.Unmarshal([]byte(resp.Payload),&product)
	fmt.Println("productNo query resp:",resp.Payload)
	return product,nil
}