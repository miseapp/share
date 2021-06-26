package files

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// -- types --

// a repo for a collection of remote files
type Files struct {
	S3 *s3.S3
	Db *dynamodb.DynamoDB
}

// -- impls --
func New() *Files {
	// init aws session
	session := session.Must(session.NewSession(&aws.Config{
		Endpoint:    aws.String(os.Getenv("AWS_ENDPOINT")),
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewEnvCredentials(),
	}))

	// init repo
	return &Files{
		S3: s3.New(session),
		Db: dynamodb.New(session),
	}
}

// -- i/commands

// creates a new file with the given content, returning the file key
func (f *Files) Create(content FileContent) (string, error) {
	// atomically increment the counter
	// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
	res, err := f.Db.UpdateItem(&dynamodb.UpdateItemInput{
		TableName: aws.String("share.count"),
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: aws.String("share-files")},
		},
		ExpressionAttributeNames: map[string]*string{
			"#C": aws.String("Count"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":incr": {N: aws.String("1")},
		},
		UpdateExpression: aws.String("SET #C = #C + :incr"),
		ReturnValues:     aws.String(dynamodb.ReturnValueUpdatedNew),
	})

	if err != nil {
		return "", err
	}

	// this shouldn't happen, UpdateItem inserts when the key is missing
	if res == nil {
		return "", &MissingCountError{}
	}

	// grab the new count as an integer
	count, err := strconv.Atoi(*res.Attributes["Count"].N)
	if err != nil {
		return "", err
	}

	// build the redirect file
	file, err := NewFile(count, content)
	if err != nil {
		return "", err
	}

	// insert the redirect file
	_, err = f.S3.PutObject(&s3.PutObjectInput{
		Key:             aws.String(fmt.Sprintf("%s.html", file.Key)),
		Body:            file.Body,
		ContentType:     aws.String("text/html"),
		ContentLength:   aws.Int64(int64(file.Length)),
		ContentLanguage: aws.String("en-US"),
		ContentMD5:      aws.String(base64.StdEncoding.EncodeToString(file.Hash[:])),
		Bucket:          aws.String("share-files"),
	})

	if err != nil {
		return "", err
	}

	return file.Key, nil
}
