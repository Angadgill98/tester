package Actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type WSOpts struct {
	URL     string            `json:"url"`
	
}

type WSAction struct {
	ActionID    int64  `json:"action_id"`
	AppID       int64  `json:"app_id"`
	ActionName  string `json:"action_name"`
	ExecutorType string `json:"executor_type"`
	Opts        WSOpts `json:"opts"`
}

// func CreateWSAction(actionID int64, appID int64, actionName string, executorType string, opts WSOpts) WSAction {
// 	return WSAction{
// 		ActionID:     actionID,
// 		AppID:        appID,
// 		ActionName:   actionName,
// 		ExecutorType: executorType,
// 		Opts:         opts,
// 	}
// }

func InsertWSAction(db *sql.DB, appID int64, actionName string, executorType string, opts WSOpts) (int64, error) {
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		return 0, fmt.Errorf("marshal WebSocket action options: %w", err)
	}

	result, err := db.Exec(
		`INSERT INTO ws_actions (app_id, action_name, executor_type, opts) VALUES (?, ?, ?, ?)`,
		appID,
		actionName,
		executorType,
		optsJSON,
	)
	if err != nil {
		return 0, fmt.Errorf("insert WebSocket action: %w", err)
	}

	actionID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get WebSocket action ID: %w", err)
	}

	return actionID, nil
}

func WSActionExists(db *sql.DB, appID int64, actionName string) (bool, error) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM ws_actions
			WHERE app_id = ? AND action_name = ?
		)`,
		appID, actionName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check WebSocket action existence: %w", err)
	}

	return exists, nil
}


func GetWSAction(db *sql.DB, actionID int64, appID int64) (WSAction, error) {
	var action WSAction
	var optsJSON []byte

	err := db.QueryRow(
		`SELECT id, app_id, action_name, executor_type, opts
		 FROM ws_actions
		 WHERE id = ? AND app_id = ?`,
		actionID, appID,
	).Scan(&action.ActionID, &action.AppID, &action.ActionName, &action.ExecutorType, &optsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return WSAction{}, fmt.Errorf("WebSocket action %d not found in application %d", actionID, appID)
		}
		return WSAction{}, fmt.Errorf("get WebSocket action: %w", err)
	}

	if err := json.Unmarshal(optsJSON, &action.Opts); err != nil {
		return WSAction{}, fmt.Errorf("parse WebSocket action options: %w", err)
	}

	return action, nil
}

func GetWSActionsByAppID(db *sql.DB, appID int64) ([]WSAction, error) {
	rows, err := db.Query(
		`SELECT id, app_id, action_name, executor_type, opts
		 FROM ws_actions
		 WHERE app_id = ?`,
		appID,
	)
	if err != nil {
		return nil, fmt.Errorf("get WebSocket actions: %w", err)
	}
	defer rows.Close()

	var actions []WSAction

	for rows.Next() {
		var action WSAction
		var optsJSON []byte

		if err := rows.Scan(&action.ActionID, &action.AppID, &action.ActionName, &action.ExecutorType, &optsJSON); err != nil {
			return nil, fmt.Errorf("scan WebSocket action: %w", err)
		}

		if err := json.Unmarshal(optsJSON, &action.Opts); err != nil {
			return nil, fmt.Errorf("parse WebSocket action options: %w", err)
		}

		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate WebSocket actions: %w", err)
	}

	return actions, nil
}

func GetWSActionIDByName(db *sql.DB, appID int64, actionName string) (int64, error) {
	var actionID int64

	err := db.QueryRow(
		`SELECT id FROM ws_actions WHERE app_id = ? AND action_name = ?`,
		appID, actionName,
	).Scan(&actionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("WebSocket action %q not found in application %d", actionName, appID)
		}
		return 0, fmt.Errorf("get WebSocket action ID by name: %w", err)
	}

	return actionID, nil
}

func UpdateWSActionOpts(db *sql.DB, actionID int64, appID int64, opts json.RawMessage) error {
    if !json.Valid(opts) {
        return fmt.Errorf("invalid WebSocket action options JSON")
    }

    query := `UPDATE ws_actions SET opts = $1 WHERE id = $2 AND app_id = $3`

    result, err := db.Exec(query, opts, actionID, appID)
    if err != nil {
        return fmt.Errorf("failed to update WebSocket action options: %w", err)
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