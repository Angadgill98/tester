package main

import (
	actions "Tester/backend/Actions"
	app "Tester/backend/App"
	sequences "Tester/backend/Sequences"
	"context"
	"database/sql"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
	Db *sql.DB
	App *app.ApplicationObj
	Action *actions.ActionObj
	Sequence *sequences.SequenceObj
}

// NewApp creates a new App application struct
func NewApp(db *sql.DB) *App {
	return &App{
		Db: db,
		App: app.CreateApplicationObject(db),
		Action: actions.NewActionObj(db),
		Sequence: sequences.NewSequenceObj(db),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
