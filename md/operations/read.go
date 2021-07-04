package operations

import (
	"errors"
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// Read echoes the caption of a media file to standard out without any
// decoration.
func Read(args []string, files []MediaFile) (err error) {
	switch len(args) {
	case 0:
		return errors.New("read: missing field name")
	case 1:
		break
	default:
		return errors.New("read: excess arguments")
	}
	if field := fields.ParseField(args[0]); field == nil {
		return fmt.Errorf("read: %q is not a recognized field name", args[0])
	} else if field != fields.CaptionField {
		return fmt.Errorf("read: not supported for %q", field.Name())
	}
	if len(files) != 1 {
		return errors.New("read caption: only one file allowed")
	}
	caption := fields.CaptionField.GetValues(files[0].Handler)
	if len(caption) != 0 {
		str := fields.CaptionField.RenderValue(caption[0])
		fmt.Print(str)
		if len(str) != 0 && str[len(str)-1] != '\n' {
			fmt.Println()
		}
	}
	return nil
}
