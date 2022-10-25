package files

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"mise-share/pkg/config"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// -- types --

// a repo for a collection of remote files
type Files struct {
	// the files config
	cfg *config.Config

	// the s3 client
	S3 *s3.Client

	// the s3 options
	s3Opts *s3.Options

	// the dynamodb client
	Db *dynamodb.Client
}

// -- lifetime --

/// create a files repo
func New(cfg *config.Config) (*Files, error) {
	// init aws config
	aws, err := cfg.InitAws()
	if err != nil {
		return nil, err
	}

	// sneakily capture the s3 options so we can get the url later
	var s3Opts *s3.Options = nil

	// init repo
	files := &Files{
		cfg: cfg,
		S3: s3.NewFromConfig(aws, func(opts *s3.Options) {
			// allow path style unless prod
			opts.UsePathStyle = !cfg.IsProd()
			// grab options ref
			s3Opts = opts
		}),
		Db: dynamodb.NewFromConfig(aws),
	}

	files.s3Opts = s3Opts

	return files, nil
}

// -- commands --

// creates a new file with the given content, returning the file key
func (f *Files) Create(content FileContent) (string, error) {
	// atomically increment the counter
	// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
	resCount, err := f.Db.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName: aws.String(f.cfg.CountName),
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberS{Value: f.cfg.FilesName},
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
		log.Println("[Files.Create] could not get next id", err)
		return "", err
	}

	// this shouldn't happen, UpdateItem inserts when the key is missing
	if resCount == nil {
		return "", &MissingCountError{}
	}

	// grab the new count as an integer
	var rec struct {
		Count string `json:"Count"`
	}

	err = attributevalue.UnmarshalMap(resCount.Attributes, &rec)
	if err != nil {
		log.Println("[Files.Create] could not unmarshal response", err)
		return "", err
	}

	count, err := strconv.Atoi(rec.Count)
	if err != nil {
		log.Println("[Files.Create] could not parse `Count` as integer", err)
		return "", err
	}

	// build the redirect input
	input, err := NewFileInput(count, content)
	if err != nil {
		log.Println("[Files.Create] could create input file", err)
		return "", err
	}

	// insert the redirect file
	_, err = f.S3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Key:             &input.Key,
			Body:            input.Body,
			ContentType:     aws.String("text/html"),
			ContentLength:   int64(input.Length),
			ContentLanguage: aws.String("en-US"),
			ContentMD5:      aws.String(base64.StdEncoding.EncodeToString(input.Hash[:])),
			Bucket:          &f.cfg.FilesName,
		},
	)

	if err != nil {
		log.Println("[Files.Create] could not create file", err)
		return "", err
	}

	// build the s3 file url
	url := fmt.Sprintf("%s/%s", f.cfg.FilesHost, input.Key)

	// and succeed w/ the url
	return url, nil
}
