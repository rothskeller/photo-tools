package operations

import "github.com/rothskeller/photo-tools/md/fields"

// fieldListOp is a base class for operations that take a list of field names
// as arguments.  It provides a parseArgs method for them.
type fieldListOp struct {
	fields []fields.Field
}

func (op *fieldListOp) parseArgs(args []string) (remainingArgs []string, err error) {
	for len(args) != 0 {
		field := fields.ParseField(args[0])
		if field == nil {
			return args, nil
		}
		found := false
		for _, f := range op.fields {
			if f.Name() == field.Name() {
				found = true
				break
			}
		}
		if !found {
			op.fields = append(op.fields, field)
		}
	}
	return nil, nil
}
