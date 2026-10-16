package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)




func ConnectUnixSocketWithRetries() net.Conn {
	for i := 0; i < global.SocketRetries; i++ {
		conn, err := net.Dial("unix", "/tmp/testflow.sock")
		if err == nil {
			return conn
		}

		time.Sleep(time.Second)
	}

	return nil
}




func SendData(conn net.Conn, data any) (Response, error) {
	var response Response

	jsonData, err := json.Marshal(data)
	if err != nil {
		return response, fmt.Errorf("failed to serialize JSON: %w", err)
	}

	if _, err := conn.Write(jsonData); err != nil {
		return response, fmt.Errorf("failed to send JSON: %w", err)
	}

	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		return response, fmt.Errorf("failed to read response: %w", err)
	}


	if len(response.Args) > 0 {
		fmt.Println("Args:")
		for _, arg := range response.Args {
			fmt.Printf("  %+v\n", arg)
		}
	}

	return response, nil
}
