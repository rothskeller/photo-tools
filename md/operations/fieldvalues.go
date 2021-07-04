package operations

import (
	"fmt"
	"strings"

	"github.com/rothskeller/photo-tools/md/fields"
)

// parseFieldValues expects the argument list to be a field name followed by one
// or more values for the field.  It returns the field and values, or an error.
func parseFieldValues(opname string, args []string) (field fields.Field, values []interface{}, err error) {
	switch len(args) {
	case 0:
		return nil, nil, fmt.Errorf("%s: missing field name", opname)
	case 1:
		return nil, nil, fmt.Errorf("%s: missing value", opname)
	}
	if field = fields.ParseField(args[0]); field == nil {
		return nil, nil, fmt.Errorf("%s: %q is not a recognized field name", opname, args[0])
	}
	var argstrings = strings.Split(strings.Join(args[1:], " "), ";")
	for _, arg := range argstrings {
		if val, err := field.ParseValue(strings.TrimSpace(arg)); err == nil {
			values = append(values, val)
		} else {
			return nil, nil, fmt.Errorf("%s %s: %s", opname, args[0], err)
		}
	}
	return field, values, nil
}
