package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Remove removes one or more values from a multi-valued field.
func Remove(args []string, files []MediaFile) (err error) {
	var field fields.Field
	var toremove []interface{}

	if field, toremove, err = parseFieldValues("remove", args); err != nil {
		return err
	}
	if !field.Multivalued() {
		return fmt.Errorf("remove: not supported for %q", field.Name())
	}
	for i, file := range files {
		// Get the current values.
		values := field.GetValues(file.Provider)
		// Remove the ones we were asked to remove.
		j := 0
		for _, v := range values {
			var found = false
			for _, rem := range toremove {
				if field.EqualValue(v, rem) {
					found = true
					break
				}
			}
			if !found {
				values[j] = v
				j++
			}
		}
		// If we found any, set the new value list that leaves it out.
		if j < len(values) {
			if err := field.SetValues(file.Provider, values[:j]); err != nil {
				return fmt.Errorf("%s: remove %s: %s", file.Path, field.Name(), err)
			}
			files[i].Changed = true
		}
	}
	return nil
}
