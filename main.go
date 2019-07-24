package main

import (
	"os"
	"fmt"
	"flag"

	"./frontend"
)

var (
	flagQuestion = flag.Bool("question", false, "Display question dialog")
	flagEntry    = flag.Bool("entry",    false, "Display text entry dialog")
)

func main() {
	flag.Parse()

	if !*flagQuestion && !*flagEntry {
		fmt.Fprintf(os.Stderr, "You must specify a dialog type. See 'shell-ui -help' for details\n")
		os.Exit(1)
	}
	if *flagQuestion {
		frontend.Question()
	}
}
