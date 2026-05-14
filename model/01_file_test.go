// /home/krylon/go/src/github.com/blicero/vorleser/model/01_file_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 09:21:41 krylon>

package model

import "testing"

func TestGetParentFolder(t *testing.T) {
	type testCase struct {
		f File
		p string
	}

	var cases = []testCase{
		{
			f: File{
				Path: "/data/Files/audio/test/test01.mp3",
			},
			p: "test",
		},
	}

	for _, c := range cases {
		if p := c.f.GetParentFolder(); p != c.p {
			t.Errorf(`Unexpected result from GetParentFolder:
Expected:	%s
Got:		%s`,
				c.p,
				p)
		}
	}
} // func TestGetParentFolder(t *testing.T)
