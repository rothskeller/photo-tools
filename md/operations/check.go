package operations

import "github.com/rothskeller/photo-tools/md/fields"

func newCheckOp() Operation { return new(checkOp) }

type checkOp struct {
	fieldListOp
}

// parseArgs parses the arguments for the operation, returning the remaining
// argument list or an error.
func (op *checkOp) parseArgs(args []string) (remainingArgs []string, err error) {
	remainingArgs, _ = op.fieldListOp.parseArgs(args)
	if len(op.fields) == 0 {
		op.fields = []fields.Field{
			fields.TitleField,
			fields.DateTimeField,
			fields.ArtistField,
			fields.GPSField,
			fields.LocationField,
			fields.PlacesField,
			fields.PeopleField,
			fields.GroupsField,
			fields.TopicsField,
			fields.OtherKeywordsField,
			fields.CaptionField,
		}
	}
	return remainingArgs, nil
}

// Check verifies that the operation is valid for the listed batches of media
// files.  (Some operations require certain numbers of batches, certain numbers
// of files per batch, certain media types, etc.).
func (op *checkOp) Check(batches [][]MediaFile) error { return nil }

// Run executes the operation against the listed media files (one batch).
func (op *checkOp) Run(files []MediaFile) error {
	panic("not implemented")
}
