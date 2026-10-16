package main

import "encoding/json"

type Response struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
	Args   []any  `json:"args"`
}

type CreateAppParams struct {
	Command     string `json:"command"`
	Operation   string `json:"operation"`
	Name        string `json:"name"`
	ParentAppID int64  `json:"parent_app_id"`
}

type CreateActionParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	Type      string `json:"type"` // http or ws
	Name      string `json:"name"`
	AppID     int64  `json:"app_id"`
	Opts      json.RawMessage `json:"opts"`
}














type CreateSequenceParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	Name      string `json:"name"`
	AppID     int64  `json:"app_id"`
}

type AddSequenceParams struct {
	Command      string `json:"command"`
	Operation    string `json:"operation"`
	AppID        int64  `json:"app_id"`
	SequenceName string `json:"sequence_name"`
	ActionName   string `json:"action_name"`
	Type         string `json:"type"` // http or ws
	Position     int    `json:"position"`
}

type RemoveSequenceParams struct {
	Command      string `json:"command"`
	Operation    string `json:"operation"`
	AppID        int64  `json:"app_id"`
	SequenceName string `json:"sequence_name"`
	Position     int    `json:"position"`
}












type CdAppParams struct {
	Command      string `json:"command"`
	Operation    string `json:"operation"`
	ChildAppName string `json:"child_app_name"`
	ParentAppID  int64  `json:"parent_app_id"`
}

type CdActionParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	AppID     int64  `json:"app_id"`
	Action_name  string  `json:"action_name"`
	Type      string `json:"type"` // "http" or "ws"
}












type LsDefaultParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	AppID     int64  `json:"app_id"`
}

type LsActionsParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	AppID     int64  `json:"app_id"`
}

type LsSequencesParams struct {
	Command   string `json:"command"`
	Operation string `json:"operation"`
	AppID     int64  `json:"app_id"`
}










type ExecuteActionParams struct {
	Command    string `json:"command"`
	Operation  string `json:"operation"`
	AppID      int64  `json:"app_id"`
	ActionID   int64  `json:"action_id"`
	ActionName string `json:"action_name"`
}

type ExecuteSequenceParams struct {
	Command    string `json:"command"`
	Operation  string `json:"operation"`
	AppID      int64  `json:"app_id"`
	SequenceID int64  `json:"sequence_id"`
}



type SaveActionParams struct {
	Command    string          `json:"command"`
	Operation  string          `json:"operation"`
	AppID      int64           `json:"app_id"`
	ActionName string          `json:"action_name"`
	Type       string          `json:"type"` // "http" or "ws"
	Opts       json.RawMessage `json:"opts"`
}