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




func SendData(conn net.Conn, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Failed to serialize JSON:", err)
		return err
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		fmt.Println("Failed to send JSON:", err)
		return err
	}

	return nil
}

