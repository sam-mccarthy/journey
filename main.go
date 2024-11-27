package main

import (
	"journey/backend"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	// Open a SQLite database. For now, we'll keep it in memory
	db, err := backend.SetupDatabase(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer backend.CloseDatabase(db)

	backend.Route(db)
}
