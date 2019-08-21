package frontend

import (
	"fmt"
	"os"
	"strconv"
)

type web struct{}

func init() {
	RegisterFrontEnd(web{})
}

func (z web) Name() string {
	return "web"
}

func (z web) Priority() int {
	return 200
}

func (z web) Available() bool {
	if os.Getenv("MUI-WEB") == "" {
		return false
	}
	return true
}

func (z web) Question() int {
	in := os.NewFile(0x10, "pipe-in")
	out := os.NewFile(0x11, "pipe-out")
	fmt.Fprint(out, "question")
	buf := make([]byte, 1024)
	n, _ := in.Read(buf)
	i, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		return 2
	}
	return i
}

func (z web) Input() string {
	in := os.NewFile(0x10, "pipe-in")
	out := os.NewFile(0x11, "pipe-out")
	fmt.Fprint(out, "input")
	buf := make([]byte, 1024)
	n, _ := in.Read(buf)
	return string(buf[:n])
}
