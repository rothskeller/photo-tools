package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var in = bufio.NewScanner(os.Stdin)
	var lat, long float64
	for in.Scan() {
		var line = in.Text()
		var cont = false
		var bits = 0
		var i32 int32
		var first = true
		var islat = true
		for idx := 0; idx < len(line); idx++ {
			var by = line[idx] - 63
			if by&0x20 != 0 {
				by &^= 0x20
				cont = true
			} else {
				cont = false
			}
			i32 |= int32(by) << bits
			bits += 5
			if !cont {
				var neg = false
				if i32&1 == 1 {
					neg = true
				}
				i32 >>= 1
				var val = float64(i32) / 100000
				if neg {
					val = -val
				}
				if first {
					if islat {
						lat = val
						islat = false
					} else {
						long = val
						first = false
						islat = true
						fmt.Printf("%f,%f\n", lat, long)
					}
				} else {
					if islat {
						lat += val
						islat = false
					} else {
						long += val
						islat = true
						fmt.Printf("%f,%f\n", lat, long)
					}
				}
				i32 = 0
				bits = 0
			}
		}
		if cont {
			fmt.Fprintf(os.Stderr, "ERROR: incomplete encoding\n")
		}
		fmt.Println()
	}
}

/*
~{ugV
126 123 117 103 86
63 60 54 40 23
31 28 22 8 23
11111 11100 10110 00100 10111

jS|CiBdAxy

*/
