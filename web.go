package mui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const (
	fd_in  = 0x10
	fd_out = 0x11
)

type web struct{
	in *os.File
	out *os.File
	buf_in *bufio.Reader
}

func init() {
	w := new(web)
	if os.Getenv("MUI_WEB") != "" {
		w.in = os.NewFile(fd_in, "pipe-in")
		w.out = os.NewFile(fd_out, "pipe-out")
		w.buf_in = bufio.NewReader(w.in)
	}
	RegisterFrontEnd(w)
}

func (w web) Name() string {
	return "web"
}

func (w web) Priority() int {
	return 200
}

var result struct {
	ok bool
}

func (w web) send(name, value string) error {
	fmt.Fprintf(os.Stderr, "DEBUG: sending (%s,%s) to pipe\n", name, value)
	_, err := fmt.Fprintf(w.out, "%s\000%s\000", name, value)
	return err
}

func (w web) recv() (value string, err error) {
	value, err = w.buf_in.ReadString('\000')
	fmt.Fprintf(os.Stderr, "DEBUG: mui.web.recv() = (%s,%v)\n", value, err)
	return
}

func (w web) Available() bool {
	if os.Getenv("MUI_WEB") == "" {
		return false
	}
	err := w.send("", "")
	if err != nil {
		return false
	}
	_, err = w.recv()
	if err != nil {
		return false
	}
	return true
}

func (w web) Question() int {
	w.send("type", "question")
	resp, err := w.recv()
	if err != nil {
		return 2
	}
	i, err := strconv.Atoi(resp)
	if err != nil {
		return 2
	}
	return i
}

func (w web) Input() string {
	return ""
}
