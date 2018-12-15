package main

import (
	"fmt"
	"log"
	"net/rpc"

	"./rpcexample" //Path to the package which contains service definition
)

func main() {

	// Try to connect to localhost:1234 using HTTP protocol (the port on which RPC server is listening)
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	// Synchronous call
	args := &rpcexample.Args{5, 6}
	// reply will store the RPC result
	var mulReply rpcexample.Result
	// Call remote procedure
	err = client.Call("Arithmetic.Multiply", args, &mulReply)
	if err != nil {
		log.Fatal("Error in Arithmetic.Multiply: ", err)
	}
	fmt.Printf("Arithmetic.Multiply: %d*%d=%d\n", args.A, args.B, mulReply)

	// Asynchronous call
	args = &rpcexample.Args{501, 100}
	divReply := new(rpcexample.Quotient)
	divCall := client.Go("Arithmetic.Divide", args, divReply, nil)
	divCall = <-divCall.Done
	if divCall.Error != nil {
		log.Fatal("Error in Arithmetic.Divide: ", divCall.Error.Error())
	}
	fmt.Printf("Arithmetic.Divide: %d/%d=%d (rem=%d)\n", args.A, args.B, divReply.Quo, divReply.Rem)
}
