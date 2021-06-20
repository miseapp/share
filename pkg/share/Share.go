package share

import "mise-share/pkg/share/files"

// -- types --

// Share is a command that shares a file
type Share struct {
	// deps
	files *files.Files
	// props
	source *Source
}

// -- impls --

// New inits a new share command
func New(source *Source) *Share {
	return &Share{
		files:  files.New(),
		source: source,
	}
}

// -- i/command

// Call invokes the share command
func (s *Share) Call() (string, error) {
	shared := NewSharedFile(s.source)
	s.files.Create(shared.AsHtml())
	return "ok", nil
}
