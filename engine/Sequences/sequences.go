package sequences

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Sequence struct {
	SequenceID   int64   `json:"sequence_id"`
	AppID        int64   `json:"app_id"`
	SequenceName string  `json:"sequence_name"`
	ActionIDs    []int64 `json:"action_ids"`
}

// func CreateSequence(sequenceID int64, appID int64, sequenceName string, actionIDs []int64) Sequence {
// 	return Sequence{
// 		SequenceID:   sequenceID,
// 		AppID:        appID,
// 		SequenceName: sequenceName,
// 		ActionIDs:    actionIDs,
// 	}
// }



func InsertSequence(db *sql.DB, appID int64, sequenceName string) (int64, error) {
	result, err := db.Exec(
		`INSERT INTO sequences (app_id, sequence_name) VALUES (?, ?)`,
		appID,
		sequenceName,
	)
	if err != nil {
		return 0, fmt.Errorf("insert sequence: %w", err)
	}

	sequenceID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get sequence ID: %w", err)
	}

	return sequenceID, nil
}

func GetSequence(db *sql.DB, sequenceID int64) (Sequence, error) {
	var sequence Sequence

	err := db.QueryRow(
		`SELECT id, app_id, sequence_name FROM sequences WHERE id = ?`,
		sequenceID,
	).Scan(&sequence.SequenceID, &sequence.AppID, &sequence.SequenceName)
	if err != nil {
		if err == sql.ErrNoRows {
			return Sequence{}, fmt.Errorf("sequence %d not found", sequenceID)
		}
		return Sequence{}, fmt.Errorf("get sequence: %w", err)
	}

	return sequence, nil
}

func SequenceExists(db *sql.DB, appID int64, sequenceName string) (bool, error) {
	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (
			SELECT 1 FROM sequences
			WHERE app_id = ? AND sequence_name = ?
		)`,
		appID, sequenceName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check sequence existence: %w", err)
	}

	return exists, nil
}


func SaveSequence(db *sql.DB, sequence Sequence) error {
	actionIDsJSON, err := json.Marshal(sequence.ActionIDs)
	if err != nil {
		return fmt.Errorf("marshal sequence action IDs: %w", err)
	}

	result, err := db.Exec(
		`UPDATE sequences SET app_id = ?, sequence_name = ?, action_ids = ? WHERE id = ?`,
		sequence.AppID,
		sequence.SequenceName,
		string(actionIDsJSON),
		sequence.SequenceID,
	)
	if err != nil {
		return fmt.Errorf("save sequence: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check saved sequence: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("sequence %d not found", sequence.SequenceID)
	}

	return nil
}

func GetSequencesByAppID(db *sql.DB, appID int64) ([]Sequence, error) {
	rows, err := db.Query(
		`SELECT id, app_id, sequence_name, action_ids
		 FROM sequences
		 WHERE app_id = ?`,
		appID,
	)
	if err != nil {
		return nil, fmt.Errorf("get sequences: %w", err)
	}
	defer rows.Close()

	var result []Sequence

	for rows.Next() {
		var sequence Sequence
		var actionIDsJSON string

		if err := rows.Scan(&sequence.SequenceID, &sequence.AppID, &sequence.SequenceName, &actionIDsJSON); err != nil {
			return nil, fmt.Errorf("scan sequence: %w", err)
		}

		if err := json.Unmarshal([]byte(actionIDsJSON), &sequence.ActionIDs); err != nil {
			return nil, fmt.Errorf("parse sequence action IDs: %w", err)
		}

		result = append(result, sequence)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sequences: %w", err)
	}

	return result, nil
}

func GetSequenceByName(db *sql.DB, appID int64, sequenceName string) (Sequence, error) {
	var sequence Sequence
	var actionIDsJSON string

	err := db.QueryRow(
		`SELECT id, app_id, sequence_name, action_ids
		 FROM sequences
		 WHERE app_id = ? AND sequence_name = ?`,
		appID, sequenceName,
	).Scan(&sequence.SequenceID, &sequence.AppID, &sequence.SequenceName, &actionIDsJSON)

	if err != nil {
		if err == sql.ErrNoRows {
			return Sequence{}, fmt.Errorf("sequence %q not found in application %d", sequenceName, appID)
		}
		return Sequence{}, fmt.Errorf("get sequence by name: %w", err)
	}

	if err := json.Unmarshal([]byte(actionIDsJSON), &sequence.ActionIDs); err != nil {
		return Sequence{}, fmt.Errorf("parse sequence action IDs: %w", err)
	}

	return sequence, nil
}