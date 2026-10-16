package main

import (
	application "engine/Application"
	"log"
)

func main() {
	db, err := SetSQLite()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	exists, err := application.IsChildApplication(db,0, "GLOBAL")
	if err != nil {
		log.Fatal(err)
	}

	if !exists {
		_, err := application.InsertApplication(db, "GLOBAL", 0)
		if err != nil {
			log.Fatal(err)
		}
	}

	CreateSocket(db).StartSocket()
	
}