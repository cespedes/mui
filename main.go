package main

import (
	"io"
	"os"
	"fmt"
	"flag"
	"strings"

	"./frontend"
)

type Command struct {
	Name string

	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(args []string)

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

var cmdCGI = Command{
	Name:  "cgi",
	Run:   runCGI,
	Short: "execute script as a CGI client",
	Long: `Usage: mui cgi [-shell interpreter] script [args]

Execute a script with an interpreter (default /bin/sh),
sending its output to a web browser as a CGI binary.
`,
}

var cmdExec = Command{
	Name:  "exec",
	Run:   runExec,
	Short: "execute a script - internal use only",
}

var cmdQuestion = Command{
	Name:  "question",
	Run:   frontend.Question,
	Short: "display question dialog",
	Long: `Usage: mui question

Display a question with two possible answers: Yes or No.

IT returns exit code zero with "Yes" and nonzero with "No".
`,
}

var cmdCalendar = Command{
	Name:  "calendar",
	Short: "display calendar dialog",
}

var cmdInput = Command{
	Name:  "input",
	Run:   frontend.Input,
	Short: "display text input dialog",
}

var cmdInfo = Command{
	Name:  "info",
	Short: "display info dialog",
}

/*
	input    display text input dialog
	error    display error dialog
	scale    display scale dialog
	progress display progress indication dialog
	password display password dialog
	list     display list dialog
	select
	radio
	checkbox
*/

var commands = []Command{
	cmdCGI,
	cmdExec,
	cmdQuestion,
//	cmdCalendar,
	cmdInput,
}

func printUsage(w io.Writer, cmd *Command) {
	if cmd == nil {
		fmt.Fprint(w, `Mui is a tool to display graphical dialog boxes.

Usage:

        mui <command> [arguments]

The commands are:

`)

		for _, c := range commands {
			fmt.Fprintf(w, "\t%-10s %s\n", c.Name, c.Short)
		}
/*
        input    display text input dialog
        error    display error dialog
        info     display info dialog
        scale    display scale dialog
        progress display progress indication dialog
        password display password dialog
        list     display list dialog
*/

		fmt.Fprintln(w, `
Use "mui help <command>" for more information about a command.`)
		return
	} else {
		fmt.Fprint(w, cmd.Long)
	}
}

var (
	flagDebug = flag.Bool("debug", false, "Show debugging information")
)

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		printUsage(os.Stderr, nil)
		os.Exit(2)
	}
	cmdname := args[0]

	for _, c := range commands {
		if cmdname == c.Name {
			c.Run(args[1:])
			return
		}
	}
	if cmdname == "help" {
		args = args[1:]
		if len(args) == 0 {
			printUsage(os.Stdout, nil)
			return
		}
		if len(args) == 1 {
			for _, c := range commands {
				if args[0] == c.Name {
					printUsage(os.Stdout, &c)
					return
				}
			}
		}
		fmt.Fprintf(os.Stderr, "mui help %s: unknown help topic. Run \"go help\".\n",
			strings.Join(args, " "))
		os.Exit(2)
	} else {
		fmt.Fprintf(os.Stderr, "mui %s: unknown command\nRun 'mui help' for usage.\n", cmdname)
		os.Exit(2)
	}
}
