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
			t <- p
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
		L:
		for {
			select {
				case p := <-t:
					output.Params = append(output.Params, p)
				default:
					break L
			}
		}
		output.Stdout, _ = ioutil.ReadAll(&buf_stdout)
		output.Stderr, _ = ioutil.ReadAll(&buf_stderr)
		output.Params = make([]param, 0)
		json, _ := json.MarshalIndent(output, "", "\t")
		notes <- fmt.Sprintf("net: sending %d bytes of stdout, %d of stderr", len(output.Stdout), len(output.Stderr))
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}))
	log.Fatal(err)
}
