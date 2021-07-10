package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Add adds one or more values to a multivalued field.
func Add(args []string, files []MediaFile) (err error) {
	var (
		field fields.Field
		toadd []interface{}
	)
	if field, toadd, err = parseFieldValues("add", args); err != nil {
		return err
	}
	if !field.Multivalued() {
		return fmt.Errorf("add: not supported for %q", field.Name())
	}
	for i, file := range files {
		// Get the current values.
		values := field.GetValues(file.Provider)
		// Add the desired values.
		for _, newv := range toadd {
			// Find out whether the value we're adding is already there.
			found := false
			for _, v := range values {
				if field.EqualValue(v, newv) {
					found = true
					break
				}
			}
			// If not, add it.
			if !found {
				values = append(values, newv)
				files[i].Changed = true
			}
		}
		if files[i].Changed {
			if err := field.SetValues(file.Provider, values); err != nil {
				return fmt.Errorf("%s: add %s: %s", file.Path, field.Name(), err)
			}
		}
	}
	return nil
}
