// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package signup

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
)

// pageID defines the home page id.
const pageID = "HOME_PAGE"

type (
	C = layout.Context
	D = layout.Dimensions
)

// Page holds the state for a page demonstrating the features of
// the NavDrawer component.
type Page struct {
	widget.List
	*router.Router
}

// New constructs a Page with the provided router.
func New(router *router.Router) *Page {
	return &Page{
		Router: router,
	}
}

func (p *Page) NavItem() component.NavItem {
	return component.NavItem{}
}

func (p *Page) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "THIS IS THE SIGNUP PAGE").Layout)
			}),
		)
	})
}
