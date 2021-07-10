package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Set sets the values of a field.
func Set(args []string, files []MediaFile) (err error) {
	var (
		field fields.Field
		toset []interface{}
	)
	if field, toset, err = parseFieldValues("set", args); err != nil {
		return err
	}
	for i, file := range files {
		if err := field.SetValues(file.Provider, toset); err != nil {
			return fmt.Errorf("%s: set %s: %s", file.Path, field.PluralName(), err)
		}
		files[i].Changed = true
	}
	return nil
}
