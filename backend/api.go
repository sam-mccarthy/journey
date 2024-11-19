package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

// Credentials - Store user credentials during login and register.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User - Store basic user information. Later, more statistics might be nice.
type User struct {
	Username   string    `json:"username"`
	JoinDate   time.Time `json:"joinDate"`
	EntryCount int       `json:"entryCount"`
}

// Entry - Store a user's journal entry.
type Entry struct {
	Username int       `json:"username"`
	Date     time.Time `json:"date"`
	Content  string    `json:"content"`
}

// Session - Stores a user's session key
type Session struct {
	Username    string    `json:"username"`
	SessionKey  string    `json:"sessionKey"`
	SessionUnix time.Time `json:"sessionUnix"`
}

func GenerateSessionKey(username string) Session {
	return Session{}
}

func RegisterUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind the posted JSON into a user struct.
	// TODO: Sanitize.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Create a hash of the POSTed password for storage.
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Insert the new user data into the table.
	_, err = db.Exec("INSERT INTO credentials (username, hash) VALUES (?, ?)", user.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Send the confirmation back to the sender.
	ctx.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func LoginUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind posted JSON into user struct.
	// TODO: Sanitize.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Query the user's username and hash from the table.
	row := db.QueryRow("SELECT username, hash FROM credentials WHERE username = ?", user.Username)

	// Retrieve the username and hash.
	var username string
	var hash string
	err := row.Scan(&username, &hash)

	if err != nil {
		// If there is an error, and it's a lack of rows, the username is invalid.
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid username or password"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Check the password against the stored hash.
	match, _, err := argon2id.CheckHash(user.Password, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	if !match {
		ctx.JSON(http.StatusUnauthorized, gin.H{"Error": "Invalid username or password"})
		return
	}

	sessionKey := GenerateSessionKey(user.Username)

	// Finally, return the confirmation back to the sender.
	ctx.JSON(http.StatusOK, gin.H{"username": username, "sessionKey": sessionKey})
}

func GetUser(ctx *gin.Context, db *sql.DB) {
	ctx.JSON(http.StatusOK, gin.H{"guh": "guh"})
}

func GetJournals(ctx *gin.Context, db *sql.DB) {

}
