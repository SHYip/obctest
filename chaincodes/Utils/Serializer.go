package main

import (
	//"errors"
	"fmt"
	//"strconv"
	"bufio"
	"os"
	//"strings"


	//"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
	"github.com/golang/protobuf/proto"
    pb "github.com/openblockchain/obc-peer/openchain/example/chaincode/KYC/protos"
)



func main() {
	var nric, name, address, mobile string 
	scanner := bufio.NewScanner(os.Stdin)
	
	
	fmt.Printf("NRIC:\n")
	scanner.Scan()
	nric=scanner.Text()
	fmt.Printf("Name:\n")
	scanner.Scan()
	name=scanner.Text()
	fmt.Printf("Address:\n")
	scanner.Scan()
	address=scanner.Text()
	fmt.Printf("Mobile:\n")
	scanner.Scan()
	mobile=scanner.Text()
	
	record:=&pb.Record{Name: name, Nric: nric, Address: address, Mobile: mobile}
	fmt.Printf("Record: %v\n", record)
	
	data,_:=proto.Marshal(record)
	test:=fmt.Sprintf("%#V", data)
	fmt.Printf("%s\n",test)
	fmt.Printf("%#v\n\n", data)
	
}
