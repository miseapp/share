package share

import (
	"context"
	"mise-share/pkg/share/files"
	"mise-share/pkg/share/test"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestShare_I(t *testing.T) {
	files, err := files.New()
	assert.Equal(t, nil, err)

	// have to insert the initial item before updating
	// TODO: add this as a build step?
	files.Db.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: test.Str("share.count"),
			Item: map[string]types.AttributeValue{
				"Id":    &types.AttributeValueMemberS{Value: "share-files"},
				"Count": &types.AttributeValueMemberN{Value: "0"},
			},
		},
	)

	share := Init(
		files,
		&Source{
			Url: test.Str("https://httpbin.org/get"),
		},
	)

	res, err := share.Call()
	assert.Equal(t, nil, err)
	assert.Equal(t, "1.html", res)
}
