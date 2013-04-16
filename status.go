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

// Process stores indentifying information on a process.
type Process struct {
	Owner string
	Exe   string
	Cpu   float64
	Mem   float64
	Time  string
}

// TopProcesses return details on the top processes based on cpu usage.
// The argument specifies the number of processes to return (not implemented, 
// currently will always return one process).
func (t *Status) TopProcesses(arg int, reply *[]Process) error {
	if arg < 1 {
		return fmt.Errorf("invalid number of processes: %i", arg)
	}
	out, err := exec.Command("ps", "auxk-c", "--no-headers").Output()
	if err != nil {
		return err
	}
	processes := strings.SplitN(string(out), "\n", arg)
	if len(processes) != arg {
		return fmt.Errorf("error parsing ps output")
	}
	*reply = make([]Process, arg)
	for i, ps := range processes {
		fields := strings.Fields(ps)
		if len(fields) < 11 {
			return fmt.Errorf("error parsing ps output")
		}
		(*reply)[i].Owner = fields[0]
		(*reply)[i].Exe = fields[10]
		(*reply)[i].Cpu, err = strconv.ParseFloat(fields[2], 64)
		if err != nil {
			return err
		}
		(*reply)[i].Mem, err = strconv.ParseFloat(fields[3], 64)
		if err != nil {
			return err
		}
		(*reply)[i].Time = fields[9]
	}
	return nil
}
