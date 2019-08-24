package main

import (
	"io"
	"fmt"
	"net/http"
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
	port := r.FormValue("port")
	id := r.FormValue("id")
	if port == "" || id == "" { // first time: let's execute the script
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
		fmt.Printf("port-id is (%s)\n", s)
		fmt.Println("TODO: output HTML template")
		os.Exit(0)
	}
	res, err := http.Get(fmt.Sprintf("http://localhost:%s?id=%s", port, id))
	if err != nil {
		// TODO: handle better...
		fmt.Printf("Error: %s\n", err)
	}
	res.Header.Write(os.Stdout)
	fmt.Print("\r\n")
	io.Copy(os.Stdout, res.Body)
}
