package main

import (
	"fmt"
	"os"
	"os/exec"
)

func StartEngine() *exec.Cmd {
	cmd := exec.Command("./engine")
	cmd.Dir = global.EnginePath

	err := cmd.Start()
	if err != nil {
		fmt.Println("Failed to start engine:", err)
		return nil
	}

	return cmd
}

func KillEngine(pid int) {
	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Println("Failed to find engine process:", err)
		return
	}

	err = process.Kill()
	if err != nil {
		fmt.Println("Failed to kill engine process:", err)
		return
	}
}