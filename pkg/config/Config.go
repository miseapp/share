package config

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
)

// -- types --

// a repo for a collection of remote files
type Config struct {
	// the aws endpoint
	Endpoint string

	// the count table name
	CountName string

	// the files bucket name
	FilesName string
}

// -- impls --

// create a new config
func New() *Config {
	// load dotenv if specified
	env := os.Getenv("SHARE_ENV")
	if env != "" {
		err := godotenv.Load("../../" + env)
		if err != nil {
			log.Fatal("[config] could not load .env @ ", env)
		}
	}

	// TODO: may neeed a separate url for files https://localhost.localstack.cloud:4566
	return &Config{
		Endpoint:  os.Getenv("AWS_ENDPOINT"),
		CountName: os.Getenv("SHARE_COUNT_NAME"),
		FilesName: os.Getenv("SHARE_FILES_NAME"),
	}
}

// init aws config
func (c *Config) InitAws() (aws.Config, error) {
	return config.LoadDefaultConfig(
		context.TODO(),
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
	// if there is an endpoint set in env
	url := c.Endpoint
	if url == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL: url,
	}

	log.Println("using endpoint", endpoint.URL)

	return endpoint, nil
}
