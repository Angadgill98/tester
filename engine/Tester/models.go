package tester

import (
	"net/http"

	"github.com/google/uuid"
)

type Action_id uuid.UUID
type Action_name string

type SequnceState struct {
	State map[Action_id]map[Action_name]ActionResponse
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
