package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func closeDatabase(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setupDatabase(dataSource string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dataSource)
	if err != nil {
		return nil, err
	}

	err = initializeDatabase(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func route(db *sql.DB) {
	router := gin.Default()
	router.POST("/api/register", func(ctx *gin.Context) { registerUser(ctx, db) })
	router.POST("/api/login", func(ctx *gin.Context) { loginUser(ctx, db) })

	router.POST("/api/user", func(ctx *gin.Context) { getUser(ctx, db) })
	router.POST("/api/journals", func(ctx *gin.Context) { getJournals(ctx, db) })

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Open a SQLite database. For now, we'll keep it in memory
	db, err := setupDatabase(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer closeDatabase(db)

	route(db)
}
