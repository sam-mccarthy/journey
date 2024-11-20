package main

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

func registerUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind the posted JSON into a user struct.
	// TODO: Sanitize.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	if !checkUsername(user.Username) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad username"})
		return
	}

	if !checkPassword(user.Password) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad password"})
		return
	}

	// Create a hash of the POSTed password for storage.
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed processing password"})
		return
	}

	err = addCredentials(db, user.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the confirmation back to the sender.
	ctx.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func loginUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind posted JSON into user struct.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	// Query the user's password hash.
	hash, err := getPasswordHash(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed processing password"})
		return
	}

	// Check the password against the stored hash.
	match, _, err := argon2id.CheckHash(user.Password, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed checking password"})
		return
	}

	if !match {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	session, err := generateSessionKey(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed generating session"})
		return
	}

	// Finally, return the session back to the sender.
	ctx.JSON(http.StatusOK, session)
}

func getUser(ctx *gin.Context, db *sql.DB) {
	ctx.JSON(http.StatusOK, gin.H{"placeholder": "getUser"})
}

func getJournals(ctx *gin.Context, db *sql.DB) {
	ctx.JSON(http.StatusOK, gin.H{"placeholder": "getJournals"})
}

func checkUsername(username string) bool {
	return len(username) >= 3
}

func checkPassword(password string) bool {
	return len(password) >= 8
}
