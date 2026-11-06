package sequences

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)


type Sequence struct {
	SequenceID   int64   
	AppID        int64   
	SequenceName string  
	ActionIDs    []int64 
}


func CreateSequence(db *sql.DB, sequence_name string, app_id int64) (Sequence, error) {
    sequence := Sequence{
        AppID:        app_id,
        SequenceName: sequence_name,
        ActionIDs:    []int64{},
    }

    actionIDs, err := json.Marshal(sequence.ActionIDs)
    if err != nil {
        return Sequence{}, err
    }

    query, values := BuildInsertQuery("sequences", map[string]any{
        "app_id":        sequence.AppID,
        "sequence_name": sequence.SequenceName,
        "action_ids":    string(actionIDs),
    }, []string{"sequence_id"})

    err = db.QueryRow(query, values...).Scan(&sequence.SequenceID)
    if err != nil {
        return Sequence{}, err
    }

    return sequence, nil
}

func CreateTemporarySequence(sequence_name string, app_id int64) Sequence {
	return Sequence{
		SequenceID:   0,
		AppID:        999999,
		SequenceName: sequence_name,
		ActionIDs:    []int64{},
	}
}

func GetSequence(db *sql.DB, sequence_id int64) (Sequence, error) {
    query, values := BuildSelectQuery("sequences", []string{
        "sequence_id", "app_id", "sequence_name", "action_ids",
    }, map[string][]any{
        "sequence_id": {"=", sequence_id},
    })

    var sequence Sequence
    var actionIDs string

    err := db.QueryRow(query, values...).Scan(
        &sequence.SequenceID,
        &sequence.AppID,
        &sequence.SequenceName,
        &actionIDs,
    )
    if err != nil {
        return Sequence{}, err
    }

    if err := json.Unmarshal([]byte(actionIDs), &sequence.ActionIDs); err != nil {
        return Sequence{}, err
    }

    return sequence, nil
}

func UpdateSequence(db *sql.DB, sequence_id int64, sequence_name string, action_ids []int64) error {
    actionIDs, err := json.Marshal(action_ids)
    if err != nil {
        return err
    }

    query, values := BuildUpdateQuery("sequences", map[string]any{
        "sequence_name": sequence_name,
        "action_ids":    string(actionIDs),
    }, map[string][]any{
        "sequence_id": {"=", sequence_id},
    })

    _, err = db.Exec(query, values...)
    return err
}

func DeleteSequence(db *sql.DB, sequence_id int64) error {
    query, values := BuildDeleteQuery("sequences", map[string][]any{
        "sequence_id": {"=", sequence_id},
    })

    _, err := db.Exec(query, values...)
    return err
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

