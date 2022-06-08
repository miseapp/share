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
func New(source *Source) (*Share, error) {
	// init files service
	files, err := files.New()
	if err != nil {
		return nil, err
	}

	// init command
	share := Init(files, source)

	return share, nil
}

func Init(files *files.Files, source *Source) *Share {
	return &Share{
		files:  files,
		source: source,
	}
}

// -- i/command

// invokes the share command
func (s *Share) Call() (string, error) {
	shared := NewSharedFile(s.source)
	res, err := s.files.Create(shared)
	return res, err
}
