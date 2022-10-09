package share

import (
	"fmt"
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
	source *Source
}

// -- lifetime --

// creates a new share command
func New(source *Source) (*Share, error) {
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
	source *Source,
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
	// create the shared file
	shared := NewSharedFile(s.source)
	file, err := s.files.Create(shared)

	// if err, short circuit
	if err != nil {
		return "", err
	}

	// if not prod, return the raw url
	if !s.cfg.IsProd() {
		return file.Url, nil
	}

	// in prod, use our custom domain
	return fmt.Sprintf("https://share.miseapp.co/%s", file.Key), nil
}
