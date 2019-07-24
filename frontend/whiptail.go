package frontend

import (
	"os"
	"os/exec"
	"syscall"
)

func init() {
	RegisterFrontEnd(whiptail)
}

type Whiptail struct {
}

func (z Whiptail) Available() bool {
	if os.Getenv("TERM") == "" {
		return false
	}
	return true
}

func (z Whiptail) Question() {
	cmd := exec.Command("whiptail", "--yesno", "Are you sure you wany yo proceed?", "7", "40")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			os.Exit(status.ExitStatus())
		}
	}
	os.Exit(0)
}

func (z Whiptail) Priority() int {
	return 50
}

var whiptail Whiptail
