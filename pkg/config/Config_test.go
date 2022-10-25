package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestNew_U(t *testing.T) {
	cfg := New()

	assert.Equal(t, "http://localhost:4566", cfg.LocalUrl)
	assert.Equal(t, "us-east-1", cfg.Region)
	assert.Equal(t, "share-count", cfg.CountName)
	assert.Equal(t, "share-files", cfg.FilesName)
	assert.Equal(t, "http://share-files.s3.localhost.localstack.cloud:4566", cfg.FilesHost)
}

func TestResolveLiveEndpoint_U(t *testing.T) {
	cfg := Config{LocalUrl: ""}

	_, err := cfg.ResolveEndpoint("DynamoDB", "us-test-1")
	assert.NotNil(t, err)
}

func TestResolveLocalEndpoint_U(t *testing.T) {
	cfg := Config{LocalUrl: "http://127.0.0.1:4566"}

	endpoint, err := cfg.ResolveEndpoint("DynamoDB", "us-test-1")
	assert.Nil(t, err)
	assert.Equal(t, "http://127.0.0.1:4566", endpoint.URL)
}
