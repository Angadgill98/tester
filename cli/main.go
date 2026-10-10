package main

import (
	"fmt"
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


