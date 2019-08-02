package frontend

import (
	"os"
	"fmt"
	"syscall"
	"os/exec"
)

func init() {
	RegisterFrontEnd(zenity)
}

type Zenity struct {
}

func (z Zenity) Available() bool {
	if os.Getenv("DISPLAY") == "" {
		return false
	}
	if _, err := exec.LookPath("zenity"); err != nil {
		return false
	}
	fmt.Println("Zenity is available")
	return true
}

func (z Zenity) Question() int {
	fmt.Println("Zenity Question")
	cmd := exec.Command("zenity", "--question")
	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func (z Zenity) Priority() int {
	return 100
}

var zenity Zenity
