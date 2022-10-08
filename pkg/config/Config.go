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
	// the url scheme
	Scheme string

	// the aws host
	Host string

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
	host := ""
	if os.Getenv("LOCAL") != "" {
		host = "localhost:4566"
	}

	// if localstack host exists, use that
	name := os.Getenv("LOCALSTACK_HOSTNAME")
	port := os.Getenv("EDGE_PORT")

	if name != "" && port != "" {
		host = fmt.Sprintf("%s:%s", name, port)
	}

	// use https, unless we have a local host
	scheme := "https"
	if host != "" {
		scheme = "http"
	}

	// build config
	// TODO: may neeed a separate url for files https://localhost.localstack.cloud:4566
	return &Config{
		Scheme:    scheme,
		Host:      host,
		Region:    os.Getenv("AWS_REGION"),
		CountName: os.Getenv("SHARE_COUNT_NAME"),
		FilesName: os.Getenv("SHARE_FILES_NAME"),
	}
}

// -- queries --

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

// the endpoint for the aws service
func (c *Config) HostForService(service string) string {
	// get the config host, if any
	host := c.Host
	if host == "" {
		return ""
	}

	// if local, s3 needs a url that can add subdomains
	if service == s3.ServiceID {
		return strings.Replace(host, "localhost", "s3.localhost.localstack.cloud", 1)
	}

	return host
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
	// get the config host, if any
	host := c.HostForService(service)
	if host == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL:           fmt.Sprintf("%s://%s", c.Scheme, host),
		SigningRegion: region,
	}

	return endpoint, nil
}
