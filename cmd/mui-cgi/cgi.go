package main

import (
	"fmt"
	"net/http/cgi"
	"os"
	"os/exec"
	"syscall"
//	"encondig/json"
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
		newargs := []string{"-exec", "-shell", path}
		newargs = append(newargs, args...)
		cmd := exec.Command("mui-cgi", newargs...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var s string
		fmt.Fscanln(stdout, &s)
		fmt.Printf("script is %q\n", args[0])
		fmt.Printf("id is (%s)\n", s)
		fmt.Println("TODO: output HTML template")
		os.Exit(0)
	}
	fmt.Println("Content-Type: application/json")
	fmt.Println()
	fmt.Printf("id = %#v\n", r.FormValue("id"))
	fmt.Println("TODO: connect to backgrounded script")
	// and now, let's execute the script:
}
