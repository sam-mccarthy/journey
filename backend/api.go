package backend

import (
	"database/sql"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

// Credentials - Store user credentials during login and register.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

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

	err = addCredentialData(db, user.Username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = addUserData(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	hash, err := getPasswordHashData(db, user.Username)
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

	session, err := newSessionData(db, user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Finally, return the session back to the sender.
	ctx.JSON(http.StatusOK, session)
}

func getUser(ctx *gin.Context, db *sql.DB) {
	data := struct {
		Username string `json:"username"`
		Session  string `json:"session"`
	}{Username: ""}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	if err := checkSessionData(db, data.Username, data.Session); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := getUserData(db, data.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func getJournals(ctx *gin.Context, db *sql.DB) {
	data := struct {
		Username string `json:"username"`
		Limit    int    `json:"limit"`
		Offset   int    `json:"offset"`
		Session  string `json:"session"`
	}{}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	if err := checkSessionData(db, data.Username, data.Session); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	entries, err := getEntryData(db, data.Username, data.Limit, data.Offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, entries)
}

func createEntry(ctx *gin.Context, db *sql.DB) {
	data := struct {
		Username string `json:"username"`
		Content  string `json:"content"`
		Session  string `json:"session"`
	}{}

	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request data"})
		return
	}

	if err := checkSessionData(db, data.Username, data.Session); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	err := addEntryData(db, data.Username, data.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Add statuses to every request, or get rid of this
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func checkUsername(username string) bool {
	return len(username) >= 3
}

func checkPassword(password string) bool {
	return len(password) >= 8
}
