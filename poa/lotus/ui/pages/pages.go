// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package pages

import (
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget/material"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/about"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/account"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/feedback"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/home"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/pages/splash"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

func loadFontCollection() ([]font.FontFace, error) {
	// universal fonts from https://github.com/satbyy/go-noto-universal

	goNotoKurrentRegularTTF, err := assets.GetFont("GoNotoKurrent-Regular.ttf")
	if err != nil {
		return nil, err
	}

	goNotoKurrentRegular, err := opentype.Parse(goNotoKurrentRegularTTF)
	if err != nil {
		return nil, err
	}

	goNotoKurrentBoldTTF, err := assets.GetFont("GoNotoKurrent-Bold.ttf")
	if err != nil {
		return nil, err
	}

	goNotoKurrentBold, err := opentype.Parse(goNotoKurrentBoldTTF)
	if err != nil {
		return nil, err
	}

	fontCollection := []font.FontFace{}
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{}, Face: goNotoKurrentRegular})
	fontCollection = append(fontCollection, font.FontFace{Font: font.Font{Weight: font.Bold}, Face: goNotoKurrentBold})
	return fontCollection, nil
}

func Loop(w *app.Window) error {
	th := material.NewTheme()
	fontCollection, _ := loadFontCollection()
	th.Shaper = text.NewShaper(text.NoSystemFonts(), text.WithCollection(fontCollection))
	var ops op.Ops

	pg := router.NewRouter()
	pg.Register(home.New(pg), account.New(pg), feedback.New(pg), about.New(pg))

	// Load the splash page and set it as the current.
	splashPg := splash.New()
	pg.SetTemporaryPage(splashPg)

	go func() {
		select {
		case <-splashPg.StartTimer():
			pg.DisableTemporaryPage()
			splashPg.StopTimer()

			// Invalidate the current frame, triggering page redrawing.
			w.Invalidate()
		}
	}()

	for {
		select {
		case e := <-w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err

			case system.FrameEvent:
				pg.ProcessEvents()
				gtx := layout.NewContext(&ops, e)

				pg.Layout(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}
