package gojob

import (
	"os/exec"
	"strconv"
	"strings"
)

type Status int

// Load represents a system load storing the one, five, and fifteen 
// minute load averages.
type Load struct {
	One, Five, Fifteen float64
}

// SystemLoad returns the system load as reported by uptime. 
func (t *Status) SystemLoad(arg int, reply *Load) error {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		return err
	}
	var replyString string
	// one minute
	replyString = strings.Split(string(out), ",")[2]
	replyString = strings.Split(strings.TrimSpace(replyString), " ")[2]
	reply.One, err = strconv.ParseFloat(strings.TrimSpace(replyString), 64)
	if err != nil {
		return err
	}
	// five minute
	replyString = strings.Split(string(out), ",")[3]
	reply.Five, err = strconv.ParseFloat(strings.TrimSpace(replyString), 64)
	if err != nil {
		return err
	}
	// fifteen minute
	replyString = strings.Split(string(out), ",")[4]
	reply.Fifteen, err = strconv.ParseFloat(strings.TrimSpace(replyString), 64)
	if err != nil {
		return err
	}
	return nil
}

// Uptime returns the output of the uptime command.
func (t *Status) Uptime(arg int, reply *[]byte) error {
	out, err := exec.Command("uptime").Output()
	if err != nil {
		return err
	}
	*reply = out
	return nil
}
