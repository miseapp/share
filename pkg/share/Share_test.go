package share

import (
	"mise-share/pkg/share/files"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestMain(m *testing.M) {
	os.Setenv("AWS_ENDPOINT", "http://localhost.localstack.cloud:4566")
	os.Setenv("AWS_ACCESS_KEY", "test")
	os.Setenv("AWS_SECRET_KEY", "test")
	os.Exit(m.Run())
}

func TestShare_I(t *testing.T) {
	files := files.New()
	share := Init(
		files,
		&Source{
			Url: strp("https://httpbin.org/get"),
		},
	)

	// have to insert the initial item before updating
	// TODO: add this as a build step?
	files.Db.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("share.count"),
		Item: map[string]*dynamodb.AttributeValue{
			"Id":    {S: aws.String("share-files")},
			"Count": {N: aws.String("0")},
		},
	})

	res, err := share.Call()
	assert.Equal(t, nil, err)
	assert.Equal(t, "1.html", res)
}
