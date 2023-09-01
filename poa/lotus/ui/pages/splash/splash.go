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
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// timerLength indicates how long the splash screen should be visible for.
const timerLength = 15 * time.Second

type splashPage struct {
	status  string
	image   *widget.Image
	loading *assets.LoadingIcon

	visibility *time.Timer
}

// New constructs a splashPage with the provided router.
func New() *splashPage {
	sequence := gween.NewSequence(
		gween.New(0, 2, 0.05, ease.Linear),
	)
	sequence.SetLoop(-1)

	return &splashPage{
		status: values.StrLoading,
		image: &widget.Image{
			Src: paint.NewImageOp(assets.SplashImage),
			Fit: widget.Contain,
		},
		loading: assets.NewLoadingIcon(),
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
	s.loading.Start()
	s.visibility = time.NewTimer(timerLength)
	return s.visibility.C
}

// StopTimer stops the timer if its set and running.
func (s *splashPage) StopTimer() {
	s.loading.Stop()
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
			lbl.Font.Weight = font.SemiBold
			lbl.Alignment = text.Middle
			lbl.Color = utils.SecondaryColor
			return lbl.Layout(gtx)
		}),
		layout.Rigid(layout.Spacer{Height: unit.Dp(15)}.Layout),
		layout.Rigid(func(gtx C) D {
			return layout.Center.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Max.X = 30
				gtx.Constraints.Max.Y = 30
				return s.loading.Layout(gtx)
			})
		}),
	)
}
