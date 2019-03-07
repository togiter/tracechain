package chaincode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/tracechain/fabric-service/models"
)

//联盟链
type AllianceChain struct {
}

//实现链码方法

func (ac *AllianceChain) Init(stub shim.ChaincodeStubInterface) peer.Response {
	log.Println("Init successed")
	return shim.Success(nil)
}

func (ac *AllianceChain) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	log.Println("=====Invoke======")
	function, args := stub.GetFunctionAndParameters()
	log.Println("========GetFunctionAndParameters========", function, len(args), args)

	if function == "addMember" {
		return ac.addMember(stub, args)
	} else if function == "queryMember" {
		return ac.queryMember(stub, args)
	} else if function == "changeMember" {
		return ac.changeMember(stub, args)
	} else {
		return shim.Error("function not defined")
	}
	return shim.Success(nil)
}

func (ac *AllianceChain) addMember(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) != 2 {
		return shim.Error("args error")
	}
	key := args[0]
	value := args[1] //json对象
	member := models.Member{}
	err := json.Unmarshal([]byte(value), &member)
	if err != nil {
		return shim.Error("add member failed;parameters cannot be parsed into json objects")
	}
	stub.PutState(key, []byte(value))
	log.Println(key, args)
	return shim.Success(nil)
}

func (ac *AllianceChain) queryMember(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	if len(args) < 1 {
		return shim.Error("args error")
	}
	if len(args) == 1 { //身份证查找
		key := args[0]
		ret, err := stub.GetState(key)
		if err != nil {
			log.Println(fmt.Sprintf("query fail key:%s err:%s", key, err))
			return shim.Error("query fail")
		}
		return shim.Success(ret)
	} else if len(args) == 2 { //范围查找[start,lenght]
		start := string(args[0])
		end := string(args[1])
		iteRet, err := stub.GetStateByRange(start, end)
		if err != nil {
			return shim.Error("range query failed!")
		}
		defer iteRet.Close()
		//解释范围结果
		var buf bytes.Buffer
		buf.WriteString("[")
		alreadyWriten := false
		for iteRet.HasNext() {
			queryResp, err := iteRet.Next()
			if err != nil {
				return shim.Error("range resoved failed")
			}
			if alreadyWriten == true {
				buf.WriteString(",")
			}
			buf.WriteString("{\"Key\":")
			buf.WriteString("\"")
			buf.WriteString(queryResp.Key)
			buf.WriteString("\"")
			buf.WriteString(",\"Record\":")
			buf.WriteString(string(queryResp.Value))
			buf.WriteString("}")
			alreadyWriten = true
		}
		buf.WriteString("]")
		fmt.Printf("range query:\n%s\n", buf.String())
		return shim.Success(buf.Bytes())
	}
	return shim.Success(nil)
}

func (ac *AllianceChain) changeMember(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	fmt.Println("function has not finished")
	return shim.Success(nil)
}
func main() {
	err := shim.Start(new(AllianceChain))
	if err != nil {
		log.Println(err)
	}

}
