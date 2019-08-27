package main

import (
	"io"
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"os/exec"
	"syscall"
	"html/template"
//	"encondig/json"
)

const html_template = `
<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>{{.Name}}</title>
    <script src="https://code.jquery.com/jquery-3.4.1.min.js"></script>
    <script>
      $(function() {
        setInterval(function() {
          $.getJSON("{{.Url}}?port={{.Port}}&id={{.Id}}", function(data) {
            stdout = atob(data.stdout);
            stderr = atob(data.stderr);
            if (stdout.length > 0) {
              $(".output").append(stdout);
            }
            if (stderr.length > 0) {
              $(".output").append($('<span class="stderr">').html(stderr));
            }
          });
        }, 500);
      });
    </script>
    <style type="text/css">
      .output {
        font-family: monospace;
        background: #ddd;
        white-space: pre-wrap;
        word-break: break-all
      }
      .output .stderr {
        font-weight: bold;
        color: red;
      }
    </style>
  </head>
  <body>
    <h1>Executing {{.Name}}...</h1>
    <div class="output"></div>
  </body>
</html>
`

func cgi_handle(path string, args []string) {
	r, err := cgi.Request()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	// id is of form <port>-<cookie>
	port := r.FormValue("port")
	id := r.FormValue("id")
	if port == "" || id == "" { // first time: let's execute the script
		fmt.Println("Content-Type: text/html")
		fmt.Println()
		newargs := []string{"-exec", "-shell", path}
		newargs = append(newargs, args...)
		cmd := exec.Command("mui-cgi", newargs...)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = cmd.Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		var port, id int
		fmt.Fscanln(stdout, &port, &id)
		if port==0 || id==0 {
			fmt.Printf("ERROR: port=%d, id=%d\n", port, id)
		}
		t, err := template.New("webpage").Parse(html_template)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		data := struct {
			Name string
			Url string
			Port int
			Id int
		}{
			Name: args[0],
			Url: os.Getenv("REQUEST_URI"),
			Port: port,
			Id: id,
		}
		err = t.Execute(os.Stdout, data)
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(0)
	}
	// TODO: handle mui requests and responses
	res, err := http.Get(fmt.Sprintf("http://localhost:%s?id=%s", port, id))
	if err != nil {
		fmt.Println("Status: 412 Precondition Failed")
		fmt.Println()
		fmt.Printf("Error: could not connect to local port %s\n", port)
		return
	}
	res.Header.Write(os.Stdout)
	fmt.Print("\r\n")
	io.Copy(os.Stdout, res.Body)
}
