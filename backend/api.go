package main

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

func RegisterUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind the posted JSON into a user struct.
	// TODO: Sanitize.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// TODO: These two should be their own validation functions.
	if len(user.Username) < 3 {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Username must be at least 3 characters."})
		return
	}

	if len(user.Password) < 8 {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": "Password must be at least 8 characters."})
		return
	}

	// Create a hash of the POSTed password for storage.
	hash, err := argon2id.CreateHash(user.Password, argon2id.DefaultParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	err = AddCredentials(db, user.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return
	}

	// Send the confirmation back to the sender.
	ctx.JSON(http.StatusOK, gin.H{"username": user.Username})
}

func LoginUser(ctx *gin.Context, db *sql.DB) {
	// Attempt to bind posted JSON into user struct.
	var user Credentials
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// Query the user's password hash.
	hash, err := GetPasswordHash(db, user.Username)
	if err != nil {
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

	// Finally, return the session back to the sender.
	ctx.JSON(http.StatusOK, GenerateSessionKey(db, user.Username))
}

func GetUser(ctx *gin.Context, db *sql.DB) {
	ctx.JSON(http.StatusOK, gin.H{"Placeholder": "GetUser"})
}

func GetJournals(ctx *gin.Context, db *sql.DB) {
	ctx.JSON(http.StatusOK, gin.H{"Placeholder": "GetJournals"})
}
