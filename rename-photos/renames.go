package main

import (
	"fmt"
	"net/http"
	"os"

	webdialogs "github.com/rothskeller/photo-tools/web-dialogs"
)

func doRenames(w http.ResponseWriter, r *http.Request) webdialogs.DialogHandler {
	if r.Method != http.MethodPost {
		return requestInfo(w, r)
	}
	if r.FormValue("confirm") != "Yes" {
		return requestInfo(w, r)
	}
	for len(renames) != 0 {
		progress := false
		for op, np := range renames {
			if op == np {
				delete(renames, op)
				progress = true
				continue
			}
			if renames[np] != "" {
				continue
			}
			if err := os.Rename(op, np); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: rename %s -> %s: %s\n", op, np, err)
				os.Exit(1)
			}
			delete(renames, op)
			progress = true
		}
		if !progress {
			fmt.Fprintf(os.Stderr, "ERROR: unable to swap names; rename not completed\n")
			break
		}
	}
	groupNames = groupNames[1:]
	groupPaths = groupPaths[1:]
	return requestInfo(w, r)
}
