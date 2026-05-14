// /home/krylon/go/src/github.com/blicero/vorleser/model/model.go
// -*- mode: go; coding: utf-8; -*-
// Created on 13. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 11:37:18 krylon>

package model

import (
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Book is a ... well, a book. Really.
type Book struct {
	ID          int64
	CurrentFile int64
	Title       string
	Author      string
	URL         *url.URL
}

// Folder is a directory containing Files.
type Folder struct {
	ID       int64
	Path     string
	LastScan time.Time
}

// File is an audio file.
type File struct {
	ID         int64
	BookID     int64
	FolderID   int64
	Path       string
	Title      string
	Ord        []int64
	Position   int64
	LastPlayed time.Time
}

// DisplayTitle returns the File's title, or its filename, if no title is known.
func (f *File) DisplayTitle() string {
	if f.Title != "" {
		return f.Title
	}

	return path.Base(f.Path)
} // func (f *File) DisplayTitle() string

// GetParentFolder returns the name of the Folder the file lives in,
// i.e. basename(dirname(path))
func (f *File) GetParentFolder() string {
	return filepath.Base(filepath.Dir(f.Path))
} // func (f *File) GetParentFolder() string

// PathURL returns the File's path as a file:// URL.
// Intended for use with DBus interfaces
func (f *File) PathURL() string {
	return "file://" + url.PathEscape(f.Path)
} // func (f *File) PathURL() string

// OrdString returns a textual representation of the File's Ord.
func (f *File) OrdString() string {
	if len(f.Ord) == 0 {
		return ""
	}

	var slist = make([]string, len(f.Ord))

	for idx, val := range f.Ord {
		slist[idx] = strconv.FormatInt(val, 10)
	}

	return strings.Join(slist, ",")
} // func (f *File) OrdString() string
