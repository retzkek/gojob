package main

// terminal-based interface

import (
	"fmt"
	"github.com/retzkek/gojob"
	"net/rpc"
	"path/filepath"
)

func runTerm(backend Backend) {
	var reply gojob.Load
	//          1234567890123456789012345678901234567890
	fmt.Printf("SERVER              LOAD MESSAGE\n")
	fmt.Printf("------------------- ---- ---------------\n")
	// poll servers
	// TODO: make asynchronous
	servers, err := backend.GetServers()
	if err != nil {
		panic(err)
	}
	for i, server := range servers {
		host := server.Address
		fmt.Printf("%-20s", host)
		client, err := rpc.DialHTTP("tcp", host+":1234")
		if err != nil {
			//log.Fatal("dialing:", err)
			fmt.Printf("---- %s\n", err)
		} else {
			err = client.Call("Status.SystemLoad", i, &reply)
			if err == nil {
				fmt.Printf("%4.2f\n", reply.Five)
				nproc := int(reply.One + 1.0)
				procs := make([]gojob.Process, nproc)
				err = client.Call("Status.TopProcesses", nproc, &procs)
				if err == nil {
					for _, ps := range procs {
						fmt.Printf("   %s running %s for %s\n", ps.Owner,
							filepath.Base(ps.Exe), ps.Time)
					}
				}
			} else {
				fmt.Printf("---- %s\n", err)
			}
		}
	}
}
