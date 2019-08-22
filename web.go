package mui

import (
	"fmt"
	"os"
	"strconv"
	"encoding/json"
)

const (
	fd_in  = 0x10
	fd_out = 0x11
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

var result struct {
	ok bool
}

func (z web) Available() bool {
	if os.Getenv("MUI_WEB") == "" {
		return false
	}
	in := os.NewFile(fd_in, "pipe-in")
	out := os.NewFile(fd_out, "pipe-out")
	_, err := fmt.Fprint(out, "{}")
	if err != nil {
		return false
	}
	buf := make([]byte, 1024)
	n, err := in.Read(buf)
	if err != nil {
		return false
	}
	err = json.Unmarshal(buf[:n], &result)
	return true
}

func (z web) Question() int {
	in := os.NewFile(fd_in, "pipe-in")
	out := os.NewFile(fd_out, "pipe-out")
	_, err := fmt.Fprint(out, "question")
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
