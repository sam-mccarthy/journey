# SQL table structure outline

For now, we need four tables - one storing users, one storing credentials, one storing user entries, and one for storing session data.

We'll use the following SQL schema.
```sql
CREATE TABLE IF NOT EXISTS Users (
  Username TEXT PRIMARY KEY,
  JoinUnix INTEGER NOT NULL,
  EntryCount INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Credentials (
  Username TEXT PRIMARY KEY,
  Hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Entries (
  EntryID INTEGER PRIMARY KEY AUTOINCREMENT,
  EntryUnix INTEGER NOT NULL,
  Username TEXT NOT NULL,
  Content TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS Sessions (
  Username TEXT NOT NULL,
  SessionKey TEXT NOT NULL,
  SessionUnix INTEGER NOT NULL
);
```