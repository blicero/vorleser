// /home/krylon/go/src/github.com/blicero/vorleser/scanner/scanner.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 13:56:58 krylon>

// Package scanner provides traversal of directory trees, looking for
// audio files.
package scanner

import (
	"io/fs"
	"log"
	"os"
	"regexp"

	"github.com/blicero/krylib"
	"github.com/blicero/vorleser/common"
	"github.com/blicero/vorleser/database"
	"github.com/blicero/vorleser/logdomain"
	"github.com/blicero/vorleser/model"
)

var audioPat = regexp.MustCompile("(?i)[.](mp3|ogg|opus|flac|m4a|m4b|mpc)$")

// Scanner traverses directory trees, looking up and down for audio files.
type Scanner struct {
	log *log.Logger
	db  *database.Database
}

// Create returns a fresh Scanner.
func Create() (*Scanner, error) {
	var (
		err error
		s   = new(Scanner)
	)

	if s.log, err = common.GetLogger(logdomain.Scanner); err != nil {
		return nil, err
	} else if s.db, err = database.Open(common.DbPath); err != nil {
		s.log.Printf("[CRITICAL] Cannot open database at %s: %s\n",
			common.DbPath,
			err.Error())
		return nil, err
	}

	return s, nil
} // func Create() (*Scanner, error)

// ScanFolder scans the given Folder for audio files. The function assumes that
// the Folder has been added to the Database already.
func (s *Scanner) ScanFolder(folder *model.Folder) error {
	s.log.Printf("[TRACE] About to scan folder %s\n", folder.Path)

	var (
		err        error
		fileSystem = os.DirFS(folder.Path)
	)

	if err = fs.WalkDir(fileSystem, ".", s.visitFolder); err != nil {
		s.log.Printf("[ERROR] Failed to scan %s: %s\n",
			folder.Path,
			err.Error())
		return err
	}

	return nil
} // func (s *Scanner) ScanFolder(folder *model.Folder) error

func (s *Scanner) visitFolder(path string, d fs.DirEntry, ex error) error {
	var (
		err  error
		file *model.File
	)

	return krylib.ErrNotImplemented
} // func (s *Scanner) visitFolder(path string, d fs.DirEntry, ex error) error
