package Global

import "net/http"

type Global struct {

	HooksPath     string	


	SocketPath		string

}

var GlobalConfig = Global{

	HooksPath:     "hooks/",


	SocketPath:"/tmp/testflow.sock",

}

type SequnceState struct {
	State map[string]ActionResponse
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
