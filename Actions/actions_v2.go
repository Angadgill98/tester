package Actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)



type Actions_obj struct {

	Http HTTP_Actions_obj
	WS WS_Actions_obj

}



type HTTP_Actions_obj struct {
	Http_actions map[string]Http_actions
	Config map[string]map[string]Executor_function_signature

}
type Http_actions struct{
	Name string
	Path string
	Executor_type string
	Executor_function_name string

}

var Http_executor_types = []string{
	"req",
}

type Executor_function_signature any

func CreateActionsObj() Actions_obj {
	return Actions_obj{
		Http: HTTP_Actions_obj{
			Http_actions: make(map[string]Http_actions),
			Config: map[string]map[string]Executor_function_signature{
				"http": make(map[string]Executor_function_signature),
			},
		},
		WS: WS_Actions_obj{
			WS_actions: make(map[string]WsAction),
			Config: map[string]map[string]Executor_function_signature{
				"ws": make(map[string]Executor_function_signature),
			},
		},
	}
}

func (this *Actions_obj) SetUpActionsObj() {
	isok, err := this.CheckHttpActionsFS()

	if err != nil {
		return
	}

	if !isok {
		fmt.Println("Http Actions FS not Setup")
		return
	}

	if isok {
		fmt.Println("Http Actions FS Setup OK")
	}

	fmt.Println("Getting http Actions files")

	actions, err := this.ReadHttpActionsFiles()
	if err != nil {
		fmt.Println("Failed to read Actions files:", err)
		return
	}

	for _, action := range actions {
		this.Http.Http_actions[action.Name] = action
		
	}

	fmt.Printf("Loaded %d actions for http\n", len(actions))



}	

var Actions_data_fs="data"
var Http_Actions_fs_path=Actions_data_fs+"/http"
func (this *Actions_obj) CheckHttpActionsFS() (bool, error) {
	_, err := os.Stat(Http_Actions_fs_path)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Directory doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	return true, nil
}


func (this *Actions_obj) ReadHttpActionsFiles() ([]Http_actions, error) {
	var actionsList []Http_actions

	for _, executorType := range Http_executor_types {
		reading_path := Http_Actions_fs_path + executorType

		entries, err := os.ReadDir(reading_path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			// if entry.Name() == "tmp.go" {
			// 	continue
			// }

			actionName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

			action := Http_actions{
				Name:          actionName,
				Path:          reading_path,
				Executor_type: executorType,
				Executor_function_name:"StartAction_"+actionName,
			}
			fmt.Printf("Founded Executor Type: %s | Action Name: %s\n for Http", action.Executor_type, action.Name)
			actionsList = append(actionsList, action)
		}
	}

	return actionsList, nil
}




type WS_Actions_obj struct {
	WS_actions map[string]WsAction
	Config map[string]map[string]Executor_function_signature

}

var Ws_executor_types = []string{
	"Send",
	"Connect",
	"Recieve",
}
type WsAction struct {
	Name string
	Path string
	Exector_type string
	Executor_function_name string
}

var Ws_Actions_fs_path=Actions_data_fs+"/ws"
func (this *Actions_obj) CheckWsActionsFS() (bool, error) {
	_, err := os.Stat(Ws_Actions_fs_path)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Directory doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	return true, nil
}

func (this *Actions_obj) ReadWsActionsFiles() ([]WsAction, error) {
	var actionsList []WsAction

	for _, executorType := range Ws_executor_types {
		reading_path := Ws_Actions_fs_path + executorType

		entries, err := os.ReadDir(reading_path)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			actionName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

			action := WsAction{
				Name:          actionName,
				Path:          reading_path,
				Exector_type:  executorType,
				Executor_function_name:      "",
			}

			fmt.Printf("Founded Executor Type: %s | Action Name: %s\n for ws", action.Exector_type, action.Name)

			actionsList = append(actionsList, action)
		}
	}

	return actionsList, nil
}






type ActionType string

const (
	ActionTypeHTTP ActionType = "http"
	ActionTypeWS   ActionType = "ws"
)

type Action_opts struct {
	Action_type ActionType
	Opts any
}

type HTTP_opts struct {
	Executor_type string
	Method string
	URL    string
}

type WS_opts struct {
	Executor_type string
	// URL   string
}




