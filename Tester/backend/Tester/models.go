package tester

import (
	"context"
	"net/http"
)



type SequnceState struct {
	State      map[int64][]ActionResponse
	Order      []int64
	CustomData map[string]any
}

type ActionResponse struct {
    Success    bool
    Message    string
    Body       string
    StatusCode int
    Status     string
    Headers    http.Header
    Cookies    []*http.Cookie
    Data       map[string]any
}





type ActionContext struct{
	ReqContextId int
	ActionId int
	ActionType string
} 

type RequestExecutionContext struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}