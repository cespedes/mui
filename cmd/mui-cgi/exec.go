package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"
	//	"syscall"
	"encoding/json"
)

func net_listen() {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	addr := ln.Addr()
	tcpaddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Listening on port %d\n", tcpaddr.Port)
	go func() {
		var output struct {
			Stdout []byte
			Stderr []byte
		}
		output.Stdout = make([]byte, 1024)
		output.Stderr = make([]byte, 1024)
		for {
			conn, _ := ln.Accept()
			output.Stdout = output.Stdout[:cap(output.Stdout)]
			output.Stderr = output.Stderr[:cap(output.Stderr)]
			n1, _ := buf_stdout.Read(output.Stdout)
			n2, _ := buf_stdout.Read(output.Stderr)
			output.Stdout = output.Stdout[:n1]
			output.Stderr = output.Stderr[:n2]
			json, _ := json.MarshalIndent(output, "", "\t")
			log.Printf("sending %d bytes of stdout, %d of stderr\n", n1, n2)
			conn.Write(json)
			conn.Close()
		}
	}()
}

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

func executeScript(shell string, args []string) {
	var err error
	cmd := exec.Command(shell, args...)
	//	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = &buf_stdout
	cmd.Stderr = &buf_stderr
	r1, w1, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	fmt.Println("os.Pipe() = ", r1.Fd(), w1.Fd())
	r2, w2, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	fmt.Println("os.Pipe() = ", r2.Fd(), w2.Fd())
	cmd.ExtraFiles = []*os.File{nil, nil, w1, r2}

	fmt.Printf("Executing: %+v\n", cmd)
	if err = cmd.Run(); err != nil {
		panic(err)
	}
	fmt.Println("Finished.")
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

func runExec(path string, args []string) {
	net_listen()
	executeScript(path, args)
	for {
		log.Printf("tick\n")
		time.Sleep(1 * time.Second)
	}
}
