package main

import (
	"database/sql"
	"encoding/json"
	"engine/Actions"
	application "engine/Application"
	sequences "engine/Sequences"
	Global "engine/global"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
)

type Socket struct {
	DB        *sql.DB

}

func  CreateSocket(db *sql.DB) *Socket {
	return &Socket{
		DB:          db,
	}
}

func (socket *Socket) StartSocket() {
	socketPath := Global.GlobalConfig.SocketPath

	if err := os.MkdirAll(filepath.Dir(socketPath), 0755); err != nil {
		fmt.Printf("Failed to create socket directory: %v\n", err)
		return
	}

	if err := os.Remove(socketPath); err != nil && !os.IsNotExist(err) {
		fmt.Printf("Failed to remove existing socket: %v\n", err)
		return
	}

	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Printf("Failed to bind socket: %v\n", err)
		return
	}
	defer listener.Close()
	defer os.Remove(socketPath)

	fmt.Printf("Listening on socket: %s\n", socketPath)

	if err := socket.ListenForClients(listener); err != nil {
		fmt.Printf("Socket listener stopped: %v\n", err)
	}


}

func (socket *Socket) ListenForClients(listener net.Listener) error {
	conn, err := listener.Accept()
	if err != nil {
		return fmt.Errorf("failed to accept client: %w", err)
	}

	socket.HandleClient(conn)
	return nil
}

func (socket *Socket) HandleClient(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Client connected: %s\n", conn.RemoteAddr())

	buffer := make([]byte, 4096)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Failed to read from client: %v\n", err)
			}
			break
		}

		if n > 0 {
			response := socket.HandleCommand(buffer[:n])

			if err := socket.SendResponse(conn, response); err != nil {
				fmt.Printf("Failed to send response: %v\n", err)
			}
		}
	}
}

func (socket *Socket) HandleCommand(data []byte) Response {
	command, err := GetJSONField(data, "command")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	commandStr, ok := command.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: command must be a string"}
	}

	switch commandStr {

	case "create":
		return socket.ProcessCreate(data)
	case "sequence":
		return socket.ProcessSequence(data)
	case "execute":
		return socket.ProcessExecute(data)
	case "cd":
		return socket.ProcessCd(data)
	case "ls":
		return socket.ProcessLs(data)
	default:
		return Response{
			Status: false,
			Msg:    fmt.Sprintf("Engine: Unknown command: %v", command),
		}
	}
}



func (socket *Socket) SendResponse(conn net.Conn, responseObj Response) error {
	data, err := json.Marshal(responseObj)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}

func (socket *Socket) ProcessCreate(data []byte) Response {
	operation, err := GetJSONField(data, "operation")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	op, ok := operation.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: operation must be a string"}
	}

	switch op {
	case "app":
		var params CreateAppParams

		if err := json.Unmarshal(data, &params); err != nil {
			return Response{
				Status: false,
				Msg:    fmt.Sprintf("Engine: failed to parse create app parameters: %v", err),
			}
		}
		if params.ParentAppID != 0 {
			isChild, err := application.IsChildApplication(socket.DB, params.ParentAppID, params.Name)
			if err != nil {
				return Response{Status: false, Msg: err.Error()}
			}

			if isChild {
				return Response{Status: false, Msg: fmt.Sprintf("Application %q already exists under parent ID %d", params.Name, params.ParentAppID)}
			}
		}

		appID, err := application.InsertApplication(socket.DB, params.Name, params.ParentAppID)
		if err != nil {
			return Response{Status: false, Msg: err.Error()}
		}

		return Response{
			Status: true,
			Msg:    fmt.Sprintf("Application created successfully with ID %d", appID),
		}


	case "action":
		var params CreateActionParams

		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse create action parameters: %v", err)}
		}

		switch params.Type {
		case "http":
			exists, err := Actions.HTTPActionExists(socket.DB, params.AppID, params.Name)
			if err != nil {
				return Response{Status: false, Msg: err.Error()}
			}
			if exists {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: HTTP action %q already exists in app %d", params.Name, params.AppID)}
			}

			var opts Actions.HTTPOpts
			if err := json.Unmarshal(params.Opts, &opts); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse HTTP action options: %v", err)}
			}

			actionID, err := Actions.InsertHTTPAction(socket.DB, params.AppID, params.Name, opts)
			if err != nil {
				return Response{Status: false, Msg: err.Error()}
			}

			app, err := application.GetApplication(socket.DB, params.AppID)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get application: %v", err)}
			}

			app.Actions = append(app.Actions, actionID)

			if err := application.SaveApplication(socket.DB, app); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save application: %v", err)}
			}

			return Response{Status: true, Msg: "HTTP action created successfully", Args: []any{actionID}}

		case "ws":
			exists, err := Actions.WSActionExists(socket.DB, params.AppID, params.Name)
			if err != nil {
				return Response{Status: false, Msg: err.Error()}
			}
			if exists {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: WebSocket action %q already exists in app %d", params.Name, params.AppID)}
			}

			var opts Actions.WSOpts
			if err := json.Unmarshal(params.Opts, &opts); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse WebSocket action options: %v", err)}
			}

			actionID, err := Actions.InsertWSAction(socket.DB, params.AppID, params.Name, "ws", opts)
			if err != nil {
				return Response{Status: false, Msg: err.Error()}
			}

			app, err := application.GetApplication(socket.DB, params.AppID)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get application: %v", err)}
			}

			app.Actions = append(app.Actions, actionID)

			if err := application.SaveApplication(socket.DB, app); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save application: %v", err)}
			}


			return Response{Status: true, Msg: "WebSocket action created successfully", Args: []any{actionID}}

		default:
			return Response{Status: false, Msg: "Engine: invalid action type. Expected http or ws"}
		}



	case "sequence":
		var params CreateSequenceParams

		if err := json.Unmarshal(data, &params); err != nil {
			return Response{
				Status: false,
				Msg:    fmt.Sprintf("Engine: failed to parse create sequence parameters: %v", err),
			}
		}

		exists, err := sequences.SequenceExists(socket.DB, params.AppID, params.Name)
		if err != nil {
			return Response{Status: false, Msg: err.Error()}
		}

		if exists {
			return Response{
				Status: false,
				Msg:    fmt.Sprintf("Engine: sequence %q already exists in app %d", params.Name, params.AppID),
			}
		}

		sequenceID, err := sequences.InsertSequence(socket.DB, params.AppID, params.Name)
		if err != nil {
			return Response{Status: false, Msg: err.Error()}
		}

		app, err := application.GetApplication(socket.DB, params.AppID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get application: %v", err)}
		}

		app.Sequences = append(app.Sequences, sequenceID)

		if err := application.SaveApplication(socket.DB, app); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save application: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Sequence created successfully",
			Args:   []any{sequenceID},
		}

		
	default:
		return Response{Status: false, Msg: "Engine: Invalid create operation"}
	}

	return Response{Status: true}
}

