package tester

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	// actions "tester/Actions"
	"tester/Actions"
	"tester/data"
)

type Tester struct {
	Client         http.Client
	WebSocket      *websocket.Conn
	Prev_res       data.ActionResponse
	Sequence_state data.SequnceState
}

func CreateTesterObj() Tester {
	return Tester{
		Client:    http.Client{},
		WebSocket: nil,
		Prev_res:  data.ActionResponse{},
		Sequence_state: data.SequnceState{
			State:      make(map[string]data.ActionResponse),
			CustomData: make(map[string]any),
		},
	}
}


func (tester *Tester) SetUpTester() {
	fmt.Println("Setting up tester");


	
}



// func (tester *Tester) ExecuteAction(action actions.Actions_obj) data.ActionResponse  {
	// var res = *action.Executor(&tester.Client, &tester.Prev_res, &tester.Sequence_state)
	// return  res
// }

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


func (tester *Tester)ExecuteHTTPAction(action_obj Actions.Actions_obj,action Actions.Http_actions)*data.ActionResponse  {
	switch action.Executor_type {
		case "req":
			executor := action_obj.Http.Config["http"][action.Name]

			httpExecutor, ok := executor.(func(*http.Client, *data.ActionResponse, *data.SequnceState) *data.ActionResponse)
			if !ok {
				fmt.Println("Invalid HTTP executor:", action.Name)
				return nil
			}

			res := httpExecutor(&tester.Client, &tester.Prev_res, &tester.Sequence_state)

			return res

		default:
			fmt.Println("Unknown HTTP executor type:", action.Executor_type)
			return nil
		}

	
}
func (tester *Tester) ExecuteWSAction(action_obj Actions.Actions_obj, action Actions.WsAction) *data.ActionResponse {
	executor := action_obj.WS.Config["ws"][action.Name]

	switch action.Exector_type {
	case "send":
		if tester.WebSocket == nil {
			fmt.Println("WebSocket connection does not exist. Connect first.")
			return nil
		}

		wsExecutor, ok := executor.(func(*websocket.Conn, *data.ActionResponse, *data.SequnceState) *data.ActionResponse)
		if !ok {
			fmt.Println("Invalid WebSocket Send executor:", action.Name)
			return nil
		}

		return wsExecutor(tester.WebSocket, &tester.Prev_res, &tester.Sequence_state)

	case "connect":
		wsExecutor, ok := executor.(func() (*websocket.Conn, error))
		if !ok {
			fmt.Println("Invalid WebSocket Connect executor:", action.Name)
			return nil
		}

		conn, err := wsExecutor()
		if err != nil {
			fmt.Println("WebSocket connection failed:", err)
			return nil
		}

		tester.WebSocket = conn
		return nil

	case "recieve":
		if tester.WebSocket == nil {
			fmt.Println("WebSocket connection does not exist. Connect first.")
			return nil
		}

		wsExecutor, ok := executor.(func(*websocket.Conn, *data.ActionResponse, *data.SequnceState) *data.ActionResponse)
		if !ok {
			fmt.Println("Invalid WebSocket Recieve executor:", action.Name)
			return nil
		}

		return wsExecutor(tester.WebSocket, &tester.Prev_res, &tester.Sequence_state)

	default:
		fmt.Println("Unknown WebSocket executor type:", action.Exector_type)
		return nil
	}
}
