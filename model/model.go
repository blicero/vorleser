// /home/krylon/go/src/github.com/blicero/vorleser/model/model.go
// -*- mode: go; coding: utf-8; -*-
// Created on 13. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-13 15:48:13 krylon>

package model

type File struct {
	ID   int64
	Path string
}

type Book struct {
	ID     int64
	Title  string
	Author string
}
