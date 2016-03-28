package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"


	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
	"github.com/golang/protobuf/proto"
    pb "github.com/openblockchain/obc-peer/openchain/example/chaincode/KYC/protos"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	

	return nil, nil
}


func (t *SimpleChaincode) invoke(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	//expecting input in the form {"NRIC", "SERIALIZEDDATA", "USER"}
	
	if(len(args) != 3) {
		return nil, errors.New("Incorrect number of arguments, expecting 3 inputs")
	}
	
	var	nric string
	var dataInput string
	var user string
	var data []byte
	var err error
	
	nric=args[0]
	dataInput=args[1]
	user=args[2]
	
	data=stringToByteArray(dataInput)
	
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
	
	//TAKE IN SERIALIZED PROTOBUF, UNMARSHAL, COMPARE TO OTHER ENTRIES WITH SAME NRIC, UPDATE STATUS FOR MATCHES/CLASHES, RETURN STATUS
	var inputString string
	var nric string
	var err error
	

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting nric of the person to query")
	}

	status:=&pb.Response{}
	inputRecord:=&pb.Record{}
	inputString = args[0]
	inputBytes:=stringToByteArray(inputString)
	proto.Unmarshal(inputBytes,inputRecord)
	
	nric=inputRecord.Nric
	
	
	keysIter, err := stub.RangeQueryState(nric, nric+"ZZZZ")
	if err != nil {
		return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
	}
	defer keysIter.Close()
	
	nextRecord := &pb.Record{}
	//var nextBytes []uint8
	if(keysIter.HasNext()==false) {
		status.Nric=pb.Response_NOENTRY
	} else {
		status.Nric=pb.Response_MATCH
	}	
	
	for keysIter.HasNext() {
		_, nextString, err := keysIter.Next()
		if err != nil {
			return nil, fmt.Errorf("keys operation failed. Error accessing state: %s", err)
		}
		//nextRecord=&pb.Record{}
		//nextBytes=stringToByteArray(nextString)
		proto.Unmarshal(nextString, nextRecord)
		status=updateStatus(inputRecord, nextRecord, status)
	}
	
	returnStatus,_:=proto.Marshal(status)
	returnString:=fmt.Sprintf("%#v", returnStatus)
	//TESTING
	if(status.Name == pb.Response_MATCH && status.Address == pb.Response_MATCH && status.Mobile == pb.Response_MATCH) {
		return []byte(returnString), nil
	} else {
		return []byte(returnString), nil	
	}
	//jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	//fmt.Printf("Query Response:%s\n", jsonResp)
	//return Avalbytes, nil
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

func updateStatus(input *pb.Record, compare *pb.Record, status *pb.Response) (*pb.Response) {
	if(input.Name == compare.Name) {
		if(status.Name == pb.Response_NOENTRY) {
			status.Name = pb.Response_MATCH
		}
	} else if(status.Name == pb.Response_NOENTRY){ 
		status.Name = pb.Response_NOMATCH
	}
		
	if(input.Address == compare.Address) {
		if(status.Address == pb.Response_NOENTRY) {
			status.Address = pb.Response_MATCH
		}
	} else if(status.Address == pb.Response_NOENTRY){
		status.Address = pb.Response_NOMATCH
	}
		
	if(input.Mobile == compare.Mobile) {
		if(status.Mobile == pb.Response_NOENTRY) {
			status.Mobile = pb.Response_MATCH
		}
	} else if(status.Mobile == pb.Response_NOENTRY){
		status.Mobile = pb.Response_NOMATCH
	}
		
	return status			
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
