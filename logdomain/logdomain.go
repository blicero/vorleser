// /home/krylon/go/src/github.com/blicero/vorleser/logdomain/logdomain.go
// -*- mode: go; coding: utf-8; -*-
// Created on 13. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-13 14:32:58 krylon>

package logdomain

//go:generate stringer -type=ID

// ID signifies a part of the application that wants to write to the log.
type ID uint8

const (
	Database ID = iota
	DBPool
	Scanner
	Player
	Terminal
)

// All returns a slice of all valid ID values.
func All() []ID {
	return []ID{
		Database,
		DBPool,
		Scanner,
		Player,
		Terminal,
	}
} // func All() []ID
