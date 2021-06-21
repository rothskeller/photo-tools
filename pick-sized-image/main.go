package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"strconv"
)

var widths = []int{320, 640, 1280, 2560, 5120, 10240}
var heights = []int{240, 480, 960, 1920, 3860}

func main() {
	cgi.Serve(http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	var img = r.FormValue("img")
	var width, _ = strconv.Atoi(r.FormValue("width"))
	var height, _ = strconv.Atoi(r.FormValue("height"))
	for _, pw := range widths {
		if pw < width {
			continue
		}
		fn := fmt.Sprintf("%s.%dw.jpg", img[:len(img)-4], pw)
		if _, err := os.Stat(fn); err == nil {
			img = fn
			goto FOUND
		}
	}
	for _, ph := range heights {
		if ph < height {
			continue
		}
		fn := fmt.Sprintf("%s.%dh.jpg", img[:len(img)-4], ph)
		if _, err := os.Stat(fn); err == nil {
			img = fn
			goto FOUND
		}
	}
FOUND:
	http.Redirect(w, r, img, http.StatusTemporaryRedirect)
}