func (socket *Socket) ProcessSequence(data []byte) Response {
	operation, err := GetJSONField(data, "operation")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	op, ok := operation.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: operation must be a string"}
	}

	switch op {
	case "add":
		var params AddSequenceParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse add sequence parameters: %v", err)}
		}

		sequence, err := sequences.GetSequenceByName(socket.DB, params.AppID, params.SequenceName)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get sequence: %v", err)}
		}

		var actionID int64

		switch params.Type {
		case "http":
			actionID, err = Actions.GetHTTPActionIDByName(socket.DB, params.AppID, params.ActionName)
		case "ws":
			actionID, err = Actions.GetWSActionIDByName(socket.DB, params.AppID, params.ActionName)
		default:
			return Response{Status: false, Msg: "Engine: invalid action type. Expected http or ws"}
		}

		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get action: %v", err)}
		}

		if params.Position < 0 || params.Position > len(sequence.ActionIDs) {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: invalid insertion position %d", params.Position)}
		}

		sequence.ActionIDs = append(sequence.ActionIDs, 0)
		copy(sequence.ActionIDs[params.Position+1:], sequence.ActionIDs[params.Position:])
		sequence.ActionIDs[params.Position] = actionID

		if err := sequences.SaveSequence(socket.DB, sequence); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save sequence: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Engine: action added to sequence successfully",
			Args:   []any{sequence.SequenceID, actionID},
		}

	case "remove":
		var params RemoveSequenceParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse remove sequence parameters: %v", err)}
		}

		sequence, err := sequences.GetSequenceByName(socket.DB, params.AppID, params.SequenceName)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get sequence: %v", err)}
		}

		if params.Position < 0 || params.Position >= len(sequence.ActionIDs) {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: invalid action index %d", params.Position)}
		}

		removedActionID := sequence.ActionIDs[params.Position]
		sequence.ActionIDs = append(sequence.ActionIDs[:params.Position], sequence.ActionIDs[params.Position+1:]...)

		if err := sequences.SaveSequence(socket.DB, sequence); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save sequence: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Engine: action removed from sequence successfully",
			Args:   []any{sequence.SequenceID, removedActionID},
		}
	case "reorder":

	case "exchange":

	default:
		return Response{
			Status: false,
			Msg:    "Engine: Invalid sequence operation. Expected add, remove, reorder, exchange, or delete",
	}
}

	return Response{Status: true}
}

