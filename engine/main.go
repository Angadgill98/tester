package main

import (
	"engine/Actions"
	"engine/Application"
	"engine/Sequences"
	"fmt"
)

func main() {
	applicationObj := Application.CreateApplicationObj()
	applicationObj.SetUpApplicationObj()

	if _, exists := (*applicationObj)["GLOBAL"]; !exists {
		(*applicationObj)["GLOBAL"] = applicationObj.CreateApplication("GLOBAL")
	}

	actionsObj,err := Actions.CreateActionsObj()
	if err != nil {
		fmt.Println("Failed to set up actions:", err)
		return
	}

	sequencesObj, err := Sequences.CreateSequenceObj()
	if err != nil {
		fmt.Println("Failed to set up sequences:", err)
		return
	}

	_ = sequencesObj
	_ = actionsObj

}