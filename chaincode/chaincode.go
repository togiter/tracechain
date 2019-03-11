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
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running" + function)
	if function == "IssueProduct" { //发布产品
		return pc.IssueProduct(stub, args)
	} else if function == "TransferProduct" { //改变产品所有权(销售)
		return pc.TransferProduct(stub, args)
	} else if function == "AlterProductPrice" { //改变产品价格
		return pc.AlterProductPrice(stub, args)
	} else if function == "QueryProductNo" { //按产品编号查询产品
		return pc.QueryProductNo(stub, args)
	} else if function == "QueryProductRange" { //批量查询产品
		return pc.QueryProductRange(stub, args)
	}
	fmt.Println("invoke did not find func:" + function)
	return shim.Error("Recevied unknown function invocation")
}

func main() {
	// Start the chaincode and make it ready for futures requests
	err := shim.Start(new(ProductChaincode))
	if err != nil {
		fmt.Printf("Error starting ProductChaincode chaincode: %s", err)
	}
}
