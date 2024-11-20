package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	// Open a SQLite database. For now, we'll keep it in memory
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
		return
	}

	InitializeDatabase(db)

	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// We need to be able to register, login, view the user, and view the user's journals.
	// Later, we'll also need a POST request available for creating journal posts, among potential others.
	router := gin.Default()
	router.POST("/api/register", func(ctx *gin.Context) { RegisterUser(ctx, db) })
	router.POST("/api/login", func(ctx *gin.Context) { LoginUser(ctx, db) })

	router.GET("/api/user", func(ctx *gin.Context) { GetUser(ctx, db) })
	router.GET("/api/journals", func(ctx *gin.Context) { GetJournals(ctx, db) })

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}
