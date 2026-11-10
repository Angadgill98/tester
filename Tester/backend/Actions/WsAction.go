package actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)


type WSOpts struct {
    URL               string
    Headers           map[string]string
    QueryParams       map[string]string
    Subprotocols      []string
    Messages          []WSMessage
    HandshakeTimeout  time.Duration
    WriteTimeout      time.Duration
}

type WSMessage struct {
    MessageID   int64
    Sequence    int
    MessageType string
    Payload     json.RawMessage
}



type WSAction struct {
    ActionID      int64
    AppID         int64
    ActionName    string
    Opts          WSOpts
    EnvironmentID int64
}



func (a *ActionObj) CreateWSAction(db *sql.DB, action_id int64, app_id int64, action_name string, executor_type string, opts WSOpts, environment_id int64) error {
    optsJSON, err := json.Marshal(opts)
    if err != nil {
        return err
    }

    query, values := BuildInsertQuery("ws_actions", map[string]any{
        "action_id":      action_id,
        "app_id":         app_id,
        "action_name":    action_name,
        "opts":           string(optsJSON),
        "environment_id": environment_id,
    }, []string{})

    _, err = db.Exec(query, values...)
    return err
}

func (a *ActionObj)GetWSAction(db *sql.DB, action_id int64) (WSAction, error) {
    var action WSAction
    var optsJSON string

    query, values := BuildSelectQuery(
        "ws_actions",
        []string{"action_id", "app_id", "action_name", "opts", "environment_id"},
        map[string][]any{
            "action_id": {"=", action_id},
        },
    )

    err := db.QueryRow(query, values...).Scan(
        &action.ActionID,
        &action.AppID,
        &action.ActionName,
        &optsJSON,
        &action.EnvironmentID,
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

func (a *ActionObj) GetWSActions(db *sql.DB, ids []int64) ([]WSAction, error) {
    if len(ids) == 0 {
        return []WSAction{}, nil
    }

    placeholders := strings.TrimRight(strings.Repeat("?,", len(ids)), ",")
    query := fmt.Sprintf(
        "SELECT action_id, app_id, action_name, opts, environment_id FROM ws_actions WHERE action_id IN (%s)",
        placeholders,
    )

    args := make([]any, len(ids))
    for i, id := range ids {
        args[i] = id
    }

    rows, err := db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var actions []WSAction

    for rows.Next() {
        var action WSAction
        var optsJSON string

        err := rows.Scan(
            &action.ActionID,
            &action.AppID,
            &action.ActionName,
            &optsJSON,
            &action.EnvironmentID,
        )
        if err != nil {
            return nil, err
        }

        if err := json.Unmarshal([]byte(optsJSON), &action.Opts); err != nil {
            return nil, err
        }

        actions = append(actions, action)
    }

    return actions, rows.Err()
}

func (a *ActionObj)UpdateWSAction(db *sql.DB, actionID int64, appID int64, actionName string, executorType string, opts json.RawMessage, environmentID int64) error {
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

func (a *ActionObj)DeleteWSAction(db *sql.DB, actionID int64) error {
    query, values := BuildDeleteQuery("ws_actions", map[string][]any{
        "action_id": {"=", actionID},
    })

    _, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to delete WebSocket action: %w", err)
    }

    return nil
}