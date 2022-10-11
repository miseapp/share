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
	// the aws url, if any
	LocalUrl string

	// the aws region
	Region string

	// the count table name
	CountName string

	// the files bucket name
	FilesName string
}

// -- lifetime --

// create a new config
func New() *Config {
	// get local url for redirection. an empty url uses the sdk default, which
	// is a live aws url
	localUrl := os.Getenv("LOCAL_URL")

	// if localstack host exists, use that instead
	name := os.Getenv("LOCALSTACK_HOSTNAME")
	port := os.Getenv("EDGE_PORT")
	if name != "" && port != "" {
		localUrl = fmt.Sprintf("http://%s:%s", name, port)
	}

	// build config
	// TODO: may neeed a separate url for files https://localhost.localstack.cloud:4566
	return &Config{
		LocalUrl:  localUrl,
		Region:    os.Getenv("AWS_REGION"),
		CountName: os.Getenv("SHARE_COUNT_NAME"),
		FilesName: os.Getenv("SHARE_FILES_NAME"),
	}
}

// -- queries --

// if this is prod
func (c *Config) IsProd() bool {
	return c.LocalUrl == ""
}

// find the aws endpoint given a service and region
func (c *Config) ResolveEndpoint(
	service string,
	region string,
	options ...interface{},
) (
	aws.Endpoint,
	error,
) {
	// if no local url, aws will use the default (live) endpoint
	if c.LocalUrl == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL:           c.LocalUrl,
		SigningRegion: region,
	}

	return endpoint, nil
}

// -- factories --

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
