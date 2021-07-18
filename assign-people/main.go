package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/rothskeller/photo-tools/metadata/filefmts"
	"github.com/webview/webview"
)

var (
	abbrevs       map[string]string
	abbrevFor     map[string]string
	personList    []string
	longestAbbrev int
	longestPerson int
	scan          *bufio.Scanner
	viewer        webview.WebView
)

func main() {
	var (
		files    []string
		handlers []filefmts.FileFormat
	)
	// Parse arguments and read files.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: assign-people file...")
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
	// Generate abbreviations.
	abbrevs = make(map[string]string)
	abbrevFor = make(map[string]string)
	for _, handler := range handlers {
		for _, person := range handler.Provider().People() {
			assignAbbrev(person)
		}
	}
	scan = bufio.NewScanner(os.Stdin)
	viewer = webview.New(true)
	defer viewer.Destroy()
	viewer.SetTitle("assign-people")
	viewer.SetSize(800, 600, webview.HintNone)
	viewer.Navigate("about:blank")
	go func() {
		for i := range files {
			handleFile(files[i], handlers[i])
		}
		viewer.Dispatch(func() { viewer.Terminate() })
	}()
	viewer.Run()
}

func handleFile(fname string, handler filefmts.FileFormat) {
	var (
		pmap   = make(map[string]bool)
		in     string
		abbrs  []string
		remove bool
		uri    url.URL
	)
	for _, person := range handler.Provider().People() {
		pmap[person] = true
	}
	uri.Scheme = "file"
	uri.Path, _ = filepath.Abs(fname)
	viewer.Dispatch(func() { viewer.Navigate(uri.String()) })
	fmt.Printf("\x1B[2J%s\n", fname)
	for _, person := range personList {
		if pmap[person] {
			fmt.Printf("  * %-*s %s\n", longestAbbrev, abbrevFor[person], person)
		} else {
			fmt.Printf("    %-*s %s\n", longestAbbrev, abbrevFor[person], person)
		}
	}
	fmt.Print("? ")
	if !scan.Scan() {
		os.Exit(1)
	}
	in = strings.TrimSpace(scan.Text())
	if in == "" {
		return
	}
	if in == "-ALL" {
		pmap = make(map[string]bool)
		goto SAVE
	}
	if in[0] == '-' {
		remove = true
		in = strings.TrimSpace(in[1:])
		if in == "" {
			return
		}
	} else if in[0] == '+' {
		in = strings.TrimSpace(in[1:])
		if in == "" {
			return
		}
	} else {
		pmap = make(map[string]bool)
	}
	abbrs = strings.Fields(in)
	for _, abbr := range abbrs {
		if person, ok := abbrevs[abbr]; ok {
			if remove {
				delete(pmap, person)
			} else {
				pmap[person] = true
			}
		} else if remove {
			fmt.Printf("ERROR: can't remove unknown person %q\n", abbr)
		} else {
			fmt.Printf("Who is %s? ", abbr)
			if !scan.Scan() {
				os.Exit(1)
			}
			in2 := strings.TrimSpace(scan.Text())
			if in2 != "" {
				pmap[in2] = true
				addAbbrev(abbr, in2)
			}
		}
	}
SAVE:
	var plist []string
	for _, person := range personList {
		if pmap[person] {
			plist = append(plist, person)
		}
	}
	if err := handler.Provider().SetPeople(plist); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", fname, err)
		os.Exit(1)
	}
	if err := filefmts.Save(handler, fname); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

func assignAbbrev(person string) {
	if _, ok := abbrevFor[person]; ok {
		return
	}
	var abbrev = strings.ToLower(strings.Map(func(r rune) rune {
		if r < 'A' || r > 'Z' {
			return -1
		}
		return r
	}, person))
	if _, ok := abbrevs[abbrev]; ok {
		var seq = 2
		for {
			var ab2 = fmt.Sprintf("%s%d", abbrev, seq)
			if _, ok := abbrevs[ab2]; ok {
				seq++
			} else {
				abbrev = ab2
				break
			}
		}
	}
	addAbbrev(abbrev, person)
}

func addAbbrev(abbrev, person string) {
	abbrevs[abbrev] = person
	abbrevFor[person] = abbrev
	personList = append(personList, person)
	sort.Strings(personList)
	if len(abbrev) > longestAbbrev {
		longestAbbrev = len(abbrev)
	}
	if len(person) > longestPerson {
		longestPerson = len(person)
	}
}
