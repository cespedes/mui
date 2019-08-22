package mui

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type whiptail struct{}

func init() {
	RegisterFrontEnd(whiptail{})
}

func (z whiptail) Name() string {
	return "whiptail"
}

func (z whiptail) Priority() int {
	return 50
}

func (z whiptail) Available() bool {
	if os.Getenv("TERM") == "" {
		return false
	}
	if _, err := exec.LookPath("whiptail"); err != nil {
		return false
	}
	fmt.Println("whiptail is available")
	return true
}

func (z whiptail) Question() int {
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

func (z whiptail) Input() string {
	return "text from whiptail"
}
