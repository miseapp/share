package config

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// -- types --

// a repo for a collection of remote files
type Config struct {
	// the aws endpoint
	Endpoint string

	// the aws region
	Region string

	// the count table name
	CountName string

	// the files bucket name
	FilesName string
}

// -- impls --

// create a new config
func New() *Config {
	endpoint := ""

	// if localstack endpoint exists, use that
	host := os.Getenv("LOCALSTACK_HOSTNAME")
	port := os.Getenv("EDGE_PORT")
	if host != "" && port != "" {
		endpoint = fmt.Sprintf("http://%s:%s", host, port)
	}

	// if not resolved and aws endpoint exists, use that
	if endpoint == "" {
		endpoint = os.Getenv("AWS_ENDPOINT")
	}

	// build config
	// TODO: may neeed a separate url for files https://localhost.localstack.cloud:4566
	return &Config{
		Endpoint:  endpoint,
		Region:    os.Getenv("AWS_REGION"),
		CountName: os.Getenv("SHARE_COUNT_NAME"),
		FilesName: os.Getenv("SHARE_FILES_NAME"),
	}
}

// init aws config
func (c *Config) InitAws() (aws.Config, error) {
	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(c.Region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(c.ResolveEndpoint),
		),
		config.WithClientLogMode(
			aws.LogRequestWithBody|aws.LogResponseWithBody,
		),
	)
}

// -- i/Endpoint

// resolve the config endpoint
func (c *Config) ResolveEndpoint(
	service string,
	region string,
	options ...interface{},
) (
	aws.Endpoint,
	error,
) {
	url := c.Endpoint

	// if there is no endpoint, error
	if url == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL:           url,
		SigningRegion: region,
	}

	return endpoint, nil
}
