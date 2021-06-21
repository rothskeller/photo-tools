// Package webdialogs provides the framework for a program which opens a web
// browser window and then asks a series of questions in a series of web forms.
package webdialogs

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// A DialogHandler is a standard Go web request handler, except that it returns
// the handler to be used for the next request.
type DialogHandler func(w http.ResponseWriter, r *http.Request) DialogHandler

// handler is the handler to be used for the next incoming request.
var handler DialogHandler

// Main is designed to be called by the main function of the program using this
// library.  It starts a server and opens a web browser talking to it.  The
// supplied firstHandler is used for the first request to that server.  When any
// dialog in the sequence returns a nil handler, the program exits.
func Main(firstHandler DialogHandler) {
	handler = firstHandler
	listener, _ := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	href := "http://" + listener.Addr().String() + "/"
	exec.Command("open", href).Start()
	http.Serve(listener, http.HandlerFunc(nextHandler))
}

func nextHandler(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	handler = handler(w, r)
	if handler == nil {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, `Dialog sequence is concluded.
Close this browser window.`)
		go func() {
			time.Sleep(time.Second)
			os.Exit(0)
		}()
	}
}
