package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
)



func main(){

	err:=Init()
	if err!=nil{
		return
	}

	socket := ConnectUnixSocketWithRetries()
	if socket == nil {
		fmt.Println("CLI: TestFlow engine socket is not initialized")
		return
	}

	if err := ReceiveGlobalAppID(socket); err != nil {
		fmt.Println(err)
		return
	}

	cli:=CreateCLI(&socket)

	cli.StartCLILoop()

	defer socket.Close()
}


func Init() error {
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = global.EnginePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	return nil
}




func ReceiveGlobalAppID(conn net.Conn) error {
    var response Response

    if err := json.NewDecoder(conn).Decode(&response); err != nil {
        return fmt.Errorf("CLI: Failed to receive GLOBAL application: %w", err)
    }

    if !response.Status {
        return fmt.Errorf("CLI: Failed to initialize GLOBAL application: %s", response.Msg)
    }

    if len(response.Args) < 2 {
        return fmt.Errorf("CLI: GLOBAL application name or ID missing from response")
    }

    appName, ok := response.Args[0].(string)
    if !ok || appName == "" {
        return fmt.Errorf("CLI: Invalid GLOBAL application name")
    }

    appID, ok := response.Args[1].(string)
    if !ok || appID == "" {
        return fmt.Errorf("CLI: Invalid GLOBAL application ID")
    }

    global.CurrentAppPath = []string{appName}
    global.CurrentAppID = appID

    return nil
}