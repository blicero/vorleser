// /home/krylon/go/src/github.com/blicero/vorleser/database/helpers.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 11:22:22 krylon>

package database

// valOrNil returns nil if the argument is equal to the zero value of its type,
// otherwise it returns the argument verbatim.
func valOrNil[T comparable](v T) any {
	var zero T

	if v == zero {
		return nil
	}

	return v
} // func valOrNil[T any](v T) any
