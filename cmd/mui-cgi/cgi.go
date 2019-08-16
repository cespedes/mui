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
