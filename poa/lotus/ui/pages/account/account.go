// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package account

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

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

func (p *AccountPage) NavItem() component.NavItem {
	return component.NavItem{
		Tag:  p.ID(),
		Name: values.StrAccount,
		Icon: assets.SettingsIcon,
	}
}

func (p *AccountPage) OnSwitchTo() {}

func (p *AccountPage) OnSwitchFrom() {}

func (p *AccountPage) HandleEvents() {}

func (p *AccountPage) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "THIS IS THE ACCOUNT PAGE").Layout)
			}),
		)
	})
}
