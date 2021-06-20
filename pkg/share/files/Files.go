package files

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

// -- types --

// a repo for a collection of remote files
type Files struct {
	s3 *s3.S3
	db *dynamodb.DynamoDB
}

// -- impls --
func New() *Files {
	// create aws session
	session := session.Must(session.NewSession())

	dynamodb.New(session)

	// create module
	return &Files{
		s3: s3.New(session),
		db: dynamodb.New(session),
	}
}

// -- i/commands

// creates a new file with the given body, returning the filename
func (f *Files) Create(body string) (string, error) {
	// atomically increment the counter
	// see: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/WorkingWithItems.html#WorkingWithItems.AtomicCounters
	res, err := f.db.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: aws.String("share.files")},
		},
		UpdateExpression: aws.String("SET Count = Count + 1"),
		ReturnValues:     aws.String(dynamodb.ReturnValueUpdatedNew),
		TableName:        aws.String("share.count"),
	})

	if err != nil {
		return "", err
	}

	// this shouldn't happen, UpdateItem inserts whent he key is missing
	if res == nil {
		return "", &MissingCountError{}
	}

	// grab the new count as an integer
	count, err := strconv.Atoi(*res.Attributes["Count"].N)
	if err != nil {
		return "", err
	}

	// encode the filename
	filename := Filename(count)

	s, err := filename.String()
	if err != nil {
		return "", err
	}

	// build a file
	file := NewFile(body)

	// create the file
	f.s3.PutObject(&s3.PutObjectInput{
		Key:             aws.String(s),
		Body:            file.Body,
		ContentType:     aws.String("text/html"),
		ContentLength:   aws.Int64(int64(file.Length)),
		ContentLanguage: aws.String("en-US"),
		ContentMD5:      aws.String(string(file.Hash[:])),
		Bucket:          aws.String("share.files"),
	})

	return s, nil
}
