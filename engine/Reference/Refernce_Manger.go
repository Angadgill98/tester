package Refernce

// import (
// 	"engine/Actions"
// 	"engine/Application"
// 	"engine/Sequences"
// )

// type ReferenceManager struct {
// 	References map[string]any
// }

// func CreateReferenceManager() *ReferenceManager {
// 	referenceManager := ReferenceManager{
// 		References: make(map[string]any),
// 	}

// 	return &referenceManager
// }


// func (rm *ReferenceManager) SetUpRefrenceManager(app *Application.Application) {
	

// 	rm.SetUpActionReferences(app.Sequences)
// 	rm.SetUpSequnceApplicationReferences(app, app.Sequences)
// }


// func (rm *ReferenceManager) SetUpActionReferences(sequenceObj *Sequences.Sequence_obj) {
// 	rm.References["actions"] = make(map[string]ActionReference)

// 	rm.SetActionsExectorRefernce(rm)
// 	rm.SetUpActionsSequenceReferences(sequenceObj)
// }

// type ActionReference struct {
// 	Function    any
// 	ReferencedBy []Reference
// }

// type Reference struct {
// 	Type string
// 	Name string
// }

// func (rm *ReferenceManager)SetActionsExectorRefernce(referenceManager *ReferenceManager) {
// 	config := Actions.Actions_config
// 	actionsMap := referenceManager.References["actions"].(map[string]ActionReference)

// 	for _, actions := range config {
// 		for action, function := range actions {
// 			actionsMap[action] = ActionReference{
// 				Function:     function,
// 				ReferencedBy: []Reference{},
// 			}
// 		}
// 	}
// }


// func (rm *ReferenceManager) SetUpActionsSequenceReferences(sequenceObj *Sequences.Sequence_obj) {
// 	actionsMap := rm.References["actions"].(map[string]ActionReference)

// 	for sequenceName, sequences := range *sequenceObj {
// 		for sequenceIndex := range sequences {
// 			sequence := &sequences[sequenceIndex]

// 			validActions := []Sequences.SequenceAction{}

// 			for _, sequenceAction := range sequence.Actions {
// 				var actionName string

// 				if sequenceAction.Action_type == "http" {
// 					actionName = sequenceAction.HTTP_action.Name
// 				}

// 				if sequenceAction.Action_type == "ws" {
// 					actionName = sequenceAction.WS_action.Name
// 				}

// 				actionReference, exists := actionsMap[actionName]
// 				if !exists {
// 					continue
// 				}

// 				actionReference.ReferencedBy = append(actionReference.ReferencedBy, Reference{
// 					Type: "sequence",
// 					Name: sequenceName,
// 				})

// 				actionsMap[actionName] = actionReference

// 				validActions = append(validActions, sequenceAction)
// 			}

// 			sequence.Actions = validActions
// 		}

// 		(*sequenceObj)[sequenceName] = sequences
// 	}
// }


// func (rm *ReferenceManager) SetUpSequnceApplicationReferences(application *Application.Application, sequenceObj *Sequences.Sequence_obj) {
// 	rm.References["sequences"] = make(map[string][]Reference)

// 	sequencesMap := rm.References["sequences"].(map[string][]Reference)


// 		validSequences := make(map[string][]Sequences.Sequence)

// 		for sequenceName, sequences := range *application.Sequences {
// 			_, exists := (*sequenceObj)[sequenceName]
// 			if !exists {
// 				continue
// 			}

// 			validSequences[sequenceName] = sequences

// 			sequencesMap[sequenceName] = append(sequencesMap[sequenceName], Reference{
// 				Type: "application",
// 				Name: sequenceName,
// 			})
// 		}

// 		*application.Sequences = validSequences
// }



