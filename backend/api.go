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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: These two should be their own validation functions.
	if len(user.Username) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "username must be at least 3 characters."})
		return
	}

	if len(user.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters."})
		return
	}

	// Create a hash of the POSTed password for storage.
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query the user's password hash.
	hash, err := getPasswordHash(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check the password against the stored hash.
	match, _, err := argon2id.CheckHash(user.Password, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !match {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	session, err := generateSessionKey(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
