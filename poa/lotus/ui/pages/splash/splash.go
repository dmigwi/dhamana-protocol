// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package splash

import (
	"time"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// timerLength indicates how long the splash screen should be visible for.
const timerLength = 10 * time.Second

type splashPage struct {
	status string
	image  *widget.Image

	visibility *time.Timer
}

// New constructs a splashPage with the provided router.
func New() *splashPage {
	return &splashPage{
		status: values.StrLoading,
		image: &widget.Image{
			Src: paint.NewImageOp(assets.SplashImage),
			Fit: widget.Contain,
		},
	}
}

func (s *splashPage) ID() string {
	return utils.SPLASH_PAGE_ID
}

func (s *splashPage) OnSwitchTo()   {}
func (s *splashPage) OnSwitchFrom() {}
func (s *splashPage) HandleEvents() {}

// NavItem returns nil because the splash screen doesn't support navigation bar
// fields.
func (s *splashPage) NavItem() *component.NavItem {
	return nil
}

// StartTimer initiates the timer.
func (s *splashPage) StartTimer() <-chan time.Time {
	s.visibility = time.NewTimer(timerLength)
	return s.visibility.C
}

// StopTimer stops the timer if its set and running.
func (s *splashPage) StopTimer() {
	if s.visibility != nil {
		s.visibility.Stop()
		s.visibility = nil
	}
}

func (s *splashPage) Layout(gtx C, th *material.Theme) D {
	// paint the background with the primary color.
	paint.Fill(gtx.Ops, utils.PrimaryColor)

	return layout.Flex{
		Axis:      layout.Vertical,
		Alignment: layout.Middle,
	}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				return s.image.Layout(gtx)
			})
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(20)}.Layout),
		layout.Rigid(func(gtx C) D {
			lbl := material.Label(th, unit.Sp(14), utils.AppVersion.String())
			lbl.Font.Weight = font.Bold
			lbl.Alignment = text.Middle
			lbl.Color = utils.SecondaryColor
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				lbl := material.Label(th, unit.Sp(20), s.status)
				lbl.Font.Weight = font.Bold
				lbl.Color = utils.SecondaryColor
				return lbl.Layout(gtx)
			})
		}),
	)
}
