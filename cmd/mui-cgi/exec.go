package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
//	"time"
//	"syscall"
//	"math/rand"
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

func executeScript(shell string, args []string, notes chan string) {
	var err error
	cmd := exec.Command(shell, args...)
	//	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = &buf_stdout
	cmd.Stderr = &buf_stderr
	r1, _, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	_, w2, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	// TODO: handle pipes from the server side
	cmd.ExtraFiles = make([]*os.File, 15)
	cmd.ExtraFiles[13] = r1
	cmd.ExtraFiles[14] = w2

	notes <- fmt.Sprintf("exec: starting command: %v", cmd)
	if err = cmd.Run(); err != nil {
		panic(err)
	}
	notes <- "exec: command finished."
	/*
		ch := make(chan error)
		go func() {
			ch <- cmd.Run()
		}()
		select {
			case err := <- ch:
				fmt.Printf("Error: %v\n", err)
		}
		close(ch)
	*/
}
