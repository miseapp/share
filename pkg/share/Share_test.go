package share

import (
	"context"
	"mise-share/pkg/share/files"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestMain(m *testing.M) {
	os.Setenv("AWS_ENDPOINT", "http://localhost:4566")
	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	os.Exit(m.Run())
}

func TestShare_I(t *testing.T) {
	files, err := files.New()
	assert.Equal(t, nil, err)

	// have to insert the initial item before updating
	// TODO: add this as a build step?
	files.Db.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String("share.count"),
			Item: map[string]types.AttributeValue{
				"Id":    &types.AttributeValueMemberS{Value: "share-files"},
				"Count": &types.AttributeValueMemberN{Value: "0"},
			},
		},
	)

	share := Init(
		files,
		&Source{
			Url: strp("https://httpbin.org/get"),
		},
	)

	res, err := share.Call()
	assert.Equal(t, nil, err)
	assert.Equal(t, "1.html", res)
}
