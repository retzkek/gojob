package main

import (
	"fmt"
	"github.com/retzkek/gojob"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type Command struct {
	Command string
	Short   string
	Run     func(cmd *Command, args []string)
	//Flags flag.FlagSet
}

func printHelp(exe string, commands []Command) {
	fmt.Printf("\nusage: %s command [args]\n", exe)
	fmt.Printf("\tcommands:\n")
	for _, c := range commands {
		fmt.Printf("\t\t%-10s %s\n", c.Command, c.Short)
	}
	fmt.Printf("\n")
}

func runHelp(cmd *Command, args []string) {
	fmt.Println("TBA")
}

func runRpc(cmd *Command, args []string) {
	status := new(gojob.Status)
	rpc.Register(status)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}

func runMongo(cmd *Command, args []string) {
	fmt.Println("TBA")
}

func runTest(cmd *Command, args []string) {
	fmt.Println("TBA")
}

func main() {
	commands := []Command{
		Command{"help", "command help", runHelp},
		Command{"rpc", "run RPC server", runRpc},
		Command{"mongo", "send status to MongoDB server", runMongo},
		Command{"test", "run diagnostics and print to terminal", runTest}}
	if len(os.Args) < 2 {
		printHelp(os.Args[0], commands)
		os.Exit(0)
	}
	cmd := os.Args[1]
	for _, c := range commands {
		if cmd == c.Command {
			c.Run(&c, os.Args[2:])
			break
		}
	}
}
