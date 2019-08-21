package frontend

import (
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
)

type dumb struct{}

func init() {
	RegisterFrontEnd(dumb{})
}

func (z dumb) Name() string {
	return "dumb"
}

func (z dumb) Priority() int {
	return 10
}

func (z dumb) Available() bool {
	return true
}

func read_letter_with_echo() (c byte, err error) {
	oldState, err := terminal.MakeRaw(0)
	if err == nil {
		defer func() {
			terminal.Restore(0, oldState)
			if c >= 32 && c <= 126 {
				fmt.Printf("%c", c)
			}
			fmt.Println()
		}()
	}
	buf := make([]byte, 1024)
	_, err = os.Stdin.Read(buf)
	if err != nil {
		return 0, err
	}
	c = buf[0]
	return c, nil
}

func (z dumb) Question() int {
	for {
		fmt.Print("Are you sure you wany yo proceed? [yn] ")
		c, err := read_letter_with_echo()
		fmt.Printf("DEBUG: Letter %d\n", c)
		if err != nil {
			return 2
		}
		if c == 'y' || c == 'Y' {
			return 0
		} else if c == 'n' || c == 'N' {
			return 1
		} else if c == 3 { // ctrl-C
			return 2
		}
	}
}

func (z dumb) Input() string {
	return "text"
}
