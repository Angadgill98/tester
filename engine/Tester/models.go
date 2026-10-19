package tester

import (
	"net/http"

)


type SequnceState struct {
	State      map[string][]ActionResponse
	Order      []string
	CustomData map[string]any
}

type ActionResponse struct {
	Body       string
	StatusCode int
	Status     string
	Headers    http.Header
	Cookies    []*http.Cookie
	Data       map[string]any
}
