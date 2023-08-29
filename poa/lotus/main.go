// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.
package main

import (
	"flag"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func main() {
	flag.Parse()
	go func() {
		w := app.NewWindow()
		if err := pages.Loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}
