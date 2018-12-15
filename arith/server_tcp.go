package main

import (
	"log"
	"net"
	"net/rpc"

	"./rpcexample" //Path to the package which contains service definition
)

func main() {

	//Create an instance of struct which implements Arith interface
	arith := new(rpcexample.Arith)

	// Register a new rpc server and the struct we created above.
	// Only structs which implement Arith interface
	// are allowed to register themselves
	server := rpc.NewServer()
	err := server.RegisterName("Arithmetic", arith)
	if err != nil {
		log.Fatal("Format of service Arith is not correct: ", err)
	}

	// Listen for incoming tcp packets on specified port.
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error:", e)
	}

	// Link rpc server to the socket, and allow rpc server to accept
	// rpc requests coming from that socket.
	server.Accept(l)
}
