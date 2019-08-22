package main

import (
	"fmt"
	"net/http/cgi"
	"os"
	"os/exec"
	"syscall"
)

func cgi_handle(path string, args []string) {
	r, err := cgi.Request()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	// id is of form <port>-<cookie>
	id := r.FormValue("id")
	if id == "" { // first time: let's execute the script
		fmt.Println("Content-Type: text/plain")
		fmt.Println()
		fmt.Println("TODO: exec script in background")
		cmd := exec.Command(path, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		}
		os.Exit(0)
	}
	fmt.Println("Content-Type: application/json")
	fmt.Println()
	fmt.Printf("id = %#v\n", r.FormValue("id"))
	fmt.Println("TODO: connect to backgrounded script")
	// and now, let's execute the script:
}
