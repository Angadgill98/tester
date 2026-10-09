package Global

import "net/http"

type Global struct {
	Http_ActionsConfigPath     string
	Ws_ActionsConfigPath     string

	SequenceConfigPath    string
	ApplicationConfigPath string
	ActionStoragePath     string

	HooksPath     string

}

var GlobalConfig = Global{
	Http_ActionsConfigPath:     "./Actions/actions_config.go",
	Ws_ActionsConfigPath:     "./Actions/ws_actions_config.go",

	SequenceConfigPath:    "./Sequence/sequence_config.json",
	ApplicationConfigPath: "./Application/application_config.json",
	ActionStoragePath:     "data/actions/",

	HooksPath:     "hooks/",

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
