// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package signup

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// signupPage holds the state for a page demonstrating the features of
// the NavDrawer component.
type signupPage struct {
	*router.Router
	nameInput    component.TextField
	addressInput component.TextField
	priceInput   component.TextField
	tweetInput   component.TextField
	numberInput  component.TextField
}

// New constructs a SignupPage with the provided router.
func New(router *router.Router) *signupPage {
	return &signupPage{
		Router: router,
	}
}

func (s *signupPage) ID() string {
	return utils.SIGNUP_PAGE_ID
}

func (s *signupPage) OnSwitchTo()   {}
func (s *signupPage) OnSwitchFrom() {}
func (s *signupPage) HandleEvents() {}

// NavItem returns nil because the signup page doesn't support navigation bar
// fields.
func (s *signupPage) NavItem() *component.NavItem {
	return nil
}

func (p *signupPage) Layout(gtx C, th *material.Theme) D {
	return layout.Inset{
		Right: unit.Dp(30),
		Left:  unit.Dp(30),
	}.Layout(gtx, func(gtx C) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(material.Body2(th, "THIS IS THE SIGNUP PAGE").Layout),
			layout.Rigid(func(gtx C) D {
				mTh := *th
				mTh.Palette.Fg = utils.PrimaryColor
				mTh.Palette.ContrastBg = utils.PrimaryColor
				// p.nameInput.
				return p.nameInput.Layout(gtx, &mTh, "Paste Wallet Private Key")
			}),
		)
	})
}
