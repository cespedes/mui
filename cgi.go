package main

import (
	"os"
//	"os/exec"
//	"bytes"
//	"net"
	"fmt"
	"flag"
	"strings"
)

func cgiUsage() {
	fmt.Fprintln(os.Stderr, `Usage: mui cgi [-shell interpreter] script [args]

Use "mui help cgi" for more information`)
}

/*
func executeScript(shell string, args []string) {
	if len(args) < 1 {
		cgiUsage()
		os.Exit(1)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("path = %s\n", path)
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	cmd := exec.Command(shell, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	fmt.Printf("Command to execute: %v\n", cmd)
	addr := ln.Addr()
	tcpaddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Listening on port %d\n", tcpaddr.Port)
}
*/

func runCGI(args []string) {
	f := flag.NewFlagSet("mui cgi", flag.ContinueOnError)
	f.Usage = cgiUsage
	pshell := f.String("shell", "/bin/sh", "Interpreter to use")
	f.Parse(args)
	args = f.Args()
	gateway := os.Getenv("GATEWAY_INTERFACE")
	if !strings.HasPrefix(gateway, "CGI/") {
		fmt.Fprintln(os.Stderr, "Error: No GATEWAY_INTERFACE defined.  Is this a CGI?")
		os.Exit(1)
	}
	fmt.Printf("shell = %v\n", *pshell)
//	executeScript(*pshell, args)
}
