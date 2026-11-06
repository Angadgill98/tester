package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
)






type WSOpts struct {
	URL     string           
	
}

type WSAction struct {
	ActionID    int64  
	AppID       int64  
	ActionName  string 
	ExecutorType string 
	Opts        WSOpts 
	EnvironmentId int64 
}


func CreateWSAction(db *sql.DB, app_id int64, action_name string, executor_type string, opts WSOpts, environment_id int64) (int64, error) {
    optsJSON, err := json.Marshal(opts)
    if err != nil {
        return 0, err
    }

    query, values := BuildInsertQuery("ws_actions", map[string]any{
        "app_id":         app_id,
        "action_name":    action_name,
        "executor_type":  executor_type,
        "opts":           string(optsJSON),
        "environment_id": environment_id,
    }, []string{"action_id"})

    var actionID int64
    err = db.QueryRow(query, values...).Scan(&actionID)
    if err != nil {
        return 0, err
    }

    return actionID, nil
}

func GetWSAction(db *sql.DB, action_id int64) (WSAction, error) {
    var action WSAction
    var optsJSON string

    query, values := BuildSelectQuery(
        "ws_actions",
        []string{"action_id", "app_id", "action_name", "executor_type", "opts", "environment_id"},
        map[string][]any{
            "action_id": {"=", action_id},
        },
    )

    err := db.QueryRow(query, values...).Scan(
        &action.ActionID,
        &action.AppID,
        &action.ActionName,
        &action.ExecutorType,
        &optsJSON,
        &action.EnvironmentId,
    )
    if err != nil {
        return action, err
    }

    err = json.Unmarshal([]byte(optsJSON), &action.Opts)
    if err != nil {
        return action, err
    }

    return action, nil
}

func UpdateWSAction(db *sql.DB, actionID int64, appID int64, actionName string, executorType string, opts json.RawMessage, environmentID int64) error {
    if !json.Valid(opts) {
        return fmt.Errorf("invalid WebSocket action options JSON")
    }

    query, values := BuildUpdateQuery(
        "ws_actions",
        map[string]any{
            "app_id":         appID,
            "action_name":    actionName,
            "executor_type":  executorType,
            "opts":           string(opts),
            "environment_id": environmentID,
        },
        map[string][]any{
            "action_id": {"=", actionID},
        },
    )

    result, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to update WebSocket action: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check updated WebSocket action rows: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("WebSocket action not found")
    }

    return nil
}

func DeleteWSAction(db *sql.DB, actionID int64) error {
    query, values := BuildDeleteQuery("ws_actions", map[string][]any{
        "action_id": {"=", actionID},
    })

    _, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to delete WebSocket action: %w", err)
    }

    return nil
}