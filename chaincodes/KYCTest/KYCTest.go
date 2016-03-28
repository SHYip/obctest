/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"

)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	
	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var NRIC string
	var serializedDetails string
	var err error
	
	NRIC=args[0]
	serializedDetails=args[1]
	err = stub.PutState(NRIC, []byte(serializedDetails))
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}


// Deletes an entity from state
func (t *SimpleChaincode) delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Run callback representing the invocation of a chaincode
// This chaincode will manage two accounts A and B and will transfer X units from A to B upon invoke
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

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}
	var NRIC string // Entities
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting NRIC of the person to query")
	}

	NRIC = args[0]
	serializedDetails, err := stub.GetState(NRIC)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}

	jsonResp := "{\"NRIC\":\"" + NRIC + "\",\"Serialized Data\":\"" + string(serializedDetails[:]) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return serializedDetails, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
