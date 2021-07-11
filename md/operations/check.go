package operations

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/rothskeller/photo-tools/md/fields"
	"github.com/rothskeller/photo-tools/metadata"
)

var checkFields = []fields.Field{
	fields.ArtistField,
	fields.DateTimeField,
	fields.GPSField,
	fields.PlacesField,
	fields.PeopleField,
	fields.FacesField,
	fields.GroupsField,
	fields.TopicsField,
	fields.TitleField,
	fields.CaptionField,
	fields.KeywordsField,
	fields.LocationField,
}

// Check displays a table giving the tagging correctness of each field.
func Check(args []string, files []MediaFile) (err error) {
	var out *tabwriter.Writer

	if len(args) != 0 {
		return errors.New("check: excess arguments")
	}
	out = tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', 0)
	fmt.Fprint(out, "FILE")
	for _, field := range checkFields {
		fmt.Fprintf(out, "\t%s", field.ShortLabel())
	}
	fmt.Fprintln(out)
	for _, file := range files {
		fmt.Fprint(out, file.Path)
		for _, field := range checkFields {
			fmt.Fprintf(out, "\t%s", checkField(file.Provider, field, true))
		}
		fmt.Fprintln(out)
	}
	out.Flush()
	return nil
}

func checkField(p metadata.Provider, field fields.Field, info bool) string {
	var (
		incorrect bool
		canon     = field.GetValues(p)
	)
	var _, tagValues = field.GetTags(p)
	for _, tvs := range tagValues {
		if equalValues(field, canon, tvs) {
			continue
		}
		if emptyValues(field, tvs) {
			incorrect = true
		} else {
			return "!=" // conflicting values
		}
	}
	if incorrect {
		return "[]"
	}
	if emptyValues(field, canon) {
		if field.Expected() {
			return "--"
		}
		return "  "
	}
	if !info {
		return "  "
	}
	if field.Multivalued() {
		return fmt.Sprintf("%2d", len(canon))
	}
	return " ✓"
}

func equalValues(field fields.Field, as, bs []interface{}) bool {
	// First, make sure every non-empty element of as is present in bs.
	for _, a := range as {
		if field.EmptyValue(a) {
			continue
		}
		var found = false
		for _, b := range bs {
			if field.EqualValue(a, b) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	// Then, make sure every non-empty element of bs is present in as.
	for _, b := range bs {
		if field.EmptyValue(b) {
			continue
		}
		var found = false
		for _, a := range as {
			if field.EqualValue(a, b) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func emptyValues(field fields.Field, vs []interface{}) bool {
	switch len(vs) {
	case 0:
		return true
	case 1:
		return field.EmptyValue(vs[0])
	default:
		return false
	}
}
