package actions

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tester/data"
)



type Actions_obj struct{
	Actions map[string]Action
	Config map[string]func(*http.Client, *data.ActionResponse, *data.SequnceState) *data.ActionResponse
	WS_congif  map[string]map[string]any
}

func CreateActionsObj() Actions_obj {
	return Actions_obj{
		Actions: make(map[string]Action),
		Config:  ActionsConfig,
		WS_congif: WSConfig,
	}
}


var ActionFS_path="data/actions/"
func (this *Actions_obj) SetupActions() {
	isok, err := this.CheckActionsFS()

	if err != nil {
		return
	}

	if !isok {
		fmt.Println("Actions FS not Setup")
		return
	}

	if isok {
		fmt.Println("Actions FS Setup OK")
	}

	fmt.Println("Getting Actions files")

	actions, err := this.ReadActionsFiles()
	if err != nil {
		fmt.Println("Failed to read Actions files:", err)
		return
	}

	for _, action := range actions {
		this.Actions[action.Name] = action
		
	}

	fmt.Printf("Loaded %d actions\n", len(actions))

	
}



func (this *Actions_obj) CheckActionsFS() (bool, error) {
	_, err := os.Stat(ActionFS_path)

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

	_, err = os.Stat("Actions/" + "actions_config.go")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("actions_config doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied for actions_config")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	wsPath := ActionFS_path + "ws"

	_, err = os.Stat(wsPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("WebSocket actions directory doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied for WebSocket actions directory")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	_, err = os.Stat("Actions/" + "actions_ws_config.go")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("actions_ws_config doesn't exist")
			return false, err
		}

		if os.IsPermission(err) {
			fmt.Println("Permission denied for actions_ws_config")
			return false, err
		}

		fmt.Println("Other error:", err)
		return false, err
	}

	return true, nil
}




func (this *Actions_obj) ReadActionsFiles() ([]Action, error) {
	entries, err := os.ReadDir(ActionFS_path)
	if err != nil {
		return nil, err
	}

	var actionsList []Action

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "tmp.go" {
			continue
		}

		fmt.Printf("In Action FS found the action File:%v\n", entry.Name())

		actionName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))

		action := Action{
			Name: actionName,
			Path: ActionFS_path,
			Executor: this.Config[actionName],
		}

		actionsList = append(actionsList, action)
	}

	return actionsList, nil
}



type Action struct{
	Name string
	Path string
	Executor func(*http.Client, *data.ActionResponse, *data.SequnceState) *data.ActionResponse

}
type ActionType string

const (
	ActionTypeHTTP ActionType = "http"
	ActionTypeWS   ActionType = "ws"
)

type HTTP_opts struct {
	Method string
	URL    string
}

type WS_opts struct {
	URL   string
	Event string
}

type Action_opts struct {
	Action_type ActionType
	HTTP        *HTTP_opts
	WS          *WS_opts
}

func (this *Actions_obj) CreateAction(action_name string,opts Action_opts) bool {
	if _, exists := this.Actions[action_name]; exists {
		fmt.Println("Action already exists:", action_name)
		return false
	}

	

	switch opts.Action_type {
	case "http":
		action := Action{
			Name: action_name,
			Path: fmt.Sprintf("%v", ActionFS_path),
			Executor: nil,
		}
		fmt.Println("Creating HTTP action")

		fmt.Println("Setitng up HTTP Action File")
		action.SetUpHttpActionFile();
		fmt.Println("Updating actions_config.go")
		err := this.UpdateConfig(action)
		if err != nil {
			fmt.Println("Failed to update actions config:", err)
			return false
		}

	case "ws":
		fmt.Printf("Creating WS action with event %v \n", opts.WS.Event)
		action := Action{
			Name: action_name,
			Path: fmt.Sprintf("%vws/", ActionFS_path),
			Executor: nil,
		}

		switch opts.WS.Event {
			case "connect":
				fmt.Println("Creating WS connect action")
				action.CreateWsConnectFile()

			case "send":
				fmt.Println("Creating WS send action")
				action.CreateWsSendFile()

			case "receive":
				fmt.Println("Creating WS receive action")
				action.CreateWsReceiveFile()

			default:
				fmt.Println("Invalid WS event")
			}

		
	default:
		fmt.Println("Invalid action type:", opts.Action_type)
		return false
	}

	


	fmt.Println("Actions created and Added to config but not loaded it memory as not USABLE")
	fmt.Println("Use SHOW for more info")
	// this.Actions[action_name] = action
	this.Config[action_name]=nil

	return true
}

