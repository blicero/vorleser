// /home/krylon/go/src/github.com/blicero/vorleser/database/query/query.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 10:40:19 krylon>

// Package query defines symbolic constants to identify the queries we
// perform on the Database.
package query

//go:generate stringer -type=ID

// ID identifies a Database query.
type ID uint8

const (
	FolderAdd ID = iota
	FolderGetByID
	FolderGetByPath
	FolderGetAll
	FolderUpdateLastScan
	FolderDelete
	FileAdd
	FileGetByFolder
	FileGetByPath
	FileGetByPathLike
	FileGetByBook
	FileUpdatePosition
	FileUpdateBook
	FileUpdateOrd
	FileDelete
)
