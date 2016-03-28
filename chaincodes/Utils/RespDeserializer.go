package main

import (
	//"errors"
	"fmt"
	"strconv"
	"bufio"
	"os"
	"strings"


	//"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
	"github.com/golang/protobuf/proto"
    pb "github.com/openblockchain/obc-peer/openchain/example/chaincode/KYC/protos"
)

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
	var inputString string
	scanner := bufio.NewScanner(os.Stdin)
	
	
	fmt.Printf("Input String:\n")
	scanner.Scan()
	inputString=scanner.Text()
	inputBytes:=stringToByteArray(inputString)
	response:=&pb.Response{}
	proto.Unmarshal(inputBytes, response)
	
	fmt.Printf("Record: %v\n", response)
		
}
