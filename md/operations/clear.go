package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Clear removes all values of the specified field.
func Clear(args []string, files []MediaFile) (err error) {
	var field fields.Field

	switch len(args) {
	case 0:
		return errors.New("clear: missing field name")
	case 1:
		break
	default:
		return errors.New("clear: excess arguments")
	}
	if field = fields.ParseField(args[0]); field == nil {
		return fmt.Errorf("clear: %q is not a recognized field name", args[0])
	}
	for i, file := range files {
		if err := field.SetValues(file.Handler, nil); err != nil {
			return fmt.Errorf("%s: clear %s: %s", file.Path, field.PluralName(), err)
		}
		files[i].Changed = true
	}
	return nil
}
