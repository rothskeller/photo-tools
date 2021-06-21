package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	webdialogs "github.com/rothskeller/photo-tools/web-dialogs"
)

func receiveInfo(w http.ResponseWriter, r *http.Request) webdialogs.DialogHandler {
	if r.Method != http.MethodPost {
		return requestInfo
	}
	var variants = make([]string, len(groupSplit))
	var varreps = make([]string, len(groupSplit))
	var varmap = make(map[string]int)
	for idx := range groupSplit {
		variants[idx] = r.FormValue(fmt.Sprintf("v%d", idx))
		if _, ok := varmap[variants[idx]]; ok {
			return requestInfo
		}
		varmap[variants[idx]] = idx
		varreps[idx] = r.FormValue(fmt.Sprintf("g%d", idx))
	}
	if ovar, ok := varmap[""]; !ok || varreps[ovar] == "" {
		return requestInfo
	}

	renames = map[string]string{}
	for idx, group := range groupSplit {
		var base = groupNames[0]
		if variants[idx] != "" {
			base = base + "." + variants[idx]
		}
		var remaining = make(map[string]bool)
		for _, file := range group {
			remaining[file] = true
		}
		var rep = varreps[idx]
		if rep == "" {
			for _, file := range group {
				if !isXMP(file) {
					rep = file
					break
				}
			}
		}
		if rep != "" {
			addRename(base, rep, remaining)
		}
		var seq = 1
		for _, file := range group {
			if !isXMP(file) && remaining[file] {
				addRename(fmt.Sprintf("%s.md%d", base, seq), file, remaining)
				seq++
			}
		}
		for _, file := range group {
			if remaining[file] {
				addRename(fmt.Sprintf("%s.md%d", base, seq), file, remaining)
				seq++
			}
		}
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := confirmTemplate.Execute(w, struct {
		Paths   []string
		Renames map[string]string
	}{
		Paths:   groupPaths[0],
		Renames: renames,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	return doRenames
}

func addRename(base, file string, remaining map[string]bool) {
	renames[file] = base + canonExt(file)
	delete(remaining, file)
	if remaining[file+".xmp"] {
		renames[file+".xmp"] = base + canonExt(file) + ".xmp"
		delete(remaining, file+".xmp")
	}
	if remaining[file+".xmp_original"] {
		renames[file+".xmp"] = base + canonExt(file) + ".xmp_original"
		delete(remaining, file+".xmp_original")
	}
}

func canonExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if strings.HasSuffix(ext, "_original") {
		ext = ext[:len(ext)-9]
	}
	switch ext {
	case "jpeg":
		ext = "jpg"
	}
	return ext
}

var confirmTemplate = template.Must(template.New("confirm").Parse(`<!DOCTYPE html><html><body>
  <form method="POST">
    <div>Confirm renames:</div>
    <table>
      {{ range .Paths }}
        <tr>
	  <td>{{.}}&nbsp;</td>
	  <td>=&gt; {{ if eq . (index $.Renames .) }}(no change){{ else }}{{ index $.Renames . }}{{ end }}</td>
	</tr>
      {{ end }}
    </table>
    <div>
      <input type="submit" name="confirm" value="Yes">&nbsp;
      <input type="submit" name="confirm" value="No">
    </div>
  </form>
</body></html>`))
