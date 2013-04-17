package main

import (
	"fmt"
	"github.com/retzkek/gojob"
	"log"
	"net/http"
	"net/rpc"
	"path/filepath"
)

func runWeb(servers []Server, port int) {
	http.HandleFunc("/", indexHandler(servers))
	log.Fatal(http.ListenAndServe(":9999", nil))
}

func indexHandler(servers []Server) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<html><head><title>Status</title><body>\n")
		fmt.Fprintf(w, "<table><tr><th>Server</th><th>Status</th></tr>\n")
		for i, server := range servers {
			host := server.Hostname
			fmt.Fprintf(w, "<tr><td>%s</td>", host)
			client, err := rpc.DialHTTP("tcp", host+":1234")
			if err != nil {
				fmt.Fprintf(w, "<td>%s</td></tr>\n", err)
			} else {
				var load gojob.Load
				err = client.Call("Status.SystemLoad", i, &load)
				if err == nil {
					fmt.Fprintf(w, "<td>%4.2f</td></tr>\n", load.Five)
					nproc := int(load.One + 1.0)
					procs := make([]gojob.Process, nproc)
					err = client.Call("Status.TopProcesses", nproc, &procs)
					if err == nil {
						for _, ps := range procs {
							fmt.Fprintf(w, "<tr><td></td><td>%s running %s for %s</td></tr>\n", ps.Owner,
								filepath.Base(ps.Exe), ps.Time)
						}
					}
				} else {
					fmt.Fprintf(w, "<tr><td></td><td></td></tr>%s\n", err)
				}
			}
		}
		fmt.Fprintf(w, "</table></body></html>")
	}
}
