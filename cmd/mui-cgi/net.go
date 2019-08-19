package main

import (
	"fmt"
	"log"
	"net"
//	"math/rand"
	"net/http"
	"encoding/json"
)

func net_listen(notes chan string) {
	ln, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		log.Fatal(err)
	}
	addr := ln.Addr()
	tcpaddr, err := net.ResolveTCPAddr(addr.Network(), addr.String())
	if err != nil {
		log.Fatal(err)
	}
	notes <- fmt.Sprintf("net: listening on port %d", tcpaddr.Port)
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		var output struct {
			Stdout []byte
			Stderr []byte
		}
		output.Stdout = make([]byte, 1024)
		output.Stderr = make([]byte, 1024)
		n1, _ := buf_stdout.Read(output.Stdout)
		n2, _ := buf_stdout.Read(output.Stderr)
		output.Stdout = output.Stdout[:n1]
		output.Stderr = output.Stderr[:n2]
		json, _ := json.MarshalIndent(output, "", "\t")
		notes <- fmt.Sprintf("net: sending %d bytes of stdout, %d of stderr", n1, n2)
		w.Write(json)
	})
	err = http.Serve(ln, nil)
	log.Fatal(err)
}