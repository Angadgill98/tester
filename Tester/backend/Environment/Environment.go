package environment

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
)




type Environment struct {
	EnvironmentId int 
	EnvironmentName string
	KeysValues map[string]any
}


func CreateEnvironment(environment_name string) Environment {
    return Environment{
        EnvironmentName: environment_name,
        KeysValues:      map[string]any{},
    }
}

func GetEnvironment(db *sql.DB, environmentID int) (Environment, error) {
    query, values := BuildSelectQuery("environments", []string{
        "environment_id",
        "environment_name",
        "keys_values",
    }, map[string][]any{
        "environment_id": {"=", environmentID},
    })

    var env Environment
    var keysValues string

    err := db.QueryRow(query, values...).Scan(&env.EnvironmentId, &env.EnvironmentName, &keysValues)
    if err != nil {
        return Environment{}, fmt.Errorf("failed to get environment: %w", err)
    }

    if err := json.Unmarshal([]byte(keysValues), &env.KeysValues); err != nil {
        return Environment{}, fmt.Errorf("failed to decode environment keys: %w", err)
    }

    return env, nil
}

func UpdateEnvironment(db *sql.DB, environmentID int, env_name string, keysValues json.RawMessage) error {
    if !json.Valid(keysValues) {
        return fmt.Errorf("invalid environment keys JSON")
    }

    query, values := BuildUpdateQuery(
        "environments",
        map[string]any{
            "environment_name": env_name,
            "keys_values":     string(keysValues),
        },
        map[string][]any{
            "environment_id": {"=", environmentID},
        },
    )

    result, err := db.Exec(query, values...)
    if err != nil {
        return fmt.Errorf("failed to update environment: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return err
    }

    if rowsAffected == 0 {
        return fmt.Errorf("environment not found")
    }

    return nil
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