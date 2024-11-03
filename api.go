package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

type Credentials struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	JoinDate   time.Time `json:"joinDate"`
	EntryCount int       `json:"entryCount"`
}

type Entry struct {
	UserID  int       `json:"username"`
	Date    time.Time `json:"date"`
	Content string    `json:"content"`
}

func RegisterUser(ctx *gin.Context, db *sql.DB) {
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	result, err := db.Exec("INSERT INTO users (username, hash) VALUES (?, ?)", user.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func LoginUser(ctx *gin.Context, db *sql.DB) {
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	row := db.QueryRow("SELECT id, hash FROM users WHERE username = ?", user.Username)

	var id int
	var hash string
	err := row.Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid username or password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	match, _, err := argon2id.CheckHash(user.Password, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !match {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid username or password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

func GetUser(ctx *gin.Context, db *sql.DB) {

}

func GetJournals(ctx *gin.Context, db *sql.DB) {

}
