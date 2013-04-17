package main

import (
	"github.com/retzkek/gojob"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	status := new(gojob.Status)
	rpc.Register(status)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
