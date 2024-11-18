# SQL table structure outline

For now, we only need three tables - one storing users, one storing user entries, and one for storing session data.

We'll use the following SQL schema for now.
```sql
CREATE TABLE Users (
  Username TEXT PRIMARY KEY,
  JoinUnix INTEGER,
  EntryCount INTEGER
);

CREATE TABLE Entries (
  EntryID INTEGER PRIMARY KEY AUTOINCREMENT,
  EntryUnix INTEGER,
  Username TEXT,
  Content TEXT
);

CREATE TABLE Sessions (
  Username TEXT,
  SessionKey TEXT,
  SessionUnix INTEGER
);
```