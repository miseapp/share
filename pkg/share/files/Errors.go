package files

// -- types --

// an error when the file count could not be determined
type MissingCountError struct {
}

// -- impls --
func (*MissingCountError) Error() string {
	return "files.count is missing"
}
