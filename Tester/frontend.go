package main

import (
	actions "Tester/backend/Actions"
	app "Tester/backend/App"
	sequences "Tester/backend/Sequences"
	"Tester/backend/models"
	"encoding/json"
	"time"
)


type Response struct {
    Success bool
    Data    any
    Error   string
}


func (a *App) GetApplication(app_id int64) Response {
    application, err := a.App.GetApplicationById(app_id)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    actions, err := a.Action.GetActions(a.Db, application.Actions)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    httpIDs := []int64{}
    wsIDs := []int64{}

    for _, action := range actions {
        switch action.ActionType {
        case "http":
            httpIDs = append(httpIDs, int64(action.ActionId))
        case "ws":
            wsIDs = append(wsIDs, int64(action.ActionId))
        }
    }

    httpActions, err := a.Action.GetHttpActions(a.Db, httpIDs)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    wsActions, err := a.Action.GetWSActions(a.Db, wsIDs)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    httpMap := make(map[int64]string)
    for _, action := range httpActions {
        httpMap[action.ActionId] = action.ActionName
    }

    wsMap := make(map[int64]string)
    for _, action := range wsActions {
        wsMap[action.ActionID] = action.ActionName
    }

    actionNames := make([]models.ActionInfo, 0, len(actions))

    for _, action := range actions {
        var name string

        switch action.ActionType {
        case "http":
            name = httpMap[int64(action.ActionId)]
        case "ws":
            name = wsMap[int64(action.ActionId)]
        }

        actionNames = append(actionNames, models.ActionInfo{
            Id:   int64(action.ActionId),
            Name: name,
            Type: action.ActionType,
        })
    }

    sequences, err := a.Sequence.GetSequences(application.Sequences)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    children, err := a.App.GetApplicationsByIds(application.Children)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    return Response{
        Success: true,
        Data: map[string]any{
            "Application": application,
            "Actions":     actionNames,
            "Sequences":   ConvertSequences(sequences),
            "Children":    ConvertChildren(children),
        },
        Error: "",
    }
}

func (a *App) GetApplications(app_ids []int64) Response {
    appList, err := a.App.GetApplicationsByIds(app_ids)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    results := make([]map[string]any, 0, len(appList))

    for _, application := range appList {
        actions, err := a.Action.GetActions(a.Db, application.Actions)
        if err != nil {
            return Response{Success: false, Error: err.Error()}
        }

        httpIDs := []int64{}
        wsIDs := []int64{}

        for _, action := range actions {
            switch action.ActionType {
            case "http":
                httpIDs = append(httpIDs, int64(action.ActionId))
            case "ws":
                wsIDs = append(wsIDs, int64(action.ActionId))
            }
        }

        httpActions, err := a.Action.GetHttpActions(a.Db, httpIDs)
        if err != nil {
            return Response{Success: false, Error: err.Error()}
        }

        wsActions, err := a.Action.GetWSActions(a.Db, wsIDs)
        if err != nil {
            return Response{Success: false, Error: err.Error()}
        }

        httpMap := make(map[int64]string)
        for _, action := range httpActions {
            httpMap[action.ActionId] = action.ActionName
        }

        wsMap := make(map[int64]string)
        for _, action := range wsActions {
            wsMap[action.ActionID] = action.ActionName
        }

        actionNames := make([]models.ActionInfo, 0, len(actions))

        for _, action := range actions {
            var name string

            switch action.ActionType {
            case "http":
                name = httpMap[int64(action.ActionId)]
            case "ws":
                name = wsMap[int64(action.ActionId)]
            }

            actionNames = append(actionNames, models.ActionInfo{
                Id:   int64(action.ActionId),
                Name: name,
                Type: action.ActionType,
            })
        }

        sequences, err := a.Sequence.GetSequences(application.Sequences)
        if err != nil {
            return Response{Success: false, Error: err.Error()}
        }

        children, err := a.App.GetApplicationsByIds(application.Children)
        if err != nil {
            return Response{Success: false, Error: err.Error()}
        }

        results = append(results, map[string]any{
            "Application": application,
            "Actions":     actionNames,
            "Sequences":   ConvertSequences(sequences),
            "Children":    ConvertChildren(children),
        })
    }

    return Response{
        Success: true,
        Data:    results,
        Error:   "",
    }
}

func ConvertActions(actions []actions.Action) []models.ActionInfo {
    result := make([]models.ActionInfo, 0, len(actions))

    for _, action := range actions {
        result = append(result, models.ActionInfo{
            Id:   int64(action.ActionId),
            Name: action.ActionType,
			Type: action.ActionType,
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

func (a *App) CreateApplication(app_name string, parentAppId int64) Response {
    app, err := a.App.CreateApplication(app_name, parentAppId)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    err = a.App.RegisterChild(parentAppId, app.AppId)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    return Response{
        Success: true,
        Data:    app,
        Error:   "",
    }
}

func (a *App) CreateHttpAction(action_name string, action_type string, app_id int64) Response {
    defaultOpts := actions.HttpActionOpts{
        Method:  "GET",
        URL:     "",
        Headers: map[string]string{},
        Query:   map[string]string{},
        Body:    json.RawMessage(`{}`),
        Cookies: map[string]string{},
    }

    environmentid := 0

    action_id, err := a.Action.CreateAction(a.Db, action_type)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    err = a.Action.CreateHttpAction(a.Db,action_id,action_name,app_id,defaultOpts,int64(environmentid))
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    err = a.App.RegisterAction(app_id, action_id)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    return Response{
        Success: true,
        Data:    action_id,
        Error:   "",
    }
}

func (a *App) CreateWSAction(action_name string, action_type string, app_id int64) Response {
    defaultOpts := actions.WSOpts{
        URL:              "",
        Headers:          map[string]string{},
        QueryParams:      map[string]string{},
        Subprotocols:     []string{},
        Messages:         []actions.WSMessage{},
        HandshakeTimeout: 5 * time.Second,
        WriteTimeout:     10 * time.Second,
    }

    environmentid := 0

    action_id, err := a.Action.CreateAction(a.Db, action_type)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    err = a.Action.CreateWSAction(a.Db, action_id, app_id, action_name, "", defaultOpts, int64(environmentid))
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    err = a.App.RegisterAction(app_id, action_id)
    if err != nil {
        return Response{Success: false, Error: err.Error()}
    }

    return Response{
        Success: true,
        Data:    action_id,
        Error:   "",
    }
}

func (a *App) GetHttpAction(action_id int64) Response {
    action, err := a.Action.GetHttpAction(a.Db, action_id)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    return Response{
        Success: true,
        Data:    action,
        Error:   "",
    }
}

func (a *App) GetWSAction(action_id int64) Response {
    action, err := a.Action.GetWSAction(a.Db, action_id)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    return Response{
        Success: true,
        Data:    action,
        Error:   "",
    }
}

func (a *App) GetSequence(sequence_id int64) Response {
    action, err := a.Sequence.GetSequence(sequence_id)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    return Response{
        Success: true,
        Data:    action,
        Error:   "",
    }
}

func (a *App) CreateSequence(sequence_name string, app_id int64) Response {
    sequence, err := a.Sequence.CreateSequence(sequence_name, app_id)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    err = a.App.RegisterSequence(app_id, sequence.SequenceID)
    if err != nil {
        return Response{
            Success: false,
            Error:   err.Error(),
        }
    }

    return Response{
        Success: true,
        Data:    sequence,
        Error:   "",
    }
}