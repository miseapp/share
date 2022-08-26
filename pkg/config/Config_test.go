package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestNew_U(t *testing.T) {
	cfg := New()
	assert.Equal(t, cfg.Endpoint, "http://localhost:4566")
	assert.Equal(t, cfg.CountName, "share-count")
	assert.Equal(t, cfg.FilesName, "share-files")
}
