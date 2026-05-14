// /home/krylon/go/src/github.com/blicero/vorleser/database/file.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 11:45:59 krylon>

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/blicero/vorleser/database/query"
	"github.com/blicero/vorleser/model"
)

const comma = ","

// FileAdd adds a File to the Database.
func (db *Database) FileAdd(f *model.File) error {
	const qid query.ID = query.FileAdd
	var (
		err  error
		stmt *sql.Stmt
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Failed to prepare query %s: %s\n",
			qid,
			err.Error())
		panic(err)
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows
EXEC_QUERY:
	if rows, err = stmt.Query(f.FolderID, valOrNil(f.BookID), f.Path); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		} else {
			err = fmt.Errorf("cannot add File %s: %w",
				f.Path,
				err)
			db.log.Printf("[ERROR] %s\n", err.Error())
			return err
		}
	} else {
		var id int64

		defer rows.Close() // nolint: errcheck

		if !rows.Next() {
			// CANTHAPPEN
			db.log.Printf("[ERROR] Query %s did not return a value\n",
				qid)
			return fmt.Errorf("query %s did not return a value", qid)
		} else if err = rows.Scan(&id); err != nil {
			var ex = fmt.Errorf("failed to get ID for newly added Folder %s: %w",
				f.Path,
				err)
			db.log.Printf("[ERROR] %s\n", ex.Error())
			return ex
		}

		f.ID = id
		return nil
	}
} // func (db *Database) FileAdd(f *model.File) error

// FileGetByFolder loads all Files that belong to the given Folder.
func (db *Database) FileGetByFolder(folder *model.Folder) ([]*model.File, error) {
	const qid query.ID = query.FileGetByFolder
	var err error
	var msg string
	var stmt *sql.Stmt
	var files []*model.File

GET_QUERY:
	if stmt, err = db.getQuery(qid); err != nil {
		if worthARetry(err) {
			time.Sleep(retryDelay)
			goto GET_QUERY
		} else {
			db.log.Printf("[ERROR] Error getting query %s: %s",
				qid,
				err.Error())
			return nil, err
		}
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows

EXEC_QUERY:
	if rows, err = stmt.Query(folder.ID); err != nil {
		if worthARetry(err) {
			time.Sleep(retryDelay)
			goto EXEC_QUERY
		} else {
			msg = fmt.Sprintf("Error querying all Feeds: %s",
				err.Error())
			db.log.Println(msg)
			return nil, errors.New(msg)
		}
	}

	defer rows.Close() // nolint: errcheck

	files = make([]*model.File, 0, 8)

	for rows.Next() {
		var (
			lastPlayed int64
			ord, title string
			pieces     []string
			file       = &model.File{FolderID: folder.ID}
		)

		if err = rows.Scan(
			&file.ID,
			&file.BookID,
			&file.Path,
			&title,
			&ord,
			&lastPlayed,
		); err != nil {
			msg = fmt.Sprintf("error scanning row: %s", err.Error())
			db.log.Printf("[ERROR] %s\n", msg)
			return nil, errors.New(msg)
		}

		pieces = strings.Split(ord, comma)
		if len(pieces) > 0 {
			file.Ord = make([]int64, len(pieces))

			for idx, piece := range pieces {
				var number int64

				if number, err = strconv.ParseInt(piece, 10, 64); err != nil {
					msg = fmt.Sprintf("Cannot parse sorting position %q: %s",
						piece,
						err.Error())
					db.log.Printf("[ERROR] %s\n",
						msg)
					return nil, errors.New(msg)
				}

				file.Ord[idx] = number
			}
		}

		file.LastPlayed = time.Unix(lastPlayed, 0)

		files = append(files, file)
	}

	return files, nil
} // func (db *Database) FileGetByFolder(folder *model.Folder) ([]*model.File, error)

// FileUpdatePosition sets a File's playback position.
func (db *Database) FileUpdatePosition(file *model.File, pos int64) error {
	const qid query.ID = query.FileUpdatePosition
	var (
		err, ex error
		stmt    *sql.Stmt
		res     sql.Result
		cnt     int64
		now     = time.Now()
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Failed to prepare query %s: %s\n",
			qid,
			err.Error())
		panic(err)
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

EXEC_QUERY:
	if res, err = stmt.Exec(pos, now.Unix(), file.ID); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		} else {
			ex = fmt.Errorf("cannot update position of File %s (%d): %w",
				file.Path,
				file.ID,
				err)
			db.log.Printf("[ERROR] %s\n", ex.Error())
			return ex
		}
	} else if cnt, err = res.RowsAffected(); err != nil {
		ex = fmt.Errorf("failed to get number of affected rows: %w",
			err)
		db.log.Printf("[ERROR] %s\n", ex.Error())
		return ex
	} else if cnt != 1 {
		ex = fmt.Errorf("unexpected number of affected rows for %s: %d (expected 1)",
			qid,
			cnt)
		db.log.Printf("[CRITICAL] %s\n", ex.Error())
		return ex
	}

	file.Position = pos
	file.LastPlayed = now
	return nil
} // func (db *Database) FileUpdatePosition(file *model.File, pos int64) error
