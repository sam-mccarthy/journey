package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/glebarez/go-sqlite"
)

func main() {
	db, sql_err := sql.Open("sqlite", ":memory:")
	if sql_err != nil {
		log.Fatal(sql_err)
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	router := gin.Default()
	router.POST("/api/register", func(ctx *gin.Context) { RegisterUser(ctx, db) })
	router.POST("/api/login", func(ctx *gin.Context) { LoginUser(ctx, db) })

	router.GET("/api/user", func(ctx *gin.Context) { GetUser(ctx, db) })
	router.GET("/api/journals", func(ctx *gin.Context) { GetJournals(ctx, db) })
}
