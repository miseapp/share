package share

// -- types --

// [value] the source of a shared file's contents.
type SharedSource interface {
	// the mime type of the source
	Type() string

	// the value of the source
	Value() string
}

// [value] a source url for a file
type SourceUrl struct {
	// the source url
	Url *string
}

// [value] a source json payload for a file
type SourceJson struct {
	// the source url
	Json *string
}

// -- impls --

func (s *SourceUrl) Type() string {
	return "application/url" // making up a mime type
}

func (s *SourceUrl) Value() string {
	return *s.Url
}

func (s *SourceJson) Type() string {
	return "application/json"
}

func (s *SourceJson) Value() string {
	return *s.Json
}
