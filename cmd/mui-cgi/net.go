package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"bufio"
	"io/ioutil"
	"strconv"
	"net/http"
	"encoding/json"
)

func http_serve(ln net.Listener, id int, notes chan string, pipe_read, pipe_write *os.File) {
	// No mux, all requests are handled by this function
	// TODO read from pipe_read and send data to http client
	// TODO read data from http client and send it to pipe_write
	buf_read := bufio.NewReader(pipe_read)
	type param struct {
		Name string  `json:"name"`
		Value string `json:"value"`
	}
	t := make(chan param)
	go func() {
		for {
			var p param
			var err error
			p.Name, err = buf_read.ReadString('\000')
			if err == nil {
				p.Value, err = buf_read.ReadString('\000')
			}
			if err != nil {
				notes <- fmt.Sprintf("error reading from pipe: %s", err)
				break
			}
			p.Name = p.Name[:len(p.Name)-1]
			p.Value = p.Value[:len(p.Value)-1]
			if len(p.Name) == 0 { // this is a "ping"
				fmt.Fprint(pipe_write, "\000")
			} else {
				debug.Printf("mui-cgi.net: Sent to channel: %v\n", p)
				t <- p
			}
		}
	}()
	err := http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i, _ := strconv.Atoi(r.FormValue("id"))
		if id != i {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		var output struct {
			Stdout []byte  `json:"stdout"`
			Stderr []byte  `json:"stderr"`
			Params []param `json:"params"`
		}
		debug.Println("mui-cgi.net: starting")
		output.Params = make([]param, 0)
		L:
		for {
			select {
				case p := <-t:
					debug.Printf("mui-cgi.net: Received from channel: %v\n", p)
					output.Params = append(output.Params, p)
				default:
					break L
			}
		}
		output.Stdout, _ = ioutil.ReadAll(&buf_stdout)
		output.Stderr, _ = ioutil.ReadAll(&buf_stderr)
//		output.Stderr = []byte(fmt.Sprintf("debug = %#v\n", debug))
		json, _ := json.MarshalIndent(output, "", "\t")
		notes <- fmt.Sprintf("net: sending %d bytes of stdout, %d of stderr", len(output.Stdout), len(output.Stderr))
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}))
	log.Fatal(err)
}
