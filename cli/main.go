package main

import "fmt"



func main(){



	socket := ConnectUnixSocketWithRetries()
	if socket == nil {
		fmt.Println("CLI: TestFlow engine socket is not initialized")
		return
	}

	cli:=CreateCLI(&socket)

	defer socket.Close()
}


func Init(){

}


func SetGlobalApp(){
	
}