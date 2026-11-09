package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)



type Application struct {
	AppId         int64
	AppName       string
	ParentAppId   int64
	EnvironmentId sql.NullInt64
	Actions       []int64
	Sequences     []int64
	Children      []int64
}

type ApplicationObj struct{
	Db *sql.DB
}

func CreateApplicationObject(db *sql.DB)*ApplicationObj{	
	return &ApplicationObj{
		Db:db,
	}
}


func (a *ApplicationObj) CreateApplication(appName string, parentAppId int64) (Application, error) {
	actions, _ := json.Marshal([]int64{})
	sequences, _ := json.Marshal([]int64{})
	children, _ := json.Marshal([]int64{})

	data := map[string]any{
		"name":         appName,
		"parent_app_id": parentAppId,
		"action_ids":   string(actions),
		"sequence_ids": string(sequences),
		"children":     string(children),
	}

	query, values := BuildInsertQuery("applications", data, []string{"id"})

	application := Application{
		AppName:     appName,
		ParentAppId: parentAppId,
		Actions:     []int64{},
		Sequences:   []int64{},
		Children:    []int64{},
	}

	err := a.Db.QueryRow(query, values...).Scan(&application.AppId)
	if err != nil {
		return Application{}, err
	}

	return application, nil
}

func (a *ApplicationObj) GetApplicationById(appId int64) (Application, error) {
	query, values := BuildSelectQuery(
		"applications",
		[]string{"id", "name", "parent_app_id", "environment_id", "action_ids", "sequence_ids", "children"},
		map[string][]any{
			"id": {"=", appId},
		},
	)

	var application Application
	var actions, sequences, children string

	err := a.Db.QueryRow(query, values...).Scan(
		&application.AppId,
		&application.AppName,
		&application.ParentAppId,
		&application.EnvironmentId,
		&actions,
		&sequences,
		&children,
	)
	if err != nil {
		return Application{}, err
	}

	if err := json.Unmarshal([]byte(actions), &application.Actions); err != nil {
		return Application{}, err
	}

	if err := json.Unmarshal([]byte(sequences), &application.Sequences); err != nil {
		return Application{}, err
	}

	if err := json.Unmarshal([]byte(children), &application.Children); err != nil {
		return Application{}, err
	}

	return application, nil
}

func (a *ApplicationObj) GetApplicationsByIds(appIds []int64) ([]Application, error) {
	if len(appIds) == 0 {
		return []Application{}, nil
	}

	placeholders := make([]string, len(appIds))
	values := make([]any, len(appIds))

	for i, id := range appIds {
		placeholders[i] = "?"
		values[i] = id
	}

	query := fmt.Sprintf(
		"SELECT id, name, parent_app_id, environment_id, action_ids, sequence_ids, children FROM applications WHERE id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	rows, err := a.Db.Query(query, values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applications := []Application{}

	for rows.Next() {
		var application Application
		var actions, sequences, children string

		err := rows.Scan(
			&application.AppId,
			&application.AppName,
			&application.ParentAppId,
			&application.EnvironmentId,
			&actions,
			&sequences,
			&children,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(actions), &application.Actions); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(sequences), &application.Sequences); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(children), &application.Children); err != nil {
			return nil, err
		}

		applications = append(applications, application)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return applications, nil
}


func BuildInsertQuery(table string, data map[string]any, returningColumns []string) (string, []any) {
    columns := []string{}
    placeholders := []string{}
    values := []any{}

    for column, value := range data {
        columns = append(columns, column)
        placeholders = append(placeholders, "?")
        values = append(values, value)
    }

    query := fmt.Sprintf(
        "INSERT INTO %s (%s) VALUES (%s)",
        table,
        strings.Join(columns, ", "),
        strings.Join(placeholders, ", "),
    )

    if len(returningColumns) > 0 {
        query += " RETURNING " + strings.Join(returningColumns, ", ")
    }

    return query, values
}

func BuildSelectQuery(table string, columns []string, conditions map[string][]any) (string, []any) {
    values := []any{}
    where := []string{}

    query := fmt.Sprintf(
        "SELECT %s FROM %s",
        strings.Join(columns, ", "),
        table,
    )

    for column, condition := range conditions {
        if len(condition) != 2 {
            return "", nil
        }

        operator, ok := condition[0].(string)
        if !ok {
            return "", nil
        }

        where = append(where, column+" "+operator+" ?")
        values = append(values, condition[1])
    }

    if len(where) > 0 {
        query += " WHERE " + strings.Join(where, " AND ")
    }

    return query, values
}

func BuildUpdateQuery(table string, data map[string]any, conditions map[string][]any) (string, []any) {
    setClauses := []string{}
    whereClauses := []string{}
    values := []any{}

    for column, value := range data {
        setClauses = append(setClauses, column+" = ?")
        values = append(values, value)
    }

    for column, condition := range conditions {
        operator := condition[0].(string)
        value := condition[1]

        whereClauses = append(whereClauses, column+" "+operator+" ?")
        values = append(values, value)
    }

    query := fmt.Sprintf(
        "UPDATE %s SET %s WHERE %s",
        table,
        strings.Join(setClauses, ", "),
        strings.Join(whereClauses, " AND "),
    )

    return query, values
}

func BuildDeleteQuery(table string, conditions map[string][]any) (string, []any) {
    values := []any{}
    where := []string{}

    if len(conditions) == 0 {
        return "", nil
    }

    for column, condition := range conditions {
        if len(condition) != 2 {
            return "", nil
        }

        operator, ok := condition[0].(string)
        if !ok {
            return "", nil
        }

        where = append(where, column+" "+operator+" ?")
        values = append(values, condition[1])
    }

    query := fmt.Sprintf(
        "DELETE FROM %s WHERE %s",
        table,
        strings.Join(where, " AND "),
    )

    return query, values
}