func (this *Actions_obj) CreateAction(action_name string, action_opts Action_opts) bool {
	switch action_opts.Action_type {
	case ActionTypeHTTP:
		httpOpts, ok := action_opts.Opts.(HTTP_opts)
		if !ok {
			fmt.Println("Invalid HTTP options")
			return false
		}

		if !ValidHTTP_exec_type(httpOpts) {
			fmt.Println("Invalid HTTP executor type:", httpOpts.Executor_type)
			return false
		}

		if _, exists := this.Http.Http_actions[action_name]; exists {
			return false
		}

		var http_action = Http_actions{
			Name:                   action_name,
			Executor_type:          httpOpts.Executor_type,
			Path:                   "",
			Executor_function_name: "",
		}

		http_action.CreateHttpAction()

		this.Http.Http_actions[action_name] = http_action
		this.Http.Config["http"][action_name] = nil

		this.Http.CreateCodeHTTPobject()

		return true

	case ActionTypeWS:
		wsOpts, ok := action_opts.Opts.(WS_opts)
		if !ok {
			fmt.Println("Invalid WS options")
			return false
		}

		if !ValidWS_exec_type(wsOpts) {
			fmt.Println("Invalid WS executor type:", wsOpts.Executor_type)
			return false
		}

		if _, exists := this.WS.WS_actions[action_name]; exists {
			return false
		}

		wsAction := WsAction{
			Name:                   action_name,
			Path:                   "",
			Exector_type:            wsOpts.Executor_type,
			Executor_function_name: "",
		}

		switch wsOpts.Executor_type {
		case "Send":
			wsAction.CreateWsSendFile()

		case "Connect":
			wsAction.CreateWsConnectFile()

		case "Recieve":
			wsAction.CreateWsReceiveFile()

		default:
			fmt.Println("Invalid WS executor type:", wsOpts.Executor_type)
			return false
		}

		

		this.WS.WS_actions[action_name] = wsAction
		this.WS.Config["ws"][action_name] = nil

		this.WS.CreateCodeWSobject()

		return true

	default:
		fmt.Println("Invalid action type")
		return false
	}
}

func ValidHTTP_exec_type(httpOpts HTTP_opts) bool {
	for _, executorType := range Http_executor_types {
		if httpOpts.Executor_type == executorType {
			return true
		}
	}

	return false
}

func ValidWS_exec_type(wsOpts WS_opts)bool{
	validExecutorType := false

	for _, executorType := range Ws_executor_types {
		if wsOpts.Executor_type == executorType {
			validExecutorType = true
			break
		}
	}

	if !validExecutorType {
		fmt.Println("Invalid WS executor type:", wsOpts.Executor_type)
		return false
	}

	return true

}


func (this *Http_actions)CreateHttpAction(){
	var action_file_name=Http_Actions_fs_path+"/"+this.Executor_type+this.Name
	this.Path=action_file_name
	switch this.Executor_type {
	case "req":
		this.SetUpHttpActionFile_req_type()
		this.Executor_function_name="StartAction_"+this.Name
	}

	

}

func (this *Http_actions) SetUpHttpActionFile_req_type() {
	code := `
	package data_actions

	import (
		"bytes"
		"encoding/json"
		"fmt"
		"io"
		"net/http"
		"tester/data"
	)

	type Action_name_Body struct {
		
	}

	type Action_name_Cookies struct {
		
	}

	func StartAction_Action_name(client *http.Client, prevReq *data.ActionResponse,sequnces *data.SequnceState) *data.ActionResponse {
		body := Action_name_Body{
			
		}

		cookies := Action_name_Cookies{
			
		}

		response, err := Action_name_SendReq(client, body, cookies)
		if err != nil {
			fmt.Printf("Failed to get appropriate response, passing empty struct and the error is %v\n", err)
			return &data.ActionResponse{}
		}

		return response
	}

	func Action_name_SendReq(client *http.Client, body any, cookies Action_name_Cookies) (*data.ActionResponse, error) {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		req, err := http.NewRequest("POST", "https://example.com", bytes.NewReader(bodyBytes))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		req.AddCookie(&http.Cookie{
			
		})

		

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return &data.ActionResponse{
			Body:       string(responseBody),
			StatusCode: resp.StatusCode,
			Status:     resp.Status,
			Headers:    resp.Header,
			Cookies:    resp.Cookies(),
			Data:       make(map[string]any),
		}, nil
	}		`


	code = strings.ReplaceAll(code, "Action_name", this.Name)
	
	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating action file:", err)
		return
	}

	fmt.Println("Action file created:", this.Path)
}



