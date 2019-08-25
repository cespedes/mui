package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"net/http"
	"encoding/json"
)

func http_serve(ln net.Listener, id int, notes chan string, pipe_read, pipe_write *os.File) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		i, _ := strconv.Atoi(r.FormValue("id"))
		if id != i {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}
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
		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	})
	err := http.Serve(ln, nil)
	log.Fatal(err)
}
