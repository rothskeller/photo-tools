package main

import (
	"fmt"
	"os"

	"github.com/rothskeller/photo-tools/filefmt"
)

func main() {
	for _, file := range os.Args[1:] {
		handler := filefmt.HandlerFor(file)
		if handler == nil {
			fmt.Printf("%s: no handler\n", file)
			continue
		}
		handler.ReadMetadata()
		if p := handler.Problems(); len(p) != 0 {
			fmt.Printf("%s: ERROR: %s\n", file, p[0])
			continue
		}
		var dirty string
		if exif := handler.EXIF(); exif != nil {
			if exif.Dirty() {
				dirty = "EXIF "
			}
		}
		if iptc := handler.IPTC(); iptc != nil {
			if iptc.Dirty() {
				dirty += "IPTC "
			}
		}
		if xmp := handler.XMP(true); xmp != nil {
			if xmp.Dirty() {
				dirty += "XMP "
			}
		}
		if dirty != "" {
			fmt.Printf("%s: %sdirty after read\n", file, dirty)
		}
	}
}
