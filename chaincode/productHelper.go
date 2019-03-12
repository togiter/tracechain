package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	// "github.com/tracechain/fabric-service/product"
)

type Product struct {
	ObjectType string `json:"objectType"`
	Name       string `json:"name"`
	Number     string `json:"number"`    //产品编号
	MillPrice  string `json:"millPrice"` //出厂价格，不可改变
	Price      string `json:"price"`
	Color      string `json:"color"`
	Owner      string `json:"owner"`     //产品拥有者
	Productor  string `json:"productor"` //厂家
}


func (pc *ProductChaincode) IssueProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 7 { //Name,number,millPrice,price,color,owner,productor
		return shim.Error("args count error")
	}
	fmt.Println("call IssueProduct.....")
	name := args[0]
	number := args[1]
	millPrice := args[2]
	price := args[3]
	color := args[4]
	owner := args[5]
	productor := args[6]
	//check if product already exists
	productBytes, err := stub.GetState(number)
	if err != nil {
		return shim.Error("failed to get product:" + err.Error())
	} else if productBytes != nil {
		fmt.Println("product No already exists:" + number)
		return shim.Error("product No already exists:" + number)
	}

	objectType := "product"
	product := Product{objectType, name, number, millPrice, price, color, owner, productor}
	productBytes, err = json.Marshal(product)
	if err != nil {
		fmt.Println("product marshal failed", err)
		return shim.Error(err.Error())
	}
	err = stub.PutState(number, productBytes)
	if err != nil {
		fmt.Println("putstate prudct failed")
		return shim.Error(err.Error())
	}
	fmt.Println(" IssueProduct.....succeed")
	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	//indexName := "number~name"
	//numberNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{product.number, prodcut.Name})
	//if err != nil {
	//	return shim.Error(err.Error())
	//}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//value := []byte{0x00}
	//stub.PutState(numberNameIndexKey, value)
	return shim.Success(nil)
}

func (pc *ProductChaincode) TransferProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//0			1 		2
	//newowner  number   price
	if len(args) < 2 {
		return shim.Error("Incorrect nember of arguments,Expecting 3")
	}
	nOwner := args[0]
	number := args[1]

	fmt.Println("--start tansfer ", nOwner, number)
	valBytes, err := stub.GetState(number)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if valBytes == nil {
		return shim.Error("product number does not exist")
	}
	prodcutTo := Product{}
	err = json.Unmarshal(valBytes, &prodcutTo)
	if err != nil {
		return shim.Error(err.Error())
	}
	prodcutTo.Owner = nOwner
	if len(args) >= 3 {
		price := args[2]
		prodcutTo.Price = price
	}
	productBytes, _ := json.Marshal(prodcutTo)
	err = stub.PutState(number, productBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (pc *ProductChaincode) AlterProductPrice(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//0		 1	   2
	//owner	number	newPrice
	if len(args) != 3 {
		return shim.Error("Incorrect number of argments")
	}
	owner := args[0]
	number := args[1]
	nPrice := args[2]
	if len(owner) <= 0 {
		shim.Error("owner can't be empty")
	}
	if len(number) <= 0 {
		shim.Error("number can't be empty")
	}
	if len(nPrice) <= 0 {
		shim.Error("newprice can't be empty")
	}
	productBytes, err := stub.GetState(number)
	if err != nil {
		return shim.Error("Failed to get product:" + err.Error())
	} else if productBytes == nil {
		return shim.Error("product does not exist")
	}
	newProduct := Product{}
	err = json.Unmarshal(productBytes, &newProduct)
	if err != nil {
		return shim.Error(err.Error())
	}
	if owner != newProduct.Owner {
		fmt.Println("owner pemissionless changed")
		return shim.Error("owner pemissionless changed")
	}
	newProduct.Price = nPrice
	newProductBytes, _ := json.Marshal(newProduct)
	err = stub.PutState(number, newProductBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("price alter finished")
	return shim.Success(nil)
}

func (pc *ProductChaincode) QueryProductNo(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments")
	}
	var jsonResp string
	number := args[0]
	valBytes, err := stub.GetState(number)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + number + "\"}"
		return shim.Error(jsonResp)
	} else if valBytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + number + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(valBytes)
}

func (pc *ProductChaincode) QueryProductRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments")
	}
	startKey := args[0]
	endKey := args[1]
	retsItr, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer retsItr.Close()
	buf, err := constructQueryResponseFromIterator(retsItr)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- get products by range:", buf.String())
	return shim.Success(buf.Bytes())
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}
