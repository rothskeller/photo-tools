package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// parseFieldList expects the argument list to be a list of field names.  It
// returns the (deduped) list of fields, or an error.
func parseFieldList(opname string, args []string) (flist []fields.Field, err error) {
	for len(args) != 0 {
		field := fields.ParseField(args[0])
		if field == nil {
			return nil, fmt.Errorf("%s: %q is not a recognized field name", opname, args[0])
		}
		found := false
		for _, f := range flist {
			if f.Name() == field.Name() {
				found = true
				break
			}
		}
		if !found {
			flist = append(flist, field)
		}
		args = args[1:]
	}
	return flist, nil
}
