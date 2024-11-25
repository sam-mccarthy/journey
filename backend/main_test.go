package main

import (
	"log"
)

func setup() {
	db, err := setupDatabase(":memory:")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer closeDatabase(db)

	route(db)
}
