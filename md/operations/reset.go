package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

func newResetOp() Operation { return new(resetOp) }

// resetOp resets the values of one or more fields to their primary value, thus
// clearing up any inconsistencies or tagging errors.
type resetOp struct {
	fieldListOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *resetOp) parseArgs(args []string) (remainingArgs []string, err error) {
	remainingArgs, _ = op.fieldListOp.parseArgs(args)
	if len(op.fields) == 0 {
		op.fields = []fields.Field{
			fields.ArtistField,
			fields.CaptionField,
			fields.DateTimeField,
			fields.FacesField,
			fields.GPSField,
			fields.GroupsField,
			fields.KeywordsField,
			fields.LocationField,
			fields.PeopleField,
			fields.PlacesField,
			fields.TitleField,
			fields.TopicsField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *resetOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *resetOp) Run(files []MediaFile) error {
	for i, file := range files {
		for _, field := range op.fields {
			values := field.GetValues(file.Handler)
			if err := field.SetValues(file.Handler, values); err != nil {
				return fmt.Errorf("%s: reset %s: %s", file.Path, field.PluralName(), err)
			}
			files[i].Changed = true
		}
	}
	return nil
	// NOTE: this doesn't actually correct all tagging errors.  If the wrong
	// tags are set, or tags are set to the wrong value, that will be fixed.
	// But if the right tags are set, but their encoding wasn't correct, the
	// encoding won't be fixed because the underlying file handler won't see
	// a change to the field value and won't rewrite the tag.  I can live
	// with that.
}
