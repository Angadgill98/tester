package main

import (

	"github.com/gorilla/websocket"
	"tester/data"
)

func StartAction_Action_name(url string) (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func WS_handle_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return &data.ActionResponse{}
	}

	return &data.ActionResponse{
		Body:   string(message),
		Status: "received",
		Data:   make(map[string]any),
	}
}

func WS_send_Action_name(conn *websocket.Conn, prevReq *data.ActionResponse, sequnces *data.SequnceState) *data.ActionResponse {
	err := conn.WriteMessage(websocket.TextMessage, []byte(prevReq.Body))
	if err != nil {
		return &data.ActionResponse{}
	}

	return &data.ActionResponse{}
}