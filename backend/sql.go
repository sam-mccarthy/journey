package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

// Credentials - Store user credentials during login and register.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// User - Store basic user information. Later, more statistics might be nice.
type User struct {
	Username string    `json:"username"`
	JoinDate time.Time `json:"joinDate"`
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

func initializeDatabase(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS Users (
		  Username TEXT PRIMARY KEY,
		  JoinUnix INTEGER NOT NULL
		)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Entries (
		  EntryID INTEGER PRIMARY KEY AUTOINCREMENT,
		  EntryUnix INTEGER NOT NULL,
		  Username TEXT NOT NULL,
		  Content TEXT NOT NULL
		)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Credentials (
		  Username TEXT PRIMARY KEY,
		  Hash TEXT NOT NULL
		)`)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS Sessions (
		  Username TEXT NOT NULL,
		  SessionKey TEXT NOT NULL,
		  SessionUnix INTEGER NOT NULL
		)`)

	return err
}

func getUser(db *sql.DB, username string) (User, error) {
	var user User

	row := db.QueryRow(`SELECT username, joinUnix FROM users WHERE Username = ?`, username)
	err := row.Scan(&user.Username, &user.JoinDate)

	if err != nil {
		return User{}, errors.New("failed to get user")
	}

	return user, nil
}

func getEntries(db *sql.DB, username string, limit int) ([]Entry, error) {
	entries := make([]Entry, 0)

	rows, err := db.Query(`SELECT username, date, content FROM Entries WHERE Username = ? LIMIT ?`, username, limit)
	if err != nil {
		return nil, errors.new("failed to get entry data")
	}

	for rows.Next() {
		entry := Entry{}
		if err := rows.Scan(&entry.Username, &entry.Date, &entry.Content); err != nil {
			return nil, errors.New("failed to scan entry")
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func addUser(db *sql.DB, username string) error {
	_, err := db.Exec(`INSERT INTO users (username, joinDate) VALUES (?, ?)`, username, time.Now())
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func addEntry(db *sql.DB, username string, content string) error {
	_, err := db.Exec(`INSERT INTO entries (username, date, content) VALUES (?, ?, ?)`, username, time.Now(), content)
	if err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func addCredentials(db *sql.DB, username string, hash string) error {
	var exists bool
	row := db.QueryRow("SELECT EXISTS(SELECT 1 FROM credentials WHERE username = ?)", username)
	err := row.Scan(&exists)

	if err == nil {
		if exists {
			return errors.New("user already exists")
		}

		// Insert the new user data into the table
		_, err = db.Exec("INSERT INTO credentials (username, hash) VALUES (?, ?)", username, hash)
		if err == nil {
			return nil
		}
	}

	log.Println(err.Error())
	return errors.New("failed to access data")
}

func getPasswordHash(db *sql.DB, username string) (string, error) {
	row := db.QueryRow("SELECT username, hash FROM credentials WHERE username = ?", username)

	var sqlUsername string
	var hash string

	err := row.Scan(&sqlUsername, &hash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("invalid username or password")
		}
		return "", errors.New("failed processing password")
	}

	return hash, nil
}

func newSession(db *sql.DB, username string) (Session, error) {
	session := Session{
		Username:    username,
		SessionUnix: time.Now(),
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return Session{}, errors.New("problem generating session key")
	}

	session.SessionKey = hex.EncodeToString(bytes)

	_, err := db.Exec("INSERT INTO sessions (username, sessionKey, sessionUnix) VALUES (?, ?, ?)",
		session.Username, session.SessionKey, session.SessionUnix)
	if err != nil {
		return Session{}, errors.New("failed creating session")
	}

	return session, nil
}
