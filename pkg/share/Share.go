package share

import (
	"mise-share/pkg/config"
	"mise-share/pkg/share/files"
)

// -- types --

// Share is a command that shares a file
type Share struct {
	// the config
	cfg *config.Config

	// the files repo
	files *files.Files

	// props
	source SharedSource
}

// -- lifetime --

// creates a new share command
func New(source SharedSource) (*Share, error) {
	// create a config
	cfg := config.New()

	// init files service
	files, err := files.New(cfg)
	if err != nil {
		return nil, err
	}

	// init command
	share := Init(
		cfg,
		files,
		source,
	)

	return share, nil
}

// creates a new share command w/ a files service
func Init(
	cfg *config.Config,
	files *files.Files,
	source SharedSource,
) *Share {
	return &Share{
		cfg:    cfg,
		files:  files,
		source: source,
	}
}

// -- command --

// invokes the share command
func (s *Share) Call() (string, error) {
	// build the shared file
	shared := NewSharedFile(s.source, &s.cfg.FilesHost)

	// and create it
	url, err := s.files.Create(shared)
	if err != nil {
		return "", err
	}

	return url, nil
}