func (socket *Socket) ProcessCd(data []byte) Response {
	operation, err := GetJSONField(data, "operation")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	op, ok := operation.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: operation must be a string"}
	}

	switch op {
	case "app":
		var params CdAppParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse cd app parameters: %v", err)}
		}

		appID, err := application.GetApplicationByName(socket.DB, params.ParentAppID, params.ChildAppName)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to find application %q: %v", params.ChildAppName, err)}
		}

		app, err := application.GetApplication(socket.DB, appID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get application: %v", err)}
		}

		return Response{Status: true, Msg: "Engine: application retrieved successfully", Args: []any{app}}

	case "action":
		var params CdActionParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse cd action parameters: %v", err)}
		}

		switch params.Type {
		case "http":
			actionID, err := Actions.GetHTTPActionIDByName(socket.DB, params.AppID, params.Action_name)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to find HTTP action %q: %v", params.Action_name, err)}
			}

			action, err := Actions.GetHTTPAction(socket.DB, actionID, params.AppID)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get HTTP action: %v", err)}
			}

			return Response{
				Status: true,
				Msg:    "Engine: HTTP action retrieved successfully",
				Args:   []any{"http", action.Opts},
			}

		case "ws":
			actionID, err := Actions.GetWSActionIDByName(socket.DB, params.AppID, params.Action_name)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to find WebSocket action %q: %v", params.Action_name, err)}
			}

			action, err := Actions.GetWSAction(socket.DB, actionID, params.AppID)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get WebSocket action: %v", err)}
			}

			return Response{
				Status: true,
				Msg:    "Engine: WebSocket action retrieved successfully",
				Args:   []any{"ws", action.Opts},
			}

		default:
			return Response{Status: false, Msg: fmt.Sprintf("Engine: invalid action type %q; expected http or ws", params.Type)}
		}

	case "save":
		var params SaveActionParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse save action parameters: %v", err)}
		}

		if params.Type != "http" && params.Type != "ws" {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: invalid action type %q; expected http or ws", params.Type)}
		}

		if len(params.Opts) == 0 || !json.Valid(params.Opts) {
			return Response{Status: false, Msg: "Engine: action options must contain valid JSON"}
		}

		switch params.Type {
		case "http":
			actionID, err := Actions.GetHTTPActionIDByName(socket.DB, params.AppID, params.ActionName)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to find HTTP action %q: %v", params.ActionName, err)}
			}

			// Implement this call using your actual HTTP action update function.
			if err := Actions.UpdateHTTPActionOpts(socket.DB, actionID, params.AppID, params.Opts); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save HTTP action options: %v", err)}
			}

		case "ws":
			actionID, err := Actions.GetWSActionIDByName(socket.DB, params.AppID, params.ActionName)
			if err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to find WebSocket action %q: %v", params.ActionName, err)}
			}

			// Implement this call using your actual WebSocket action update function.
			if err := Actions.UpdateWSActionOpts(socket.DB, actionID, params.AppID, params.Opts); err != nil {
				return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to save WebSocket action options: %v", err)}
			}
		}

		return Response{Status: true, Msg: "Engine: action options saved successfully"}

	default:
		return Response{Status: false, Msg: "Engine: invalid cd operation. Expected app, action, or save"}
	}
}

func (socket *Socket) ProcessLs(data []byte) Response {
	operation, err := GetJSONField(data, "operation")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	op, ok := operation.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: operation must be a string"}
	}

	switch op {
		case "":
		var params LsDefaultParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse ls parameters: %v", err)}
		}

		children, err := application.GetChildApplications(socket.DB, params.AppID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get child applications: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Engine: applications retrieved successfully",
			Args:   []any{children},
		}

	case "actions":
		var params LsActionsParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse ls actions parameters: %v", err)}
		}

		httpActions, err := Actions.GetHTTPActionsByAppID(socket.DB, params.AppID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get HTTP actions: %v", err)}
		}

		wsActions, err := Actions.GetWSActionsByAppID(socket.DB, params.AppID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get WebSocket actions: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Engine: actions retrieved successfully",
			Args:   []any{httpActions, wsActions},
		}

	case "sequences":
		var params LsSequencesParams
		if err := json.Unmarshal(data, &params); err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to parse ls sequences parameters: %v", err)}
		}

		sequenceList, err := sequences.GetSequencesByAppID(socket.DB, params.AppID)
		if err != nil {
			return Response{Status: false, Msg: fmt.Sprintf("Engine: failed to get sequences: %v", err)}
		}

		return Response{
			Status: true,
			Msg:    "Engine: sequences retrieved successfully",
			Args:   []any{sequenceList},
		}

	default:
		return Response{Status: false, Msg: "Engine: Invalid ls operation. Expected empty, actions, or sequences"}
	}

	return Response{Status: true}
}

func (socket *Socket) ProcessExecute(data []byte) Response {
	operation, err := GetJSONField(data, "operation")
	if err != nil {
		return Response{Status: false, Msg: err.Error()}
	}

	op, ok := operation.(string)
	if !ok {
		return Response{Status: false, Msg: "Engine: operation must be a string"}
	}

	switch op {
	case "action":
	case "sequence":
	default:
		return Response{Status: false, Msg: "Engine: Invalid execute operation. Expected action or sequence"}
	}

	return Response{Status: true}
}