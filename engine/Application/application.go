package application

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Application struct {
	AppId       int64
	AppName     string
	ParentAppId int64
	Actions     []int64
	Sequences   []int64
	Children    []int64
}

// func CreateApplication(appId int64, appName string, parentAppId int64, actions []int64, sequences []int64, children []int64) Application {
// 	return Application{
// 		AppId:       appId,
// 		AppName:     appName,
// 		ParentAppId: parentAppId,
// 		Actions:     actions,
// 		Sequences:   sequences,
// 		Children:    children,
// 	}
// }



func SaveApplication(db *sql.DB, app Application) error {
	actions, err := json.Marshal(app.Actions)
	if err != nil {
		return fmt.Errorf("marshal action IDs: %w", err)
	}

	sequences, err := json.Marshal(app.Sequences)
	if err != nil {
		return fmt.Errorf("marshal sequence IDs: %w", err)
	}

	children, err := json.Marshal(app.Children)
	if err != nil {
		return fmt.Errorf("marshal child application IDs: %w", err)
	}

	var parentAppID any
	if app.ParentAppId != 0 {
		parentAppID = app.ParentAppId
	}

	result, err := db.Exec(
		`UPDATE apps
		SET name = ?, parent_app_id = ?, action_ids = ?, sequence_ids = ?, children = ?
		WHERE id = ?`,
		app.AppName, parentAppID, string(actions), string(sequences), string(children), app.AppId,
	)
	if err != nil {
		return fmt.Errorf("save application: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check saved application: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("application %d not found", app.AppId)
	}

	return nil
}

func InsertApplication(db *sql.DB, appName string, parentAppID int64) (int64, error) {
	var parentID any
	if parentAppID != 0 {
		var exists bool
		err := db.QueryRow(`SELECT EXISTS (SELECT 1 FROM apps WHERE id = ?)`, parentAppID).Scan(&exists)
		if err != nil {
			return 0, fmt.Errorf("check parent application: %w", err)
		}
		if !exists {
			return 0, fmt.Errorf("parent application %d not found", parentAppID)
		}

		parentID = parentAppID
	}

	result, err := db.Exec(
		`INSERT INTO apps (name, parent_app_id) VALUES (?, ?)`,
		appName, parentID,
	)
	if err != nil {
		return 0, fmt.Errorf("insert application: %w", err)
	}

	appID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get application ID: %w", err)
	}

	return appID, nil
}

func GetApplication(db *sql.DB, appID int64) (Application, error) {
	var app Application
	var parentAppID sql.NullInt64
	var actionsJSON string
	var sequencesJSON string
	var childrenJSON string

	err := db.QueryRow(
		`SELECT id, name, parent_app_id, action_ids, sequence_ids, children
		FROM apps WHERE id = ?`,
		appID,
	).Scan(
		&app.AppId, &app.AppName, &parentAppID,
		&actionsJSON, &sequencesJSON, &childrenJSON,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return Application{}, fmt.Errorf("application %d not found", appID)
		}
		return Application{}, fmt.Errorf("get application: %w", err)
	}

	if parentAppID.Valid {
		app.ParentAppId = parentAppID.Int64
	}

	if err := json.Unmarshal([]byte(actionsJSON), &app.Actions); err != nil {
		return Application{}, fmt.Errorf("parse action IDs: %w", err)
	}

	if err := json.Unmarshal([]byte(sequencesJSON), &app.Sequences); err != nil {
		return Application{}, fmt.Errorf("parse sequence IDs: %w", err)
	}

	if err := json.Unmarshal([]byte(childrenJSON), &app.Children); err != nil {
		return Application{}, fmt.Errorf("parse child application IDs: %w", err)
	}

	return app, nil
}

func IsChildApplication(db *sql.DB, parentAppID int64, childAppName string) (bool, error) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1
			FROM apps AS parent, json_each(parent.children) AS child
			JOIN apps AS child_app ON child_app.id = child.value
			WHERE parent.id = ? AND child_app.name = ?
		)`,
		parentAppID, childAppName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check child application: %w", err)
	}

	return exists, nil
}

func GetChildApplications(db *sql.DB, parentAppID int64) ([]Application, error) {
	rows, err := db.Query(
		`SELECT id, name, parent_app_id, action_ids, sequence_ids, children
		 FROM apps
		 WHERE parent_app_id = ?`,
		parentAppID,
	)
	if err != nil {
		return nil, fmt.Errorf("get child applications: %w", err)
	}
	defer rows.Close()

	var children []Application

	for rows.Next() {
		var app Application
		var parentID sql.NullInt64
		var actionsJSON, sequencesJSON, childrenJSON string

		if err := rows.Scan(&app.AppId, &app.AppName, &parentID, &actionsJSON, &sequencesJSON, &childrenJSON); err != nil {
			return nil, fmt.Errorf("scan child application: %w", err)
		}

		if parentID.Valid {
			app.ParentAppId = parentID.Int64
		}

		if err := json.Unmarshal([]byte(actionsJSON), &app.Actions); err != nil {
			return nil, fmt.Errorf("parse child application action IDs: %w", err)
		}
		if err := json.Unmarshal([]byte(sequencesJSON), &app.Sequences); err != nil {
			return nil, fmt.Errorf("parse child application sequence IDs: %w", err)
		}
		if err := json.Unmarshal([]byte(childrenJSON), &app.Children); err != nil {
			return nil, fmt.Errorf("parse nested child application IDs: %w", err)
		}

		children = append(children, app)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate child applications: %w", err)
	}

	return children, nil
}

func GetApplicationByName(db *sql.DB, parentAppID int64, appName string) (int64, error) {
	var appID int64

	err := db.QueryRow(
		`SELECT id FROM apps WHERE name = ? AND parent_app_id = ?`,
		appName, parentAppID,
	).Scan(&appID)

	if err != nil {
		return 0, err
	}

	return appID, nil
}