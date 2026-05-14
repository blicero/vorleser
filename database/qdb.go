// /home/krylon/go/src/github.com/blicero/vorleser/database/qdb.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 10:42:13 krylon>

package database

import "github.com/blicero/vorleser/database/query"

var qdb = map[query.ID]string{
	query.FolderAdd:            "INSERT INTO folder (path) VALUES (?)",
	query.FolderGetByID:        "SELECT path, last_scan FROM folder WHERE id = ?",
	query.FolderGetByPath:      "SELECT id, last_scan FROM folder where path = ?",
	query.FolderGetAll:         "SELECT id, path, last_scan FROM folder",
	query.FolderUpdateLastScan: "UPDATE folder SET last_scan = ? WHERE id = ?",
	query.FolderDelete:         "DELETE FROM folder WHERE id = ?",
	query.FileAdd:              "INSERT INTO file (folder_id, book_id, path) VALUES (?, ?, ?)",
	query.FileGetByFolder: `
SELECT
    id,
    book_id,
    path,
    COALESCE(title, ''),
    COALESCE(ord, ''),
    position,
    last_played
FROM file
WHERE folder_id = ?
`,
	query.FileGetByBook: `
SELECT
    id,
    folder_id,
    path,
    COALESCE(title, ''),
    COALESCE(ord, ''),
    position,
    last_played
FROM file
WHERE book_id = ?
`,
	query.FileGetByPath: `
SELECT
    id,
    folder_id,
    book_id,
    COALESCE(title, ''),
    COALESCE(ord, ''),
    position,
    last_played
FROM file
WHERE path = ?
`,
	query.FileGetByPathLike: `
SELECT
    id,
    folder_id,
    book_id,
    COALESCE(title, ''),
    COALESCE(ord, ''),
    position,
    last_played
FROM file
WHERE path LIKE ?
`,
	query.FileUpdatePosition: "UPDATE file SET position = ?, last_played = ? WHERE id = ?",
	query.FileUpdateBook:     "UPDATE file SET book_id = ? WHERE id = ?",
	query.FileUpdateOrd:      "UPDATE file SET ord = ? WHERE id = ?",
	query.FileDelete:         "DELETE FROM file WHERE id = ?",
}
