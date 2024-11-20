package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
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

func InitializeDatabase(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Users (
		  Username TEXT PRIMARY KEY,
		  Hash TEXT NOT NULL,
		  JoinUnix INTEGER NOT NULL,
		  EntryCount INTEGER NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Entries (
		  EntryID INTEGER PRIMARY KEY AUTOINCREMENT,
		  EntryUnix INTEGER NOT NULL,
		  Username TEXT NOT NULL,
		  Content TEXT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Credentials (
		  Username TEXT PRIMARY KEY,
		  Hash TEXT NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Sessions (
		  Username TEXT NOT NULL,
		  SessionKey TEXT NOT NULL,
		  SessionUnix INTEGER NOT NULL
		)`)

	if err != nil {
		log.Fatal(err)
	}
}

func AddUser(db *sql.DB, user User) {

}

func AddCredentials(db *sql.DB, username string, hash string) error {
	// Insert the new user data into the table.
	_, err := db.Exec("INSERT INTO credentials (username, hash) VALUES (?, ?)", username, hash)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
		return err
	}

	return nil
}

func GetPasswordHash(db *sql.DB, username string) (string, error) {
	row := db.QueryRow("SELECT username, hash FROM credentials WHERE username = ?", username)

	var sqlUsername string
	var hash string

	err := row.Scan(&sqlUsername, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	return hash, nil
}

func GenerateSessionKey(db *sql.DB, username string) Session {
	session := Session{
		Username:    username,
		SessionUnix: time.Now(),
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Problem generating session key")
		return Session{}
	}

	session.SessionKey = hex.EncodeToString(bytes)

	return session
}
