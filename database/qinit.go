// /home/krylon/go/src/github.com/blicero/vorleser/database/qinit.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 10:28:02 krylon>

package database

var qinit = []string{
	`
CREATE TABLE book (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    author TEXT,
    url TEXT,
    current_file INTEGER NOT NULL DEFAULT 0
) STRICT
`,
	`
CREATE TABLE folder (
    id        INTEGER PRIMARY KEY,
    path      TEXT UNIQUE NOT NULL,
    last_scan INTEGER NOT NULL DEFAULT 0
) STRICT
`,
	"CREATE INDEX dir_scan_idx ON folder (last_scan)",
	`
CREATE TABLE file (
    id		INTEGER PRIMARY KEY,
    folder_id	INTEGER NOT NULL,
    book_id	INTEGER,
    path	TEXT UNIQUE NOT NULL,
    title	TEXT,
    ord		TEXT,
    position	INTEGER NOT NULL DEFAULT 0,
    last_played INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (folder_id) REFERENCES folder (id)
        ON UPDATE RESTRICT
        ON DELETE CASCADE,
    FOREIGN KEY (book_id) REFERENCES book (id)
        ON UPDATE RESTRICT
        ON DELETE CASCADE
) STRICT
`,
}
