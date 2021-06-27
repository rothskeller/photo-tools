package operations

import (
	"fmt"

	"github.com/rothskeller/photo-tools/md/fields"
)

// fieldValueOp is a base class for operations that take a single field name and
// a value for that field.  It provides a parseArgs method for them.
type fieldValueOp struct {
	name  string
	field fields.Field
	value interface{}
}

func (op *fieldValueOp) parseArgs(args []string) (remainingArgs []string, err error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("%s: missing field name", op.name)
	}
	if op.field = fields.ParseField(args[0]); op.field == nil {
		return nil, fmt.Errorf("%s: missing field name", op.name)
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("%s: missing value", op.name)
	}
	if op.value, err = op.field.ParseValue(args[1]); err != nil {
		return nil, fmt.Errorf("%s %s: %s", op.name, args[0], err)
	}
	return args[2:], nil
}
