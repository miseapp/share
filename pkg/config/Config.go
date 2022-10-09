package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// -- types --

// a repo for a collection of remote files
type Config struct {
	// the aws url, if any
	Url string

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
	// host is for local redirection. an empty host will use the sdk default,
	// which is a live aws url
	url := ""
	if os.Getenv("LOCAL") != "" {
		url = os.Getenv("LOCAL_URL")
	}

	// if localstack host exists, use that
	name := os.Getenv("LOCALSTACK_HOSTNAME")
	port := os.Getenv("EDGE_PORT")

	if name != "" && port != "" {
		url = fmt.Sprintf("%s:%s", name, port)
	}

	// build config
	// TODO: may neeed a separate url for files https://localhost.localstack.cloud:4566
	return &Config{
		Url:       url,
		Region:    os.Getenv("AWS_REGION"),
		CountName: os.Getenv("SHARE_COUNT_NAME"),
		FilesName: os.Getenv("SHARE_FILES_NAME"),
	}
}

// -- queries --

// if this is prod
func (c *Config) IsProd() bool {
	return c.Url == ""
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
	// get the config url, if any
	url := c.Url
	if url == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// if local, s3 needs a url that can add subdomains
	if service == s3.ServiceID {
		url = strings.Replace(url, "localhost", "s3.localhost.localstack.cloud", 1)
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL:           url,
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
