package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
)



type HttpAction struct{
	ActionName string
	AppId int64
	ActionId int64 
	Opts HttpActionOpts
	EnvironmentId int64 
}	

type HttpActionOpts struct{
	Method  string            
	URL     string            
	Headers map[string]string 
	Query   map[string]string 
	Body    json.RawMessage            
	Cookies map[string]string
}	


func CreateHttpAction(db *sql.DB, action_name string, app_id int64, opts HttpActionOpts, environment_id int64) (int64, error) {
    optsJSON, err := json.Marshal(opts)
    if err != nil {
        return 0, err
    }

    query, values := BuildInsertQuery("http_actions", map[string]any{
        "action_name":    action_name,
        "app_id":         app_id,
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


func GetHttpAction(db *sql.DB, action_id int64) (HttpAction, error) {
    var action HttpAction
    var optsJSON string

    query, values := BuildSelectQuery(
        "http_actions",
        []string{"action_id", "app_id", "action_name", "opts", "environment_id"},
        map[string][]any{
            "action_id": {"=", action_id},
        },
		
    )

    err := db.QueryRow(query, values...).Scan(
        &action.ActionId,
        &action.AppId,
        &action.ActionName,
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


func UpdateHTTPAction(db *sql.DB, actionID int64, actionName string, appID int64, opts json.RawMessage, environmentID int64) error {
    if !json.Valid(opts) {
        return fmt.Errorf("invalid HTTP action options JSON")
    }

    query, values := BuildUpdateQuery(
        "http_actions",
        map[string]any{
            "action_name":    actionName,
            "app_id":         appID,
            "opts":           string(opts),
            "environment_id": environmentID,
        },
        map[string][]any{
            "action_id": {"=", actionID},
        },
    )

    result, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to update HTTP action: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to check updated HTTP action rows: %w", err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("HTTP action not found")
    }

    return nil
}

func DeleteHTTPAction(db *sql.DB, actionID int64) error {
    query, values := BuildDeleteQuery("http_actions", map[string][]any{
        "action_id": {"=", actionID},
    })

    _, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to delete HTTP action: %w", err)
    }

    return nil
}