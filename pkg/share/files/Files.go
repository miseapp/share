package files

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// -- types --

// a repo for a collection of remote files
type Files struct {
	S3 *s3.Client
	Db *dynamodb.Client
}

// -- impls --
func New() (*Files, error) {
	// init aws config
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(ResolveEndpoint)),
		config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)

	if err != nil {
		return nil, err
	}

	// init repo
	files := &Files{
		S3: s3.NewFromConfig(cfg),
		Db: dynamodb.NewFromConfig(cfg),
	}

	return files, nil
}

// -- i/commands

// creates a new file with the given content, returning the file key
func (f *Files) Create(content FileContent) (string, error) {
	// atomically increment the counter
	// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
	res, err := f.Db.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName: aws.String("share.count"),
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberS{Value: "share-files"},
			},
			ExpressionAttributeNames: map[string]string{
				"#C": "Count",
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":incr": &types.AttributeValueMemberN{Value: "1"},
			},
			UpdateExpression: aws.String("SET #C = #C + :incr"),
			ReturnValues:     types.ReturnValueUpdatedNew,
		},
	)

	if err != nil {
		log.Println("[Files.Create] update failed", err)
		return "", err
	}

	// this shouldn't happen, UpdateItem inserts when the key is missing
	if res == nil {
		return "", &MissingCountError{}
	}

	// grab the new count as an integer
	var rec struct {
		Count string `json:"Count"`
	}

	err = attributevalue.UnmarshalMap(res.Attributes, &rec)
	if err != nil {
		log.Println("[Files.Create] could not unmarshal response", err)
		return "", err
	}

	count, err := strconv.Atoi(rec.Count)
	if err != nil {
		log.Println("[Files.Create] could not parse `Count` as integer", err)
		return "", err
	}

	// build the redirect file
	file, err := NewFile(count, content)
	if err != nil {
		return "", err
	}

	// insert the redirect file
	_, err = f.S3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Key:             aws.String(fmt.Sprintf("%s.html", file.Key)),
			Body:            file.Body,
			ContentType:     aws.String("text/html"),
			ContentLength:   int64(file.Length),
			ContentLanguage: aws.String("en-US"),
			ContentMD5:      aws.String(base64.StdEncoding.EncodeToString(file.Hash[:])),
			Bucket:          aws.String("share-files"),
		},
	)

	if err != nil {
		return "", err
	}

	return file.Key, nil
}

// -- i/Endpoint
func ResolveEndpoint(service, region string, options ...interface{}) (aws.Endpoint, error) {
	// if there is an endpoint set in env
	url := os.Getenv("AWS_ENDPOINT")
	if url == "" {
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}

	// use it instead of the default
	endpoint := aws.Endpoint{
		URL: url,
	}

	return endpoint, nil
}
