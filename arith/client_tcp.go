package main

import (
	"fmt"
	"log"
	"net/rpc"

	"./rpcexample" //Path to the package which contains service definition
)

func main() {

	// Try to connect to localhost:1234 (the port on which RPC server is listening)
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()
	

	// Synchronous call
	//args := &rpcexample.Args{0, 0}

	array := rpcexample.Args{make ([] int, 5)}

	x := 0
	for i:=0; i < 5; i++ {
		fmt.Scanf("%d", &x)
		array.Vect[i] = x
	}

	for i:=0; i < 5; i++ {
		fmt.Printf("Elementi: %d\n", array.Vect[i])
	}

	fmt.Printf("struct: %d\n", array)
	fmt.Printf("array: %d\n", array.Vect)

	// reply will store the RPC result
	//var mulReply rpcexample.Result
	var sum rpcexample.Sum
	// Call remote procedure
	err = client.Call("Arithmetic.Sum", array, &sum)
	if err != nil {
		log.Fatal("Error in Arithmetic.Sum: ", err)
	}
	fmt.Printf("Arithmetic.Sum: %d\n", sum)

}

	/* Asynchronous call
	args = &rpcexample.Args{501, 100}
	divReply := new(rpcexample.Quotient)
	divCall := client.Go("Arithmetic.Divide", args, divReply, nil)
	divCall = <-divCall.Done
	if divCall.Error != nil {
		log.Fatal("Error in Arithmetic.Divide: ", divCall.Error.Error())
	}
	fmt.Printf("Arithmetic.Divide: %d/%d=%d (rem=%d)\n", args.A, args.B, divReply.Quo, divReply.Rem)
}*/
