package main

import (
	"io"
	"os"
	"fmt"
	"flag"

	"./frontend"
)

type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'mui help' output.
	Short string

	// Long is the long message shown in the 'mui help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

var (
	flagQuestion = flag.Bool("question", false, "Display question dialog")
	flagEntry    = flag.Bool("entry",    false, "Display text entry dialog")
)

var Mui = &Command{
	UsageLine: "mui",
	Long: `Mui is a tool to display graphical dialog boxes.

Usage:

	mui <command> [arguments]

The commands are:

	question display question dialog
	calendar display calendar dialog
	entry    display text entry dialog
	error    display error dialog
	info     display info dialog
	scale    display scale dialog
	progress display progress indication dialog
	password display password dialog
	list     display list dialog

Use "mui help <command>" for more information about a command.`,
//	Long: "Mui is a tool to display terminal, X11 (Gtk) or web graphical dialog boxes from shell scripts.",
}

func printUsage(w io.Writer, cmd *Command) {
	fmt.Fprintln(w, cmd.Long)
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage(os.Stderr, Mui)
		os.Exit(1)
	}
	cmdname := args[0]

	if cmdname == "question" {
		frontend.Question()
	} else {
		fmt.Fprintf(os.Stderr, "mui %s: unknown command\nRun 'mui help' for usage.\n", cmdname)
	}
}
