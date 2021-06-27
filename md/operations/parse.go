package operations

// ParseOperation parses the operation, if any, specified at the beginning of
// the provided argument list.  If it finds one, it returns an operation
// instance for the parsed operation, the arguments remaining on the list after
// the operation's keyword and parameters have been removed, and a nil error.
// If it doesn't find one, it returns a nil operation instance, the unchanged
// argument list, and a nil error.  If it finds an operation but its parameters
// are invalid, it returns nil, nil, and an error.
func ParseOperation(args []string) (op Operation, remainingArgs []string, err error) {
	switch args[0] {
	case "add":
		op = newAddOp()
	case "check":
		op = newCheckOp()
	case "choose":
		op = newChooseOp()
	case "clear":
		op = newClearOp()
	case "copy":
		op = newCopyOp()
	case "read":
		op = newReadOp()
	case "remove":
		op = newRemoveOp()
	case "reset":
		op = newResetOp()
	case "set":
		op = newSetOp()
	case "show":
		op = newShowOp()
	case "tags":
		op = newTagsOp()
	case "write":
		op = newWriteOp()
	default:
		return nil, args, nil
	}
	if remainingArgs, err = op.parseArgs(args[1:]); err != nil {
		return nil, nil, err
	}
	return op, remainingArgs, nil
}
