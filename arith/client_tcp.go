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

	/*TODO: prova RPC con somma di un array*/
	array := rpcexample.Args{make([]int, 5)}

	fmt.Printf("Inserire 5 valori\n")
	x := 0
	for i := 0; i < 5; i++ {
		fmt.Scanf("%d", &x)
		array.Vect[i] = x
	}

	for i := 0; i < 5; i++ {
		fmt.Printf("Elementi: %d\n", array.Vect[i])
	}

	fmt.Printf("struct: %d\n", array)
	fmt.Printf("array: %d\n", array.Vect)

	// reply will store the RPC result
	var sum rpcexample.Sum
	// Call remote procedure
	err = client.Call("arithmetic.Somma", array, &sum)
	if err != nil {
		log.Fatal("Error in arithmetic.Somma: ", err)
	}
	fmt.Printf("arithmetic.Somma: %d\n", sum)

	/* TODO:Synchronous call*/
	args := &rpcexample.Fattori{5, 6}
	// reply will store the RPC result
	var mulReply rpcexample.Result
	// Call remote procedure
	err = client.Call("arithmetic.Multiply", args, &mulReply)
	if err != nil {
		log.Fatal("Error in arithmetic.Multiply: ", err)
	}
	fmt.Printf("arithmetic.Multiply: %d*%d=%d\n", args.A, args.B, mulReply) /**/

	/* TODO:Asynchronous call
	args = &rpcexample.Args{501, 100}
	divReply := new(rpcexample.Quotient)
	divCall := client.Go("arithmetic.Divide", args, divReply, nil)
	divCall = <-divCall.Done
	if divCall.Error != nil {
		log.Fatal("Error in arithmetic.Divide: ", divCall.Error.Error())
	}
	fmt.Printf("arithmetic.Divide: %d/%d=%d (rem=%d)\n", args.A, args.B, divReply.Quo, divReply.Rem)
	}*/
}
