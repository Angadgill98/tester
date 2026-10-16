package Actions

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type HTTPOpts struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Query   map[string]string `json:"query"`
	Body    []byte            `json:"body"`
}

type HTTPAction struct {
	ActionID   int64     `json:"action_id"`
	AppID      int64     `json:"app_id"`
	ActionName string    `json:"action_name"`
	Opts       HTTPOpts  `json:"opts"`
}

// func CreateHTTPAction(actionID int64, appID int64, actionName string, opts HTTPOpts) HTTPAction {
// 	return HTTPAction{
// 		ActionID:   actionID,
// 		AppID:      appID,
// 		ActionName: actionName,
// 		Opts:       opts,
// 	}
// }


func InsertHTTPAction(db *sql.DB, appID int64, actionName string, opts HTTPOpts) (int64, error) {
	optsJSON, err := json.Marshal(opts)
	if err != nil {
		return 0, fmt.Errorf("marshal HTTP action options: %w", err)
	}

	result, err := db.Exec(
		`INSERT INTO http_actions (app_id, action_name, opts) VALUES (?, ?, ?)`,
		appID,
		actionName,
		optsJSON,
	)
	
	if err != nil {
		return 0, fmt.Errorf("insert HTTP action: %w", err)
	}

	actionID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get HTTP action ID: %w", err)
	}

	return actionID, nil
}

func HTTPActionExists(db *sql.DB, appID int64, actionName string) (bool, error) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM http_actions
			WHERE app_id = ? AND action_name = ?
		)`,
		appID, actionName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check HTTP action existence: %w", err)
	}

	return exists, nil
}


func GetHTTPAction(db *sql.DB, actionID int64, appID int64) (HTTPAction, error) {
	var action HTTPAction
	var optsJSON []byte

	err := db.QueryRow(
		`SELECT id, app_id, action_name, opts
		 FROM http_actions
		 WHERE id = ? AND app_id = ?`,
		actionID, appID,
	).Scan(&action.ActionID, &action.AppID, &action.ActionName, &optsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return HTTPAction{}, fmt.Errorf("HTTP action %d not found in application %d", actionID, appID)
		}
		return HTTPAction{}, fmt.Errorf("get HTTP action: %w", err)
	}

	if err := json.Unmarshal(optsJSON, &action.Opts); err != nil {
		return HTTPAction{}, fmt.Errorf("parse HTTP action options: %w", err)
	}

	return action, nil
}

func GetHTTPActionsByAppID(db *sql.DB, appID int64) ([]HTTPAction, error) {
	rows, err := db.Query(
		`SELECT id, app_id, action_name, opts
		 FROM http_actions
		 WHERE app_id = ?`,
		appID,
	)
	if err != nil {
		return nil, fmt.Errorf("get HTTP actions: %w", err)
	}
	defer rows.Close()

	var actions []HTTPAction

	for rows.Next() {
		var action HTTPAction
		var optsJSON []byte

		if err := rows.Scan(&action.ActionID, &action.AppID, &action.ActionName, &optsJSON); err != nil {
			return nil, fmt.Errorf("scan HTTP action: %w", err)
		}

		if err := json.Unmarshal(optsJSON, &action.Opts); err != nil {
			return nil, fmt.Errorf("parse HTTP action options: %w", err)
		}

		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate HTTP actions: %w", err)
	}

	return actions, nil
}

func GetHTTPActionIDByName(db *sql.DB, appID int64, actionName string) (int64, error) {
	var actionID int64

	err := db.QueryRow(
		`SELECT id FROM http_actions WHERE app_id = ? AND action_name = ?`,
		appID, actionName,
	).Scan(&actionID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("HTTP action %q not found in application %d", actionName, appID)
		}
		return 0, fmt.Errorf("get HTTP action ID by name: %w", err)
	}

	return actionID, nil
}

func UpdateHTTPActionOpts(db *sql.DB, actionID int64, appID int64, opts json.RawMessage) error {
    if !json.Valid(opts) {
        return fmt.Errorf("invalid HTTP action options JSON")
    }

    query := `UPDATE http_actions SET opts = $1 WHERE id = $2 AND app_id = $3`

    result, err := db.Exec(query, opts, actionID, appID)
    if err != nil {
        return fmt.Errorf("failed to update HTTP action options: %w", err)
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

