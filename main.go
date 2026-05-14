// /home/krylon/go/src/github.com/blicero/vorleser/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 09. 03. 2026 by Benjamin Walkenhorst
// (c) 2026 Benjamin Walkenhorst
// Time-stamp: <2026-05-14 09:49:16 krylon>

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/blicero/vorleser/common"
)

func main() {
	var (
		err     error
		ticker  *time.Ticker
		addr    string
		profOut string
		sigQ    = make(chan os.Signal, 1)
	)

	fmt.Printf("%s %s, built on %s\n",
		common.AppName,
		common.Version,
		common.BuildStamp.Format(common.TimestampFormat))

	flag.StringVar(&profOut, "prof", "", "if non-empty, write profiling information to the named file")
	flag.Parse()

	if profOut != "" {
		var profH *os.File

		fmt.Printf("Writing profiling data to %s\n",
			profOut)

		if profH, err = os.Create(profOut); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Failed to open %s: %s\n",
				profOut,
				err.Error())
			os.Exit(1)
		}

		defer profH.Close() // nolint: errcheck

		if err = pprof.StartCPUProfile(profH); err != nil {
			fmt.Fprintf(
				os.Stderr,
				"Error starting CPU profile: %s\n",
				err.Error())
			os.Exit(1)
		}

		defer pprof.StopCPUProfile()
	}

	ticker = time.NewTicker(common.ActiveTimeout)
	defer ticker.Stop()

	signal.Notify(sigQ, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
		case s := <-sigQ:
			fmt.Fprintf(
				os.Stderr,
				"Caught signal: %s\n",
				s)
			return
		}
	}
}
