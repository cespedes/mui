package frontend

import (
	"os"
	"fmt"
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
        if _, err := exec.LookPath("whiptail"); err != nil {
                return false
        }
        fmt.Println("Whiptail is available")
	return true
}

func (z Whiptail) Question() int {
	cmd := exec.Command("whiptail", "--yesno", "Are you sure you wany yo proceed?", "7", "40")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func (z Whiptail) Priority() int {
	return 50
}

var whiptail Whiptail
