package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type ProductChaincode struct {
}

func (pc *ProductChaincode) Init(stu shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (pc *ProductChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	invokeFunc, args := stub.GetFunctionAndParameters()
	
	if invokeFunc != "invoke" {
		return shim.Error("Recevied unknown function invoke invocation")
	}
	if len(args) < 1 {
		return shim.Error("Recevied args error!!!!")
	}
	function := args[0]
	params := args[:]
	fmt.Println("invoke is running" + function)
	if invokeFunc == "invoke" { //发布产品
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
	fmt.Println("invoke did not find func:" + function)
	return shim.Error("Recevied unknown function invocation:"+function)
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode chaincode: %s", err)
	}
}
