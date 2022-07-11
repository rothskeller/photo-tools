package main

import (
	"bufio"
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/filefmts"
)

//go:embed "page.html"
var pageHTML []byte

var (
	scan      *bufio.Scanner
	files     []string
	handlers  []filefmts.FileFormat
	listener  net.Listener
	index     int
	prevPlace string
)

func main() {
	// Parse arguments and read files.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: assign-places file...")
		os.Exit(2)
	}
	for _, file := range os.Args[1:] {
		handler, err := filefmts.HandlerForName(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
			continue
		}
		if handler == nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s: unsupported file type\n", file)
			continue
		}
		files = append(files, file)
		handlers = append(handlers, handler)
	}
	if len(files) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: no files to act on")
		os.Exit(1)
	}
	listener, _ = net.Listen("tcp", "localhost:0")
	go http.Serve(listener, http.HandlerFunc(handleHTTP))
	time.Sleep(100 * time.Millisecond)
	exec.Command("open", fmt.Sprintf("http://%s/", listener.Addr())).Start()
	scan = bufio.NewScanner(os.Stdin)
	for index = range files {
		prevPlace = handleFile(files[index], handlers[index], prevPlace)
	}
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == fmt.Sprintf("/%s", filepath.Base(files[index])) {
		http.ServeFile(w, r, files[index])
		return
	}
	if r.URL.Path == "/url" {
		fmt.Fprintf(w, "/%s", filepath.Base(files[index]))
		return
	}
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html")
		w.Write(pageHTML)
		return
	}
	http.Error(w, "404 Not Found", http.StatusNotFound)
}

func handleFile(fname string, handler filefmts.FileFormat, defPlace string) string {
	var (
		places    []metadata.HierValue
		currPlace string
		err       error
		in        string
		uri       url.URL
	)
	uri.Scheme = "file"
	uri.Path, _ = filepath.Abs(fname)
	gps := handler.Provider().GPS()
RESTART:
	fmt.Printf("\x1B[2J%s (%f, %f)\n", fname, gps.Latitude(), gps.Longitude())
	if places = handler.Provider().Places(); len(places) != 0 {
		currPlace = "/" + places[0].String()
		defPlace = "/" + places[0].String()
	}
	if defPlace == "" {
		fmt.Print("? ")
	} else {
		fmt.Printf("%s / ? ", defPlace)
	}
	if !scan.Scan() {
		os.Exit(1)
	}
	in = strings.TrimSpace(scan.Text())
	if in == "" {
		if defPlace == currPlace || defPlace == "" {
			return defPlace
		}
		in = defPlace
	}
	if strings.HasPrefix(in, "/") {
		currPlace = in
	} else {
		currPlace = path.Clean(path.Join(defPlace, in))
	}
	if len(places) == 0 {
		places = append(places, metadata.HierValue{})
	}
	places[0], err = metadata.ParseHierValue(currPlace[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		goto RESTART
	}
	if err := handler.Provider().SetPlaces(places); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
		os.Exit(1)
	}
	if err := filefmts.Save(handler, fname); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	return currPlace
}
