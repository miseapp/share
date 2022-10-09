package request

import "encoding/json"

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

// deocde the request body from the string
func DecodeRequestBody(encoded string) (*RequestBody, error) {
	var body RequestBody

	// decode the body
	err := json.Unmarshal([]byte(encoded), &body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}
