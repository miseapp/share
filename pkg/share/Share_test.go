package share

import (
	"context"
	"fmt"
	"mise-share/pkg/share/files"
	"mise-share/pkg/share/test"
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

// -- tests --
func TestShare_I(t *testing.T) {
	files, err := files.New()
	assert.Equal(t, nil, err)

	next := count(t, files) + 1

	share := Init(
		files,
		&Source{
			Url: test.Str("https://httpbin.org/get"),
		},
	)

	url, err := share.Call()
	assert.Equal(t, nil, err)
	assert.Equal(t, fmt.Sprintf("http://share-files.s3.localhost.localstack.cloud:4566/%d.html", next), url)
}

// -- helpers --
// count the current number of shares
func count(t *testing.T, files *files.Files) int {
	// grab the current count as an integer
	var rec struct {
		Count string `json:"Count"`
	}

	// get the current count
	res, err := files.Db.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			TableName: aws.String(
				os.Getenv("SHARE_COUNT_NAME"),
			),
			Key: map[string]types.AttributeValue{
				"Id": &types.AttributeValueMemberS{
					Value: os.Getenv("SHARE_FILES_NAME"),
				},
			},
		},
	)
	assert.Equal(t, nil, err)

	// decode it
	err = attributevalue.UnmarshalMap(res.Item, &rec)
	assert.Equal(t, nil, err)

	// as an integer
	count, err := strconv.Atoi(rec.Count)
	assert.Equal(t, nil, err)

	return count
}
