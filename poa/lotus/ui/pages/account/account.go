// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package account

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	alo "gioui.org/example/component/applayout"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"
)

// pageID defines the account page id.
const pageID = utils.ACCOUNT_PAGE_ID

type (
	C = layout.Context
	D = layout.Dimensions
)

// AccountPage holds the state for a page demonstrating the features of
// the NavDrawer component.
type AccountPage struct {
	nonModalDrawer widget.Bool
	widget.List
	*router.Router
}

// New constructs a AccountPage with the provided router.
func New(router *router.Router) *AccountPage {
	return &AccountPage{
		Router: router,
	}
}

func (p *AccountPage) ID() string {
	return pageID
}

func (p *AccountPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *AccountPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *AccountPage) NavItem() component.NavItem {
	return component.NavItem{
		Tag:  p.ID(),
		Name: values.StrAccount,
		Icon: assets.SettingsIcon,
	}
}

func (p *AccountPage) OnSwitchTo()   {}
func (p *AccountPage) OnSwitchFrom() {}

func (p *AccountPage) HandleEvents() {}

func (p *AccountPage) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DefaultInset.Layout(gtx, material.Body1(th, `The nav drawer widget provides a consistent interface element for navigation.

	The controls below allow you to see the various features available in our Navigation Drawer implementation.`).Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Use non-modal drawer").Layout,
					func(gtx C) D {
						if p.nonModalDrawer.Changed() {
							p.Router.NonModalDrawer = p.nonModalDrawer.Value
							if p.nonModalDrawer.Value {
								p.Router.NavAnim.Appear(gtx.Now)
							} else {
								p.Router.NavAnim.Disappear(gtx.Now)
							}
						}
						return material.Switch(th, &p.nonModalDrawer, "Use Non-Modal Navigation Drawer").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Drag to Close").Layout,
					material.Body2(th, "You can close the modal nav drawer by dragging it to the left.").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Touch Scrim to Close").Layout,
					material.Body2(th, "You can close the modal nav drawer touching anywhere in the translucent scrim to the right.").Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Bottom content anchoring").Layout,
					material.Body2(th, "If you toggle support for the bottom app bar in the App Bar settings, nav drawer content will anchor to the bottom of the drawer area instead of the top.").Layout)
			}),
		)
	})
}
