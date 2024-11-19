# SQL table structure outline

For now, we only need three tables - one storing users, one storing user entries, and one for storing session data.

We'll use the following SQL schema for now.
```sql
CREATE TABLE IF NOT EXISTS Users (
  Username TEXT PRIMARY KEY,
  JoinUnix INTEGER NOT NULL,
  EntryCount INTEGER NOT NULL
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