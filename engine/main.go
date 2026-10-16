
package main

import "log"

func main() {
	db, err := SetSQLite()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer db.Close()
}
