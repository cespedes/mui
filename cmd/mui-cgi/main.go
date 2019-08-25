package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

var (
	flagDebug = flag.Bool("debug", false, "Show debugging information")
	flagExec  = flag.Bool("exec", false, "Internal use only")
	flagShell = flag.String("shell", "/bin/sh", "Interpreter to use")
)

//        Short: "execute script as a CGI client",
//        Long: `Usage: mui cgi [-shell interpreter] script [args]

//Execute a script with an interpreter (default /bin/sh),
//sending its output to a web browser as a CGI binary.

func printUsage(w io.Writer) {
	fmt.Fprint(w, `mui-cgi executes a script as a CGI client.

It is meant to be inserted as the first line of a shell script as:

#!/usr/bin/mui-cgi [-shell interpreter] [arguments]

This will execute the script with am interpreter (default /bin/sh).

If the script is run from a web browser, as a CGI client, it will generate
a web page which shows the script's standard output as it is generated.

It also shows dialogs created with "mui" as HTML pop-ups.
`)
}

func init() {
	flag.Usage = func() {
		printUsage(flag.CommandLine.Output())
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage(os.Stderr)
		os.Exit(2)
	}

	_, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatal(err)
	}

	path, err := exec.LookPath(*flagShell)
	if err != nil {
		log.Fatal(err)
	}

	if *flagExec {
		// Create TCP listener:
		ln, err := net.Listen("tcp", "127.0.0.1:")
		if err != nil {
			log.Fatal(err)
		}
		addr := ln.Addr()
		tcpaddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
		if err != nil {
			log.Fatal(err)
		}
		rand.Seed(time.Now().UTC().UnixNano())
		id := rand.Int()
		fmt.Printf("%d-%d\n", tcpaddr.Port, id)

		// 2 pipe to communicate with script:
		r1, w1, err := os.Pipe()
		if err != nil {
			panic(err)
		}
		r2, w2, err := os.Pipe()
		if err != nil {
			panic(err)
		}

		notes := make(chan string)
		go http_serve(ln, id, notes, r1, w2)
		go executeScript(path, args, notes, r2, w1)
		for {
			s := <-notes
			log.Printf("note: %s\n", s)
		}
		return
	}

	// calling as a CGI?
	if gi := os.Getenv("GATEWAY_INTERFACE"); strings.HasPrefix(gi, "CGI/") {
		cgi_handle(path, args)
		os.Exit(0)
	}

	// this is not a CGI: let's execute the script:
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
