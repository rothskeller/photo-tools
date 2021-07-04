package operations

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Write sets the caption on a media file, reading it from standard input.
func Write(args []string, files []MediaFile) (err error) {
	var value string

	switch len(args) {
	case 0:
		return errors.New("write: missing field name")
	case 1:
		break
	default:
		return errors.New("write: excess arguments")
	}
	if field := fields.ParseField(args[0]); field == nil {
		return errors.New("write: missing field name")
	} else if field != fields.CaptionField {
		return fmt.Errorf("write: not supported for %q", field.Name())
	}
	if by, err := io.ReadAll(os.Stdin); err == nil {
		value = string(by)
	} else {
		return fmt.Errorf("write: standard input: %s", err)
	}
	for i, file := range files {
		if err := fields.CaptionField.SetValues(file.Handler, []interface{}{value}); err != nil {
			return fmt.Errorf("%s: write caption: %s", file.Path, err)
		}
		files[i].Changed = true
	}
	return nil
}
