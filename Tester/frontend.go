package main

import (
	actions "Tester/backend/Actions"
	app "Tester/backend/App"
	sequences "Tester/backend/Sequences"
	"Tester/backend/models"
)



func (a *App) GetApplication(app_id int64) (models.ApplicationFrontend, error) {
    application, err := a.App.GetApplicationById(app_id)
    if err != nil {
        return models.ApplicationFrontend{}, err
    }

    actions, err := a.Action.GetActions(a.Db,application.Actions)
    if err != nil {
        return models.ApplicationFrontend{}, err
    }

    sequences, err := a.Sequence.GetSequences(application.Sequences)
    if err != nil {
        return models.ApplicationFrontend{}, err
    }

    children, err := a.App.GetApplicationsByIds(application.Children)
    if err != nil {
        return models.ApplicationFrontend{}, err
    }

    return models.ApplicationFrontend{
        AppId:         application.AppId,
        AppName:       application.AppName,
        ParentAppId:   application.ParentAppId,
        EnvironmentId: application.EnvironmentId,
        Actions:       ConvertActions(actions),
        Sequences:     ConvertSequences(sequences),
        Children:      ConvertChildren(children),
    }, nil
}

func ConvertActions(actions []actions.Action) []models.ActionInfo {
    result := make([]models.ActionInfo, 0, len(actions))

    for _, action := range actions {
        result = append(result, models.ActionInfo{
            Id:   int64(action.ActionId),
            Name: action.ActionType,
        })
    }

    return result
}

func ConvertSequences(sequences []sequences.Sequence) []models.SequenceInfo {
    result := make([]models.SequenceInfo, 0, len(sequences))

    for _, sequence := range sequences {
        result = append(result, models.SequenceInfo{
            Id:   sequence.SequenceID,
            Name: sequence.SequenceName,
        })
    }

    return result
}

func ConvertChildren(children []app.Application) []models.ChildApplicationInfo {
    result := make([]models.ChildApplicationInfo, 0, len(children))

    for _, child := range children {
        result = append(result, models.ChildApplicationInfo{
            Id:   child.AppId,
            Name: child.AppName,
        })
    }

    return result
}