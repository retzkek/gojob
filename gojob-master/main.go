package main

import (
	"flag"
)

//////////////////// flags ////////////////////
var webPort int

func init() {
	const (
		defaultWebPort = -1 // no web interface
		usage          = "start web interface on specified port"
	)
	flag.IntVar(&webPort, "web", defaultWebPort, usage)
	flag.IntVar(&webPort, "w", defaultWebPort, usage+" (shorthand)")
}

var addServer = flag.String("add", "", "add server")

//////////////////// main ////////////////////
func main() {
	flag.Parse()

	mongo, err := InitMongo("localhost", "test")
	if err != nil {
		panic(err)
	}

	if *addServer != "" {
		err = mongo.AddNewServer(*addServer)
		if err != nil {
			panic(err)
		}
	}

	if webPort > 0 {
		runWeb(mongo, webPort)
	} else {
		runTerm(mongo)
	}
}
