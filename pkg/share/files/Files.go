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
	// the s3 client
	S3 *s3.Client

	// the dynamodb client
	Db *dynamodb.Client

	// the files config
	cfg *config.Config
}

// -- impls --
func New() (*Files, error) {
	// init our config
	cfg := config.New()

	// init aws config
	aws, err := cfg.InitAws()
	if err != nil {
		return nil, err
	}

	// init repo
	files := &Files{
		S3:  s3.NewFromConfig(aws),
		Db:  dynamodb.NewFromConfig(aws),
		cfg: cfg,
	}

	return files, nil
}

// -- i/commands

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
		log.Println("[Files.Create] update failed", err)
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

	// build the redirect file
	file, err := NewFile(count, content)
	if err != nil {
		return "", err
	}

	// insert the redirect file
	_, err = f.S3.PutObject(
		context.TODO(),
		&s3.PutObjectInput{
			Key:             aws.String(file.Key),
			Body:            file.Body,
			ContentType:     aws.String("text/html"),
			ContentLength:   int64(file.Length),
			ContentLanguage: aws.String("en-US"),
			ContentMD5:      aws.String(base64.StdEncoding.EncodeToString(file.Hash[:])),
			Bucket:          aws.String(f.cfg.FilesName),
		},
	)

	if err != nil {
		return "", err
	}

	// format the url
	url := fmt.Sprintf(
		"%s://%s.%s/%s",
		f.cfg.Scheme,
		f.cfg.FilesName,
		f.cfg.HostForService(s3.ServiceID),
		file.Key,
	)

	return url, nil
}
