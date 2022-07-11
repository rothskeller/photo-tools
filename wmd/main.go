package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/rothskeller/photo-tools/metadata"
	"github.com/rothskeller/photo-tools/metadata/filefmts"
)

//go:embed "dist/*"
var dist embed.FS

//go:embed "page.html"
var pageHTML []byte

var (
	files    []string
	handlers []filefmts.FileFormat
	listener net.Listener
)

func main() {
	var seenBase = make(map[string]bool)
	// Parse arguments and read files.
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: wmd file...")
		os.Exit(2)
	}
	for _, file := range os.Args[1:] {
		base := filepath.Base(file)
		if seenBase[base] {
			fmt.Fprintf(os.Stderr, "ERROR: multiple files named %q\n", base)
			os.Exit(2)
		}
		seenBase[base] = true
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
	time.Sleep(365 * 24 * time.Hour)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(pageHTML)
		return
	}
	if r.URL.Path == "/" {
		r.URL.Path = "/index.html"
	}
	if f, err := dist.Open(path.Join("dist", r.URL.Path)); err == nil {
		http.ServeContent(w, r, r.URL.Path, time.Time{}, rdsk{f})
		return
	}
	if r.URL.Path == "/metadata.json" {
		sendMetadata(w)
		return
	}
	var index = -1
	for i, fname := range files {
		base := filepath.Base(fname)
		if r.URL.Path == "/"+base {
			index = i
			break
		}
	}
	if index < 0 {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, files[index])
		return
	}
	applyChanges(w, r, index)
}

type imgmd struct {
	Filename string
	Artist   string
	Caption  string
	DateTime string
	GPS      string
	Groups   []string
	Keywords []string
	Location string
	People   []string
	Places   []string
	Title    string
	Topics   []string
}

type hier struct {
	Name     string
	Children []*hier
}

func sendMetadata(w http.ResponseWriter) {
	var places = []*hier{}
	var topics = []*hier{}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, `{"images":[`)
	enc := json.NewEncoder(w)
	for index, fname := range files {
		if index != 0 {
			fmt.Fprint(w, ",")
		}
		md := metadataForImage(filepath.Base(fname), handlers[index].Provider())
		enc.Encode(&md)
		for _, place := range handlers[index].Provider().Places() {
			places = addToHierarchy(places, place)
		}
		for _, topic := range handlers[index].Provider().Topics() {
			topics = addToHierarchy(topics, topic)
		}
	}
	fmt.Fprint(w, `],"placeHierarchy":`)
	enc.Encode(places)
	fmt.Fprint(w, `,"topicHierarchy":`)
	enc.Encode(topics)
	fmt.Fprint(w, `}`)
}

func addToHierarchy(h []*hier, v metadata.HierValue) []*hier {
	if len(v) == 0 {
		return h
	}
	var c = v[0]
	for i, n := range h {
		if c == n.Name {
			h[i].Children = addToHierarchy(n.Children, v[1:])
			return h
		}
	}
	var n hier
	n.Name = c
	n.Children = addToHierarchy(n.Children, v[1:])
	h = append(h, &n)
	sort.Slice(h, func(i, j int) bool {
		return h[i].Name < h[j].Name
	})
	return h
}

type errors struct {
	Errors []string
}

func applyChanges(w http.ResponseWriter, r *http.Request, index int) {
	var errs errors
	handler := handlers[index]
	provider := handler.Provider()
	r.ParseMultipartForm(1048576)
	var save bool
	if r.Form["artist"] != nil {
		save = true
		if err := provider.SetCreator(r.FormValue("artist")); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	if r.Form["caption"] != nil {
		save = true
		if err := provider.SetCaption(r.FormValue("caption")); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	if r.Form["title"] != nil {
		save = true
		if err := provider.SetTitle(r.FormValue("title")); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	if r.Form["gps"] != nil {
		var gps metadata.GPSCoords
		if err := gps.Parse(r.FormValue("gps")); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		} else if err := provider.SetGPS(gps); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		} else {
			save = true
		}
	}
	if r.Form["places"] != nil {
		save = true
		var hvs []metadata.HierValue
		if len(r.Form["places"]) > 1 || r.Form["places"][0] != "" {
			for _, p := range r.Form["places"] {
				hv, err := metadata.ParseHierValue(p)
				if err != nil {
					errs.Errors = append(errs.Errors, err.Error())
				} else {
					hvs = append(hvs, hv)
				}
			}
		}
		if err := provider.SetPlaces(hvs); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	if r.Form["topics"] != nil {
		save = true
		var hvs []metadata.HierValue
		if len(r.Form["topics"]) > 1 || r.Form["topics"][0] != "" {
			for _, p := range r.Form["topics"] {
				hv, err := metadata.ParseHierValue(p)
				if err != nil {
					errs.Errors = append(errs.Errors, err.Error())
				} else {
					hvs = append(hvs, hv)
				}
			}
		}
		if err := provider.SetTopics(hvs); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if len(errs.Errors) == 0 && save {
		if err := filefmts.Save(handler, files[index]); err != nil {
			errs.Errors = append(errs.Errors, err.Error())
		}
	}
	if len(errs.Errors) != 0 {
		json.NewEncoder(w).Encode(&errs)
		return
	}
	json.NewEncoder(w).Encode(metadataForImage(filepath.Base(files[index]), provider))
}

func metadataForImage(filename string, provider metadata.Provider) (md *imgmd) {
	md = new(imgmd)
	md.Filename = filename
	md.Artist = provider.Creator()
	md.Caption = provider.Caption()
	md.DateTime = provider.DateTime().String()
	if len(md.DateTime) > 10 {
		md.DateTime = md.DateTime[:10] + " " + md.DateTime[11:]
		dot := strings.IndexByte(md.DateTime, '.')
		tz := strings.IndexAny(md.DateTime[10:], "-+Z")
		if dot >= 0 && tz < 0 {
			md.DateTime = md.DateTime[:dot]
		} else if dot >= 0 && tz >= 0 {
			md.DateTime = md.DateTime[:dot] + " " + md.DateTime[tz+10:]
		} else if tz >= 0 {
			md.DateTime = md.DateTime[:tz+10] + " " + md.DateTime[tz+10:]
		}
	}
	if len(md.DateTime) >= 10 {
		if t, err := time.Parse("2006-01-02", md.DateTime[:10]); err == nil {
			md.DateTime = t.Format("Mon ") + md.DateTime
		}
	}
	md.GPS = provider.GPS().String()
	for _, g := range provider.Groups() {
		md.Groups = append(md.Groups, g.String())
	}
	if md.Groups == nil {
		md.Groups = make([]string, 0)
	}
	for _, k := range provider.Keywords() {
		md.Keywords = append(md.Keywords, k.String())
	}
	if md.Keywords == nil {
		md.Keywords = make([]string, 0)
	}
	md.Location = provider.Location().String()
	for _, p := range provider.People() {
		md.People = append(md.People, p)
	}
	if md.People == nil {
		md.People = make([]string, 0)
	}
	for _, p := range provider.Places() {
		md.Places = append(md.Places, p.String())
	}
	if md.Places == nil {
		md.Places = make([]string, 0)
	}
	md.Title = provider.Title()
	for _, t := range provider.Topics() {
		md.Topics = append(md.Topics, t.String())
	}
	if md.Topics == nil {
		md.Topics = make([]string, 0)
	}
	return md
}

type rdsk struct {
	fs.File
}

func (f rdsk) Seek(a int64, b int) (c int64, d error) {
	return f.File.(io.Seeker).Seek(a, b)
}
