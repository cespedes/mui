package main

import (
	"io"
	"os"
	"os/exec"
	"fmt"
	"log"
	"flag"
	"strings"
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
	fmt.Print("cgi: debug=", *flagDebug, "; exec=", *flagExec, "; shell=", *flagShell, "\n")

	_, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatal(err)
	}

	// calling as a CGI?
	gateway := os.Getenv("GATEWAY_INTERFACE")
	if !strings.HasPrefix(gateway, "CGI/") {
		// this is not a CGI: let's execute the script:
		cmd := exec.Command(*flagShell, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		os.Exit(0)
	}
}
