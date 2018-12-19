package main

import (
	"./rpcservice"
	"log"
	"net"
	"net/rpc"
)

func main() {

	//Create an instance of struct which implements Wordcounter interface
	main := new(rpcservice.Wordcounter)

	// Register a new rpc server and the struct we created above.
	// Only structs which implement Wordcounter interface
	// are allowed to register themselves
	server := rpc.NewServer()
	err := server.RegisterName("wordcounter", main)
	if err != nil {
		log.Fatal("Format of service Wordcounter is not correct: ", err)
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
