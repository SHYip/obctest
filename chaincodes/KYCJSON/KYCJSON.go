package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"	


	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"

)

type SimpleChaincode struct {
}

type Rating struct {
	Risk string `json:"Risk"`
	Qualifier string `json:"Qualifier"`
}

type Record struct {
	User string `json:"User"`
	Nric string `json:"Nric"`
	Name string `json:"Name"`
	Address string `json:"Address"`
	Mobile string `json:"Mobile"`
	RiskRating Rating `json:"RiskRating"`
}


type Response struct {
	Records []Record `json:"Records"`
}


func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	

	return nil, nil
}


func (t *SimpleChaincode) invoke(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//expecting input in the form {"JSONDATA"}
	
	if(len(args) != 1) {
		return nil, errors.New("Incorrect number of arguments, expecting 10 input")
	}
	
	data:= stringToByteArray(args[0])
	var input Record
	var err error
	json.Unmarshal(data, &input)
	
	nric := input.Nric
	user := input.User
	
	err = stub.PutState(nric+user, data)
	if(err!=nil){
		return nil, err
	}
	
	return nil, nil
	
	

}

func (t *SimpleChaincode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Run callback representing the invocation of a chaincode
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" {
		// Initialize the entities and their asset holdings
		return t.init(stub, args)
	} else if function == "invoke" {
		// Transaction makes payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	}
    
	return nil, errors.New("Received unknown function invocation")
}

func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting nric of the person to query")
	}

	nric := args[0]
	var response Response
	
	keysIter, err := stub.RangeQueryState(nric, nric+"ZZZZZ")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()
	
	var nextRecord Record
	
	for keysIter.HasNext() {
		_, nextRecordBytes, err := keysIter.Next()
		json.Unmarshal(nextRecordBytes, &nextRecord)
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		response.Records=append(response.Records, nextRecord)
	}
	dataBytes,_:=json.Marshal(response)
	
	//dataJSON := json.unmarshal(dataBytes) 
	return dataBytes, nil
}

func stringToByteArray(input string) ([]uint8) {
	var output []uint8
	
	inputArray:=strings.Split(input, ", ")
	for _, i := range inputArray {
		j,_:=strconv.ParseUint(i, 0, 8)
		k:=uint8(j)
		output=append(output,k)
	}
		return output
	
}



func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