func (this *Action) SetUpHttpActionFile() {
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


var file_path="Actions"
func (this *Actions_obj) UpdateConfig(action Action) error {
	configPath := file_path + "/actions_config.go"

	code, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	functionName := "data_actions.StartAction_"+action.Name  

	entry := fmt.Sprintf("\t\"%s\": %s,\n", action.Name, functionName)

	configCode := string(code)

	insertPosition := strings.LastIndex(configCode, "}")

	if insertPosition == -1 {
		return fmt.Errorf("invalid actions_config.go")
	}

	configCode = configCode[:insertPosition] + entry + configCode[insertPosition:]

	err = os.WriteFile(configPath, []byte(configCode), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Actions config updated:", functionName)

	return nil
}

func (this *Actions_obj) DeleteAction(actionName string) bool {
	action, exists := this.Actions[actionName]

	if exists {
		err := action.DeleteActionFile()
		if err != nil {
			fmt.Println("Failed to delete action file:", err)
			return false
		}
	} else {
		err := this.DeleteActionFile(actionName)
		if err != nil {
			fmt.Println("Failed to delete action file:", err)
			return false
		}
	}

	err := this.DeleteConfig(actionName)
	if err != nil {
		fmt.Println("Failed to remove action from config:", err)
		return false
	}

	delete(this.Actions, actionName)
	delete(this.Config, actionName)

	return true
}

func (this *Actions_obj) DeleteActionFile(actionName string) error {
	filePath := ActionFS_path + actionName + ".go"

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Action file deleted:", filePath)

	return nil
}

func (this *Action) DeleteActionFile() error {
	filePath := this.Path + this.Name + ".go"

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	fmt.Println("Action file deleted:", filePath)

	return nil
}

func (this *Actions_obj) DeleteConfig(actionName string) error {
	configPath := file_path + "/actions_config.go"

	code, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	functionName := "data_actions.StartAction_" + actionName

	entry := fmt.Sprintf("\t\"%s\": %s,\n", actionName, functionName)

	configCode := strings.Replace(string(code), entry, "", 1)

	err = os.WriteFile(configPath, []byte(configCode), 0644)
	if err != nil {
		return err
	}

	fmt.Println("Action removed from config:", actionName)

	return nil
}






func (this *Action) CreateWsConnectFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
	"tester/data"
)

func StartAction_Action_name(url string) (*websocket.Conn, error) {

}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket connect action file:", err)
		return
	}

	fmt.Println("WebSocket connect action file created:", this.Path)
}

func (this *Action) CreateWsSendFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
	"tester/data"
)

func WS_send_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {

}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket send action file:", err)
		return
	}

	fmt.Println("WebSocket send action file created:", this.Path)
}

func (this *Action) CreateWsReceiveFile() {
	code := `
package data_ws_actions

import (
	"github.com/gorilla/websocket"
	"tester/data"
)

func WS_handle_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {

}
`

	code = strings.ReplaceAll(code, "Action_name", this.Name)

	err := os.WriteFile(this.Path+this.Name+".go", []byte(code), 0644)
	if err != nil {
		fmt.Println("Error creating WebSocket receive action file:", err)
		return
	}

	fmt.Println("WebSocket receive action file created:", this.Path)
}
