package main

import (
	"fmt"
	"github.com/retzkek/gojob"
	"html/template"
	"log"
	"net/http"
	"net/rpc"
)

func runWeb(backend Backend, port int) {
	http.HandleFunc("/", indexHandler(backend))
	address := fmt.Sprintf(":%d", port)
	fmt.Printf("Starting server on %s\n", address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func indexHandler(backend Backend) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		servers, err := backend.GetServers()
		if err != nil {
			panic(err)
		}
		for i, server := range servers {
			host := server.Address
			client, err := rpc.DialHTTP("tcp", host+":1234")
			if err != nil {
				servers[i].Status = fmt.Sprintf("%s", err)
			} else {
				var load gojob.Load
				err = client.Call("Status.SystemLoad", i, &load)
				if err == nil {
					//fmt.Fprintf(w, "<td>%4.2f</td></tr>\n", load.Five)
					servers[i].Load = load
					if load.One < 0.5 {
						servers[i].Status = "Idle"
					} else {
						servers[i].Status = "Busy"
					}
					nproc := int(load.One + 1.0)
					procs := make([]gojob.Process, nproc)
					err = client.Call("Status.TopProcesses", nproc, &procs)
					if err != nil {
						log.Print(err)
					} else {
						servers[i].Processes = procs
					}
				}
			}
		}
		t, _ := template.ParseFiles("templates/list.html")
		t.Execute(w, servers)
	}
}