var actions_config_file_path="Actions/actions_config.go"
func UpdateConfigFile(actionType string, code string) {
	content, err := os.ReadFile(actions_config_file_path)
	if err != nil {
		fmt.Println("Error reading actions config:", err)
		return
	}

	contentString := string(content)

	start := fmt.Sprintf("\t\"%s\": {", actionType)
	startIndex := strings.Index(contentString, start)

	if startIndex == -1 {
		fmt.Println("Action type not found in actions config:", actionType)
		return
	}

	endIndex := strings.Index(contentString[startIndex:], "\n\t},")
	if endIndex == -1 {
		fmt.Println("Action type block not found in actions config:", actionType)
		return
	}

	endIndex = startIndex + endIndex + len("\n\t},")

	newBlock := fmt.Sprintf("\t\"%s\": {\n%s\t},", actionType, code)

	contentString = contentString[:startIndex] + newBlock + contentString[endIndex:]

	err = os.WriteFile(actions_config_file_path, []byte(contentString), 0644)
	if err != nil {
		fmt.Println("Error updating actions config:", err)
		return
	}
}


func (this *HTTP_Actions_obj) CreateCodeHTTPobject() {
	var code string

	for actionName, action := range this.Http_actions {
		code += fmt.Sprintf("\t\"%s\": %s,\n", actionName, action.Executor_function_name)
	}

	UpdateConfigFile("http",code)
}


func (this *WS_Actions_obj) CreateCodeWSobject() {
	var code string

	for actionName, action := range this.WS_actions {
		code += fmt.Sprintf("\t\"%s\": %s,\n", actionName, action.Executor_function_name)
	}

	UpdateConfigFile("ws",code)
}


func (this *Actions_obj) DeleteHttpAction(actionName string) bool {
	action, exists := this.Http.Http_actions[actionName]

	if !exists {
		fmt.Println("HTTP action not found:", actionName)
		return false
	}

	err := action.DeleteActionFile()
	if err != nil {
		fmt.Println("Failed to delete action file:", err)
		return false
	}

	delete(this.Http.Http_actions, actionName)
	delete(this.Http.Config["http"], actionName)

	this.Http.CreateCodeHTTPobject()

	return true
}

func (this *Actions_obj) DeleteWsAction(actionName string) bool {
	action, exists := this.WS.WS_actions[actionName]

	if !exists {
		fmt.Println("WS action not found:", actionName)
		return false
	}

	err := action.DeleteActionFile()
	if err != nil {
		fmt.Println("Failed to delete action file:", err)
		return false
	}

	delete(this.WS.WS_actions, actionName)
	delete(this.WS.Config["ws"], actionName)

	this.WS.CreateCodeWSobject()

	return true
}

func (this *Http_actions) DeleteActionFile() error {
	filePath := this.Path + this.Name + ".go"

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Action file deleted:", filePath)

	return nil
}

func (this *WsAction) DeleteActionFile() error {
	filePath := this.Path + this.Name + ".go"

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Action file deleted:", filePath)

	return nil
}





func (this *WsAction) CreateWsConnectFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
)

func StartAction_Action_name(url string) (*websocket.Conn, error) {
	return nil, nil
}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket connect action file:", err)
		return
	}

	this.Executor_function_name = "StartAction_" + this.Name

	fmt.Println("WebSocket connect action file created:", this.Path+this.Name+".go")
}

func (this *WsAction) CreateWsSendFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
	"tester/data"
)

func WS_send_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {

	return nil
}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket send action file:", err)
		return
	}

	this.Executor_function_name = "WS_send_" + this.Name

	fmt.Println("WebSocket send action file created:", this.Path+this.Name+".go")
}

func (this *WsAction) CreateWsReceiveFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
	"tester/data"
)

func WS_handle_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {

	return nil
}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket receive action file:", err)
		return
	}

	this.Executor_function_name = "WS_handle_" + this.Name

	fmt.Println("WebSocket receive action file created:", this.Path+this.Name+".go")
}