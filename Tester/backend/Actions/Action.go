package actions

import (
	"database/sql"
	"fmt"
	"strings"
)



type Action struct {
    ActionType string
    ActionId   int
}

type ActionObj struct {
	Db *sql.DB
}

func NewActionObj(db *sql.DB) *ActionObj {
	return &ActionObj{
		Db: db,
	}
}

func (a *ActionObj) CreateAction(db *sql.DB, action_type string) (int64, error) {
    query, values := BuildInsertQuery("actions", map[string]any{
        "action_type": action_type,
    }, []string{})

    result, err := db.Exec(query, values...)
    if err != nil {
        return 0, err
    }

    actionID, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return actionID, nil
}

func (a *ActionObj) GetActions(db *sql.DB, actionIDs []int64) ([]Action, error) {
    actionsList := make([]Action, 0, len(actionIDs))

    for _, actionID := range actionIDs {
        query, values := BuildSelectQuery("actions", []string{"action_id", "action_type"}, map[string][]any{
            "action_id": {"=", actionID},
        })

        var action Action
        err := db.QueryRow(query, values...).Scan(&action.ActionId, &action.ActionType)
        if err != nil {
            return nil, fmt.Errorf("failed to get action ID %d: %w", actionID, err)
        }

        actionsList = append(actionsList, action)
    }

    return actionsList, nil
}


func (a *ActionObj) GetActionType(db *sql.DB, action_id int) (string, error) {
    query, values := BuildSelectQuery(
        "actions",
        []string{"action_type"},
        map[string][]any{
            "action_id": {"=", action_id},
        },
    )

    var action_type string

    err := db.QueryRow(query, values...).Scan(&action_type)
    if err != nil {
        return "", err
    }

    return action_type, nil
}

func (a *ActionObj) DeleteAction(db *sql.DB, actionID int64) error {
    query, values := BuildDeleteQuery("actions", map[string][]any{
        "action_id": {"=", actionID},
    })

    _, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to delete action: %w", err)
    }

    return nil
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

