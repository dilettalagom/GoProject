package main

import (
	"./rpcservice"
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

//utils
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func convertToInt(num_s string) int64 {

	num_int, err := strconv.ParseInt(num_s, 10, 64)
	check(err)
	return num_int

}

func main() {

	// Try to connect to localhost:1234 (the port on which RPC server is listening)
	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error in dialing: ", err)
	}
	defer client.Close()

	//get input from command line
	if len(os.Args[1:]) < 3 {
		fmt.Println("Error in input. Usage: " + os.Args[0] + "<N of Mappers> <M of Reduces> <directory_path> ")
		os.Exit(1)
	}

	N := int(convertToInt(os.Args[1]))
	M := int(convertToInt(os.Args[2]))

	//get and open Directory
	DIR_PATH, err := os.Open(os.Args[3])
	files, err := ioutil.ReadDir(os.Args[3])
	check(err)

	//creating a worker for each file.txt
	for _, f := range files {

		name_file := DIR_PATH.Name() + "/" + f.Name()
		file_stream, err := os.Open(name_file)
		check(err)

		file_stream = file_stream
		//setup rpc arg
		args_rpc := rpcservice.Args{}
		args_rpc.File = name_file //perché prima metteva tutta la directory + file (file_stream.Name())
		args_rpc.N = N
		args_rpc.M = M

		//rpc call -> file.txt
		// reply will store the RPC result

		var wordcount rpcservice.Result = make(map[string]int)

		fmt.Printf("filename %s\n", args_rpc.File)
		fmt.Printf("N e M %d - %d\n", args_rpc.N, args_rpc.M)

		// Call remote procedure
		err = client.Call("wordcounter.Map", args_rpc, &wordcount)
		if err != nil {
			log.Fatal("Error in wordcounter.Map: ", err)
		}

		fmt.Printf("_________________________RESULTS of %s_________________________\n", args_rpc.File)
		for key, value := range wordcount {
			fmt.Println("Key:", key, "Value:", value)
		}

	}

}
