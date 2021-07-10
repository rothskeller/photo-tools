package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/photo-tools/metadata/filefmts"
)

func main() {
	var file string
	defer func() {
		if p := recover(); p != nil {
			println(file)
			panic(p)
		}
	}()
	// os.Args = []string{"", "/Users/stever/Timeline/1998-05-03-HMC-Reunion/2-MVC-018F.jpg"}
	for _, file = range os.Args[1:] {
		fh, err := os.Open(file)
		if err != nil {
			fmt.Printf("%s: %s\n", file, err)
			continue
		}
		handler, err := filefmts.HandlerFor(fh)
		fh.Close()
		if err != nil {
			fmt.Printf("%s: %s\n", file, err)
		} else if handler == nil {
			fmt.Printf("%s: no handler\n", file)
		}
	}
}
