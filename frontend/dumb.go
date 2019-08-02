package frontend

import (
	"os"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
)

func init() {
	RegisterFrontEnd(dumb)
}

type Dumb struct {}

func (z Dumb) Available() bool {
	return true
}

func read_letter() (byte,error) {
	oldState, err := terminal.MakeRaw(0)
	if err == nil {
		defer func() {
			terminal.Restore(0, oldState)
			fmt.Println()
		}()
	}
	buf := make([]byte, 1024)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (z Dumb) Question() int {
	for {
		fmt.Print("Are you sure you wany yo proceed? [yn] ")
		c, err := read_letter()
		fmt.Printf("DEBUG: Letter %d\n", c)
		if err != nil {
			return 2
		}
		if c=='y' || c=='Y' {
			return 0
		} else if c=='n' || c=='N' {
			return 1
		} else if c==3 { // ctrl-C
			return 2
		}
	}
}

func (z Dumb) Priority() int {
	return 10
}

var dumb Dumb
