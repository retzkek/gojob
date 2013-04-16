package main

import (
	"fmt"
	"github.com/retzkek/gojob"
	"net/rpc"
)

func main() {
	servers := []string{"mach47", "127.0.0.1", "localhost"}
	var reply gojob.Load
	//          1234567890123456789012345678901234567890
	fmt.Printf("SERVER              LOAD\n")
	fmt.Printf("------------------- ----\n")
	// poll servers
	// TODO: make asynchronous
	for i, server := range servers {
		fmt.Printf("%-20s", server)
		client, err := rpc.DialHTTP("tcp", server+":1234")
		if err != nil {
			//log.Fatal("dialing:", err)
			fmt.Printf("down\n")
		} else {
			err = client.Call("Status.SystemLoad", i, &reply)
			if err == nil {
				fmt.Printf("%4.2f\n", reply.Five)
			} else {
				fmt.Printf("err\n")
			}
		}
	}
}
