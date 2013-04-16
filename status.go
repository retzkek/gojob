package gojob

import (
	"fmt"
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
	// parse uptime output (this *should* work with most uptime formats)
	var replyString string
	if strings.Contains(string(out), "load average:") {
		// linux format
		replyString = strings.Split(string(out), "load average:")[1]
	} else if strings.Contains(string(out), "load averages:") {
		// BSD format
		replyString = strings.Split(string(out), "load averages:")[1]
	} else {
		// unknown format
		return fmt.Errorf("trouble parsing uptime string [%s]", string(out))
	}
	replyString = strings.TrimSpace(replyString)
	loadStrings := strings.Split(replyString, ", ")
	if len(loadStrings) != 3 {
		return fmt.Errorf("trouble parsing uptime string [%s]", replyString)
	}
	// one minute
	reply.One, err = strconv.ParseFloat(loadStrings[0], 64)
	if err != nil {
		return err
	}
	// five minute
	reply.Five, err = strconv.ParseFloat(loadStrings[1], 64)
	if err != nil {
		return err
	}
	// fifteen minute
	reply.Fifteen, err = strconv.ParseFloat(loadStrings[2], 64)
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
