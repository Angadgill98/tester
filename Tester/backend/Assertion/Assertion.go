package assertion

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)



type Assertion struct {
	AssertionId int64
	ActionId int64
	StatusCode int64
	Assertion []json.RawMessage
}


func CreateAssertion(db *sql.DB, actionId int64, statusCode int64, assertions []json.RawMessage) (int64, error) {
	assertionJSON, err := json.Marshal(assertions)
	if err != nil {
		return 0, err
	}

	query, values := BuildInsertQuery("assertions", map[string]any{
		"action_id":   actionId,
		"status_code": statusCode,
		"assertion":   string(assertionJSON),
	}, []string{"assertion_id"})

	var assertionID int64

	err = db.QueryRow(query, values...).Scan(&assertionID)
	if err != nil {
		return 0, err
	}

	return assertionID, nil
}


func InsertActionAssertion(db *sql.DB, assertionId int64, assertion json.RawMessage) error {
	query := `
		UPDATE assertions
		SET assertion = json_insert(assertion, '$[#]', json(?))
		WHERE assertion_id = ?
	`

	_, err := db.Exec(query, string(assertion), assertionId)
	if err != nil {
		return err
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

