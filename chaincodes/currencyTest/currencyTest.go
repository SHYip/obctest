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
	"strconv"

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"

)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var creationEntity string
	var initValue int
	var key string
	var err error

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// Initialize the chaincode
	creationEntity = args[0]
	
	initValue, err = strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}
	
	key = args[2]
	
	// Write the state to the ledger
	err = stub.PutState(creationEntity, []byte(strconv.Itoa(initValue)))
	err = stub.PutState(creationEntity+"key", []byte(key))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Transaction makes payment of X units from A to B
func (t *SimpleChaincode) invoke(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var cmd string
	
	cmd=args[0]
	
	switch cmd{
		case "create":
			t.createAcc(stub, args)
		case "transfer":
			t.transferFunds(stub, args)
		default:
			return nil, errors.New("Unrecognised command")			
	}
	return nil, nil
}

func (t *SimpleChaincode) createAcc(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var entity string
	var key string
	var err error
	
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}
	
	entity=args[1]
	key=args[2]
	
	err = stub.PutState(entity, []byte(strconv.Itoa(0)))
	if err!= nil {
		return nil, errors.New("Failed to put state")
	}
	err = stub.PutState(entity+"key", []byte(key))
	if err!= nil {
		return nil, errors.New("Failed to put state")
	}
	
	return nil, nil
	
}

func (t *SimpleChaincode) transferFunds(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var A, B string
	var amount int
	var key string
	var Aval, Bval int
	var err error
	
	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}
	A=args[1]
	B=args[2]
	amount, err = strconv.Atoi(args[3])
	key=args[4]
	
	storedKey, err := stub.GetState(A+"key")
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if storedKey == nil {
		return nil, errors.New("Entity not found")
	}
	
	if string(storedKey[:]) != key{
		return nil,errors.New("Incorrect key")
	}
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))
	
	if amount>Aval {
		return nil, errors.New("Insufficient funds")
	}
	
	if amount<0 {
		return nil, errors.New("Please input only positive values")
	}
	
	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if Avalbytes == nil {
		return nil, errors.New("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))
	
	
	Aval-=amount
	Bval+=amount
	
	err = stub.PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(B, []byte(strconv.Itoa(Bval)))
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
	var A string // Entities
	var key string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]
	key = args[1]
	storedKey, err := stub.GetState(A+"key")
	if err != nil {
		return nil, errors.New("Failed to get state")
	}
	if storedKey == nil {
		return nil, errors.New("Entity not found")
	}
	
	if string(storedKey[:]) != key{
		return nil,errors.New("Incorrect key")
	}

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return nil, errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return Avalbytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
