package tester

import (
	"bytes"
	"engine/Actions"
	"engine/Sequences"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Tester struct {
	Client         http.Client
	WebSocket      map[uuid.UUID]*websocket.Conn
	Prev_res       ActionResponse
	Sequence_state SequnceState
}

func CreateTesterObj() Tester {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	return Tester{
		Client: http.Client{Jar: jar},
		WebSocket: make(map[uuid.UUID]*websocket.Conn),
		Prev_res:  ActionResponse{},
		Sequence_state: SequnceState{
			State:      make(map[Action_name][]ActionResponse),
			CustomData: make(map[string]any),
		},
	}
}


func (tester *Tester) SetUpTester() {
	*tester = CreateTesterObj()
}


// func (tester *Tester) ExecuteAction(action actions.Actions_obj) data.ActionResponse  {
	// var res = *action.Executor(&tester.Client, &tester.Prev_res, &tester.Sequence_state)
	// return  res
// }


func (state SequnceState) ToJSON() map[string]any {
	return map[string]any{
		"state":       state.State,
		"custom_data": state.CustomData,
	}
}

func (tester *Tester) ShowState() {
	fmt.Println()
	fmt.Println("========== TESTER STATE ==========")

	fmt.Println("Previous Response:")
	fmt.Printf("%+v\n", tester.Prev_res)

	fmt.Println()
	fmt.Println("Sequence State:")

	for key, value := range tester.Sequence_state.State {
		fmt.Printf("  %s: %+v\n", key, value)
	}

	fmt.Println()
	fmt.Println("==================================")
	fmt.Println()
}


func (tester *Tester) ExecuteHTTPAction(action *Actions.Http_action,action_id uuid.UUID) *ActionResponse {
	opts := action.Opts

	reqURL, err := url.Parse(opts.Url)
	if err != nil {
		return nil
	}

	query := reqURL.Query()
	for _, item := range opts.Query {
		parts := strings.SplitN(item, "=", 2)
		if len(parts) == 2 {
			query.Set(parts[0], parts[1])
		}
	}
	reqURL.RawQuery = query.Encode()

	req, err := http.NewRequest(opts.Method, reqURL.String(), bytes.NewReader(opts.Body))
	if err != nil {
		return nil
	}

	for key, value := range opts.Headers {
		req.Header.Set(key, value)
	}

	resp, err := tester.Client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	response := &ActionResponse{
		Body:       string(body),
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		Data:       make(map[string]any),
	}

	tester.Prev_res = *response

	

	return response
}

func (tester *Tester) ExecuteSequence(sequence Sequences.Sequence, http_action_config *Actions.Http_Actions_config, ws_action_config *Actions.WS_Actions_config) map[string]any {
	tester.SetUpTester()

	for _, action := range sequence.Action {
		switch action.Type {
		case "http":
			httpAction, exists := (*http_action_config)[action.ID]
			if !exists {
				continue
			}

			actionResponse := tester.ExecuteHTTPAction(httpAction, action.ID)
			if actionResponse == nil {
				continue
			}

			actionName := Action_name(httpAction.Name)

			if tester.Sequence_state.State == nil {
				tester.Sequence_state.State = make(map[Action_name][]ActionResponse)
			}

			tester.Sequence_state.State[actionName] = append(
				tester.Sequence_state.State[actionName],
				*actionResponse,
			)

		// case "ws":
		// 	wsAction, exists := (*ws_action_config)[action.ID]
		// 	if !exists {
		// 		continue
		// 	}
		// 	tester.ExecuteWSAction(&wsAction, action.ID)
		}
	}

	return tester.Sequence_state.ToJSON()
}
// func (tester *Tester) ExecuteWSAction(action Actions.WsAction) *ActionResponse {
// 	executor := action_obj.WS.Config["ws"][action.Name]

// 	switch action.Exector_type {
// 	case "send":
// 		if tester.WebSocket == nil {
// 			fmt.Println("WebSocket connection does not exist. Connect first.")
// 			return nil
// 		}

// 		wsExecutor, ok := executor.(func(*websocket.Conn, *data.ActionResponse, *data.SequnceState) *data.ActionResponse)
// 		if !ok {
// 			fmt.Println("Invalid WebSocket Send executor:", action.Name)
// 			return nil
// 		}

// 		return wsExecutor(tester.WebSocket, &tester.Prev_res, &tester.Sequence_state)

// 	case "connect":
// 		wsExecutor, ok := executor.(func() (*websocket.Conn, error))
// 		if !ok {
// 			fmt.Println("Invalid WebSocket Connect executor:", action.Name)
// 			return nil
// 		}

// 		conn, err := wsExecutor()
// 		if err != nil {
// 			fmt.Println("WebSocket connection failed:", err)
// 			return nil
// 		}

// 		tester.WebSocket = conn
// 		return nil

// 	case "recieve":
// 		if tester.WebSocket == nil {
// 			fmt.Println("WebSocket connection does not exist. Connect first.")
// 			return nil
// 		}

// 		wsExecutor, ok := executor.(func(*websocket.Conn, *data.ActionResponse, *data.SequnceState) *data.ActionResponse)
// 		if !ok {
// 			fmt.Println("Invalid WebSocket Recieve executor:", action.Name)
// 			return nil
// 		}

// 		return wsExecutor(tester.WebSocket, &tester.Prev_res, &tester.Sequence_state)

// 	default:
// 		fmt.Println("Unknown WebSocket executor type:", action.Exector_type)
// 		return nil
// 	}
// }
