package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"io/ioutil"
	"strconv"
	"net/http"
	"encoding/json"
)

func http_serve(ln net.Listener, id int, notes chan string, pipe_read, pipe_write *os.File) {
	// No mux, all requests are handled by this function
	err := http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i, _ := strconv.Atoi(r.FormValue("id"))
		if id != i {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
		type param struct {
			Name string  `json:"name"`
			Value string `json:"value"`
		}
		var output struct {
			Stdout []byte  `json:"stdout"`
			Stderr []byte  `json:"stderr"`
			Params []param `json:"params"`
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
