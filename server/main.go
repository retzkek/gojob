package main

import (
	"github.com/retzkek/gojob"
	"log"
	"net"
	"net/rpc"
	"net/http"
)

func main() {
	arith := new(gojob.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
			log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

