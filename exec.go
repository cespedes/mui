package main

import (
	"fmt"
	"log"
//	"time"
	"os"
	"os/exec"
	"net"
	"flag"
	"syscall"
)

func execUsage() {
	fmt.Fprintln(os.Stderr, `Usage: mui exec [-shell interpreter] script [args]

Use "mui help exec" for more information`)
}

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
}

func executeScript(shell string, args []string) {
	if len(args) < 1 {
		execUsage()
		os.Exit(1)
	}
	path, err := exec.LookPath(args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("path = %s\n", path)
	net_listen()
	cmd := exec.Command(shell, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
//	r, w, err := os.Pipe()
//	cmd.Stdout = w
	fmt.Printf("Executing: %+v\n", cmd)
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}
	for {
		out := make([]byte, 1024)
		n, err := stdout.Read(out)
		if err != nil {
			break
		}
		fmt.Printf("%s", out[:n])
	}
	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}

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

func runExec(args []string) {
	f := flag.NewFlagSet("mui exec", flag.ContinueOnError)
	f.Usage = execUsage
	pshell := f.String("shell", "/bin/sh", "Interpreter to use")
	f.Parse(args)
	args = f.Args()
	fmt.Printf("shell = %v\n", *pshell)
	executeScript(*pshell, args)
}
