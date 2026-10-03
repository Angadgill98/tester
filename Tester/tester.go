package tester

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	// actions "tester/Actions"
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
// 	var res = *action.Executor(&tester.Client, &tester.Prev_res, &tester.Sequence_state)
// 	return  res
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