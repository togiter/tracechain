package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
var logger = shim.NewLogger("ProductChaincode")
type ProductChaincode struct {
}

func (pc *ProductChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debugf("[txID %s] ########### ProductChaincode Init ###########\n", stub.GetTxID())
	return shim.Success([]byte("成功初始化"))
}

func (pc *ProductChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	logger.Debugf("[txID %s] ########### ProductChaincode Invoke ###########\n", stub.GetTxID())
	invokeFunc, args := stub.GetFunctionAndParameters()
	
	if invokeFunc != "invoke" {
		return shim.Error("Recevied unknown function invoke invocation")
	}
	// if len(args) < 1 {
	// 	return shim.Error("Recevied args error!!!!")
	// }
	//args := stub.GetStringArgs();
	function := args[0]
	params := args[1:]
	fmt.Printf("function Invoke:%s",function)
	fmt.Printf("Invoke is running:%s" ,function)
	if function == "IssueProduct" { //发布产品
		return pc.IssueProduct(stub, params)
	} else if function == "TransferProduct" { //改变产品所有权(销售)
		return pc.TransferProduct(stub, params)
	} else if function == "AlterProductPrice" { //改变产品价格
		return pc.AlterProductPrice(stub, params)
	} else if function == "QueryProductNo" { //按产品编号查询产品
		return pc.QueryProductNo(stub, params)
	} else if function == "QueryProductRange" { //批量查询产品
		return pc.QueryProductRange(stub, params)
	}
	fmt.Printf("invoke did not find func:%s",function)
	return shim.Error("Recevied unknown function invocation:"+function)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode chaincode: %s", err)
	}
}
