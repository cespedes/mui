package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type Buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}

var buf_stdout, buf_stderr Buffer

func executeScript(shell string, args []string, notes chan string, pipe_read, pipe_write *os.File) {
	var err error
	cmd := exec.Command(shell, args...)
	cmd.Stdout = &buf_stdout
	cmd.Stderr = &buf_stderr
	cmd.Env = append(os.Environ(), "MUI_WEB=1")

	cmd.ExtraFiles = make([]*os.File, 15)
	cmd.ExtraFiles[13] = pipe_read
	cmd.ExtraFiles[14] = pipe_write

	notes <- fmt.Sprintf("exec: starting command: %v", cmd)
	err = cmd.Run()
	notes <- fmt.Sprintf("exec: command finished with error %v", err)
}
