package frontend

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type zenity struct{}

func init() {
	RegisterFrontEnd(zenity{})
}

func (z zenity) Name() string {
	return "zenity"
}

func (z zenity) Priority() int {
	return 100
}

func (z zenity) Available() bool {
	if os.Getenv("DISPLAY") == "" {
		return false
	}
	if _, err := exec.LookPath("zenity"); err != nil {
		return false
	}
	fmt.Println("zenity is available")
	return true
}

func (z zenity) Question() int {
	fmt.Println("zenity Question")
	cmd := exec.Command("zenity", "--question")
	err := cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus()
		}
	}
	return 0
}

func (z zenity) Input() string {
	return "text from zenity"
}
