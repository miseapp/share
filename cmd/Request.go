package main

// -- types --

// the api event structure
type Request struct {
	Source *RequestSource `json:"source"`
}

// the api event's source for a shared file
type RequestSource struct {
	Url *string `json:"url"`
}
