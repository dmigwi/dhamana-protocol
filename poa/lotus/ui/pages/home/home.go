// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package home

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

// pageID defines the home page id.
const pageID = utils.HOME_PAGE_ID

type (
	C = layout.Context
	D = layout.Dimensions
)

// HomePage holds the state for a page demonstrating the features of
// the Menu component.
type HomePage struct {
	*router.Router
	widget.List
}

// New constructs a Page with the provided router.
func New(router *router.Router) *HomePage {
	return &HomePage{
		Router: router,
	}
}

func (p *HomePage) ID() string {
	return pageID
}

func (p *HomePage) NavItem() component.NavItem {
	return component.NavItem{
		Tag:  p.ID(),
		Name: values.StrHome,
		Icon: assets.HomeIcon,
	}
}

func (p *HomePage) OnSwitchTo() {}

func (p *HomePage) OnSwitchFrom() {}

func (p *HomePage) HandleEvents() {}

func (p *HomePage) Layout(gtx C, th *material.Theme) D {
	return layout.Flex{}.Layout(gtx,
		layout.Rigid(func(gtx C) D {
			return layout.Inset{}.Layout(gtx, material.Body2(th, "THIS IS THE HOME PAGE").Layout)
		}),
	)
}
