// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package router

import (
	"time"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"
)

// Page defines the interface all pages should implement to navigation from on
// page to another one.
type Page interface {
	// ID returns the unique page identity string
	ID() string
	// OnSwitchTo implements the functionality to be invoked once the page loads
	// and before its actual functionality is invoked.
	OnSwitchTo()
	// OnSwitchFrom implements the functionality be invoked before the page
	// exists.
	OnSwitchFrom()
	// HandleEvents handles all the events triggered by the components in the
	// current page.
	HandleEvents()
	// Layout arranges the ui components in the current page.
	Layout(gtx layout.Context, th *material.Theme) layout.Dimensions
	NavItem() component.NavItem
}

type Router struct {
	// current defines the page ID of the current page in view.
	currentPage Page
	// pageStack track of the order in which pages are visited for the purpose
	// of backtracking to a previous page in the order in which they were
	// accessed.
	pageStack []Page
	// registeredPages maintains references to already initialised pages mapped
	// to their respective page IDs
	registeredPages map[interface{}]Page

	backbutton *widget.Clickable

	navDrawer  *component.ModalNavDrawer
	menuBar    *component.AppBar
	modalLayer *component.ModalLayer
}

type (
	C = layout.Context
	D = layout.Dimensions
)

func NewRouter() *Router {
	modal := component.NewModal()
	bar := component.NewAppBar(modal)
	bar.NavigationIcon = assets.MenuIcon
	nav := component.NewNav(utils.AppName, values.StrAppDescription)

	return &Router{
		modalLayer: modal,
		navDrawer:  component.ModalNavFrom(&nav, modal),
		menuBar:    bar,
	}
}

func (r *Router) Register(pages ...Page) {
	if r.pageStack != nil {
		// Page stack has already been populated
		return
	}

	r.registeredPages = make(map[interface{}]Page, len(pages))
	for _, p := range pages {
		if r.currentPage == nil {
			r.currentPage = p
		}

		r.registeredPages[p.ID()] = p
		r.navDrawer.AddNavItem(p.NavItem())
	}
}

func (r *Router) SetBackNavButton() {
	r.menuBar.NavigationIcon = assets.BackIcon
}

// OnDisplay appends the passed page onto of the current page stack. It helps
// maintain page access history for future backtracking if need be.
func (r *Router) OnDisplay(pageID interface{}) {
	p, ok := r.registeredPages[pageID]
	if !ok {
		return
	}

	if r.currentPage != nil && p.ID() == r.currentPage.ID() {
		return
	}

	r.currentPage = p
	r.pageStack = append(r.pageStack, p)
}

// OnDisplayNew clears the current page stuck before pushing the passed page
// as the only one in the stack.
func (r *Router) OnDisplayNew(pageID interface{}) {
	p, ok := r.registeredPages[pageID]
	if !ok {
		return
	}

	if r.currentPage != nil && p.ID() == r.currentPage.ID() {
		// Do not push the current page onto itself.
		return
	}

	r.currentPage = p
	// empty out the previous pages.
	r.pageStack = append(r.pageStack[:0], p)
}

func (r *Router) ProcessEvents() {
	// if current page isn't sent ignore processing the events.
	if r.currentPage == nil {
		return
	}

	if r.menuBar.NavigationButton.Clicked() {
		if r.menuBar.NavigationIcon == assets.BackIcon {
			count := len(r.pageStack)
			if count > 0 {
				r.currentPage = r.pageStack[count-1]
			}
		} else {
			r.navDrawer.Appear(time.Now())
		}
	}

	if r.navDrawer.NavDestinationChanged() {
		p, ok := r.registeredPages[r.navDrawer.CurrentNavDestination()]
		if ok {
			r.currentPage = p
		}
	}

	r.currentPage.HandleEvents()
}

func (r *Router) AddBackButton() {
	r.menuBar.NavigationIcon = assets.BackIcon
}

// Layout handles ploting the componnets of the current page by calling the
// actual page Layout method.
func (r *Router) Layout(gtx C, th *material.Theme) D {
	paint.Fill(gtx.Ops, th.Palette.Bg)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		// Top bar
		layout.Rigid(func(gtx C) D {
			r.menuBar.Title = r.currentPage.NavItem().Name
			return r.menuBar.Layout(gtx, th, "Menu", "Actions")
		}),

		// Navigation bar
		layout.Rigid(func(gtx C) D {
			return r.modalLayer.Layout(gtx, th)
		}),

		// Page content.
		layout.Flexed(1, func(gtx C) D {
			if r.currentPage == nil {
				return D{}
			}
			return r.currentPage.Layout(gtx, th)
		}),
	)
}
