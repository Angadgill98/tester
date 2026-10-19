package main

import (
	tester "engine/Tester"
	"log"
)

func main() {
	db, err := SetSQLite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// _, err = db.Exec(`DELETE FROM apps`)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	_, err = db.Exec(`INSERT OR IGNORE INTO apps (id, name, parent_app_id, action_ids, sequence_ids, children) VALUES (0, 'GLOBAL', NULL, '[]', '[]', '[]')`)
	if err != nil {
		log.Fatal(err)
	}
	tester:=tester.CreateTesterObj()
	CreateSocket(db,&tester).StartSocket()
}