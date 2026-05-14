// /home/krylon/go/src/github.com/blicero/vorleser/database/02_folder_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 14. 05. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 12:09:03 krylon>

package database

import (
	"fmt"
	"testing"

	"github.com/blicero/vorleser/model"
)

const folderCnt = 10

var folders []*model.Folder

func TestFolderAdd(t *testing.T) {
	if tdb == nil {
		t.SkipNow()
	}

	folders = make([]*model.Folder, folderCnt)

	for i := range folderCnt {
		var (
			err    error
			folder = &model.Folder{
				Path: fmt.Sprintf("/vault/audio/books_%02d",
					i),
			}
		)

		if err = tdb.FolderAdd(folder); err != nil {
			t.Fatalf("Failed to add Folder %s to Database: %s",
				folder.Path,
				err.Error())
		} else {
			folders[i] = folder
		}
	}
} // func TestFolderAdd(t *testing.T)
