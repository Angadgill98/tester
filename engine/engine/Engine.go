package Engine

import "database/sql"

type Engine struct {
	DB *sql.DB
}

func NewEngine(db *sql.DB) *Engine {
	return &Engine{
		DB: db,
	}
}


