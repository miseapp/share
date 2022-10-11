package request

import (
	"encoding/base64"
	"encoding/json"
)

// -- types --

// the api event structure
type RequestBody struct {
	Source *RequestSource `json:"source"`
}

// the api event's source for a shared file
type RequestSource struct {
	Url  *string `json:"url,omitempty"`
	Json *string `json:"json,omitempty"`
}

// -- impls --
// deocde a request's body string
func DecodeRequestBody(str string, isBase64 bool) (*RequestBody, error) {
	// convert to data
	var data []byte
	if !isBase64 {
		data = []byte(str)
	} else {
		d, err := base64.StdEncoding.DecodeString(str)
		if err != nil {
			return nil, err
		}

		data = d
	}

	// decode the data
	var body RequestBody

	err := json.Unmarshal(data, &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
