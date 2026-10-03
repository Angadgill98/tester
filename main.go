package main

import (
	"fmt"
	Actions "tester/Actions"
	"tester/CLI"
	sequences "tester/Sequences"
	Services "tester/Services"
)

func main() {
	var services=Services.CreateServices()
	fmt.Println("========== File System Setup START ==========")
	services.File_system.SetUpFiles()
	fmt.Printf("========== File System Setup END ========== \n\n\n\n\n")

	var actions = Actions.CreateActionsObj()

	fmt.Println("========== Actions Setup START ==========")
	actions.SetUpActionsObj()
	fmt.Printf("========== Actions Setup End ==========\n\n\n\n\n")

	var sequences=sequences.CreateSequencesObj()

	fmt.Println("========== Sequence Setup START ==========")
	sequences.SetupSequences(actions)
	fmt.Printf("========== Sequence Setup START ==========\n\n\n\n\n")



	cli := CLI.CreateCLI(&actions,&sequences)

	cli.Init()


}














type Tester struct {
}

func CreateTesterObj() {

}

func SetUpTester(tester Tester) {
	fmt.Println("Setting up tester");



}