// identical-photos examines a set of photos and prints which ones of them
// contain the same pixels.  (The metadata may be different.)
package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"

	_ "golang.org/x/image/tiff"
)

func main() {
	var pmap = make(map[string][]string)
	for _, file := range os.Args[1:] {
		if fh, err := os.Open(file); err == nil {
			if img, _, err := image.Decode(fh); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %s: %s\n", file, err)
			} else {
				b := img.Bounds()
				h := md5.New()
				buf := make([]byte, 16)
				for x := b.Min.X; x < b.Max.X; x++ {
					for y := b.Min.Y; y < b.Max.Y; y++ {
						r, g, b, a := img.At(x, y).RGBA()
						binary.BigEndian.PutUint32(buf[0:], r)
						binary.BigEndian.PutUint32(buf[4:], g)
						binary.BigEndian.PutUint32(buf[8:], b)
						binary.BigEndian.PutUint32(buf[12:], a)
						h.Write(buf)
					}
				}
				sum := hex.EncodeToString(h.Sum(nil))
				pmap[sum] = append(pmap[sum], file)
			}
			fh.Close()
		}
	}
	for _, files := range pmap {
		if len(files) > 1 {
			fmt.Print("identical:")
			for _, file := range files {
				fmt.Printf(" %s", file)
			}
			fmt.Println()
		}
	}
	for _, files := range pmap {
		if len(files) == 1 {
			fmt.Printf("unique: %s\n", files[0])
		}
	}
}
