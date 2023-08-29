// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package pages

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/about"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/account"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/feedback"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/home"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func Loop(w *app.Window) error {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	var ops op.Ops

	pg := router.NewRouter()
	pg.Register(0, home.New(&pg))
	pg.Register(1, account.New(&pg))
	pg.Register(2, feedback.New(&pg))
	pg.Register(3, about.New(&pg))

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				gtx := layout.NewContext(&ops, e)
				pg.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}
