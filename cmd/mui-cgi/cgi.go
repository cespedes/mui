package main

import (
	"os"
	"fmt"
	"syscall"
	"os/exec"
	"net/http/cgi"
)

func cgi_handle(path string, args []string) {
	fmt.Println("Content-Type: text/plain")
	fmt.Println()
	r,err := cgi.Request()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
//	err = r.ParseForm()
//	if err != nil {
//		fmt.Printf("Error: %s\n", err)
//	}
	// r.ParseForm & r.Form["id"]
	// r.FormValue("id")
//	fmt.Printf("request = %#v\n", r)
//	fmt.Printf("URL = %#v\n", r.URL)
	fmt.Printf("id = %#v\n", r.FormValue("id"))
	// and now, let's execute the script:
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
