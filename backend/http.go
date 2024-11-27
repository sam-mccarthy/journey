package backend

import (
	"database/sql"
	"github.com/gin-contrib/graceful"
	"github.com/gin-gonic/gin"
	"log"
)

func SignalRoute(db *sql.DB, channel chan struct{}) {
	router, err := graceful.Default()
	if err != nil {
		log.Fatal(err)
	}

	router.POST("/api/register", func(ctx *gin.Context) { registerUser(ctx, db) })
	router.POST("/api/login", func(ctx *gin.Context) { loginUser(ctx, db) })

	router.POST("/api/user", func(ctx *gin.Context) { getUser(ctx, db) })
	router.POST("/api/journals", func(ctx *gin.Context) { getJournals(ctx, db) })

	router.POST("/api/new_entry", func(ctx *gin.Context) { createEntry(ctx, db) })

	go func() {
		err := router.Run()

		if err != nil {
			log.Fatal(err)
		}
	}()

	defer router.Close()

	<-channel
}

func Route(db *sql.DB) {
	SignalRoute(db, make(chan struct{}))
}
