// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.
package main

import (
	"log"
	"os"
	"runtime/trace"

	"gioui.org/app"
	"gioui.org/unit"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
)

func main() {
	minSizeX := unit.Dp(375)
	minSizeY := unit.Dp(600)
	maxSizeX := unit.Dp(500)
	maxSizeY := unit.Dp(1000)

	w := app.NewWindow(
		app.Title(utils.AppName),
		app.MinSize(minSizeX, minSizeY),
		app.Size(minSizeX, minSizeY),
		app.MaxSize(maxSizeX, maxSizeY),
		app.PortraitOrientation.Option(),
		app.NavigationColor(utils.HighlightColor),
		app.StatusColor(utils.DarkPriColor),
	)

	file := "file.trace"
	f, err := os.Create(file)
	if err != nil {
		log.Fatal(err)
	}

	trace.Start(f)

	cleanUp := func() {
		f.Close()
		trace.Stop()
	}

	go func() {
		if err := pages.Loop(w); err != nil {
			log.Fatal(err)
		}
		cleanUp()
		os.Exit(0)
	}()

	app.Main()
}
