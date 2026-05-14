// /home/krylon/go/src/github.com/blicero/vorleser/database/folder.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 11:08:46 krylon>

package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/blicero/vorleser/database/query"
	"github.com/blicero/vorleser/model"
)

// FolderAdd adds a folder to the Database.
func (db *Database) FolderAdd(folder *model.Folder) error {
	const qid query.ID = query.FolderAdd
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
	if rows, err = stmt.Query(); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		} else {
			err = fmt.Errorf("cannot add Folder %s: %w",
				folder.Path,
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
				folder.Path,
				err)
			db.log.Printf("[ERROR] %s\n", ex.Error())
			return ex
		}

		folder.ID = id
		return nil
	}
} // func (db *Database) FolderAdd(folder *model.Folder) error

// FolderGetByID looks up a Folder by its ID.
func (db *Database) FolderGetByID(id int64) (*model.Folder, error) {
	const qid query.ID = query.FolderGetByID
	var (
		err  error
		stmt *sql.Stmt
	)

	if stmt, err = db.getQuery(qid); err != nil {
		db.log.Printf("[ERROR] Cannot prepare query %s: %s\n",
			qid,
			err.Error())
		return nil, err
	} else if db.tx != nil {
		stmt = db.tx.Stmt(stmt)
	}

	var rows *sql.Rows

EXEC_QUERY:
	if rows, err = stmt.Query(id); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		}

		return nil, err
	}

	defer rows.Close() // nolint: errcheck,gosec

	if rows.Next() {
		var (
			f        = &model.Folder{ID: id}
			lastScan int64
		)

		if err = rows.Scan(&f.Path, &lastScan); err != nil {
			var ex = fmt.Errorf("failed to scan row: %w", err)
			db.log.Printf("[ERROR] %s\n", ex.Error())
			return nil, ex
		}

		f.LastScan = time.Unix(lastScan, 0)
		return f, nil
	}

	return nil, nil
} // func (db *Database) FolderGetByID(id int64) (*model.Folder, error)

// FolderGetAll loads all Folders from the Database.
func (db *Database) FolderGetAll() ([]*model.Folder, error) {
	const qid query.ID = query.FolderGetAll
	var err error
	var msg string
	var stmt *sql.Stmt
	var folders []*model.Folder

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
	if rows, err = stmt.Query(); err != nil {
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

	folders = make([]*model.Folder, 0, 8)

	for rows.Next() {
		var (
			lastScan int64
			folder   = new(model.Folder)
		)

		if err = rows.Scan(
			&folder.ID,
			&folder.Path,
			&lastScan,
		); err != nil {
			msg = fmt.Sprintf("error scanning row: %s", err.Error())
			db.log.Printf("[ERROR] %s\n", msg)
			return nil, errors.New(msg)
		}

		folder.LastScan = time.Unix(lastScan, 0)

		folders = append(folders, folder)
	}

	return folders, nil
} // func (db *Database) FolderGetAll() ([]*model.Folder, error)

// FolderUpdateLastScan updates a Folder's timestamp.
func (db *Database) FolderUpdateLastScan(folder *model.Folder) error {
	const qid query.ID = query.FolderUpdateLastScan
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
	if res, err = stmt.Exec(now.Unix(), folder.ID); err != nil {
		if worthARetry(err) {
			waitForRetry()
			goto EXEC_QUERY
		} else {
			ex = fmt.Errorf("cannot update last_scan of Folder %s (%d): %w",
				folder.Path,
				folder.ID,
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

	folder.LastScan = now
	return nil
} // func (db *Database) FolderUpdateLastScan(folder *model.Folder) error
