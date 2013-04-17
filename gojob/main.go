package main

import (
	"flag"
	"time"
)

type Server struct {
	Hostname  string
	Status    string
	Timestamp time.Time
}

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

//////////////////// main ////////////////////
func main() {
	flag.Parse()

	servers := make([]Server, 3)
	servers[0].Hostname = "mach47"
	servers[1].Hostname = "127.0.0.1"
	servers[2].Hostname = "localhost"

	if webPort > 0 {
		runWeb(servers, webPort)
	} else {
		runTerm(servers)
	}
}
