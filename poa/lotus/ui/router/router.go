// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package router

import (
	"time"

	"gioui.org/layout"
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
	// NavItem returns the navigation fields linked to each page.
	NavItem() *component.NavItem
}

type Router struct {
	// current defines the page ID of the current page in view.
	currentPage Page
	// tempPage defines a page that is used until the window is invalidated
	// moving to the regular pages. They include but not limited to: splash and
	// sign up pages.
	tempPage Page
	// pageStack track of the order in which pages are visited for the purpose
	// of backtracking to a previous page in the order in which they were
	// accessed. It holds the page ID strings.
	pageStack []string
	// registeredPages maintains references to already initialised pages mapped
	// to their respective page IDs
	registeredPages map[string]Page

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

// SetTemporaryPage sets a one-time-use then throw-away page after the current
// window is invalidated. Temp pages should not be registered or pushed to the
// page stack.
func (r *Router) SetTemporaryPage(p Page) {
	if p != nil {
		r.tempPage = p
	}
}

// DisableTemporaryPage disables the currently set temporary page enabling the
// regular pages already registered to be displayed.
func (r *Router) DisableTemporaryPage() {
	r.tempPage = nil
}

// Register receives all the pages that will be regulary accessed.
func (r *Router) Register(pages ...Page) {
	if r.pageStack != nil {
		// Page stack has already been populated
		return
	}

	r.registeredPages = make(map[string]Page, len(pages))
	for _, p := range pages {
		if r.currentPage == nil {
			r.currentPage = p
		}

		r.registeredPages[p.ID()] = p
		navItem := p.NavItem()
		if navItem != nil {
			r.navDrawer.AddNavItem(*navItem)
		}
	}
}

func (r *Router) SetBackNavButton() {
	r.menuBar.NavigationIcon = assets.BackIcon
}

// OnDisplay appends the passed page onto of the current page stack. It helps
// maintain page access history for future backtracking if need be.
func (r *Router) OnDisplay(pageID interface{}) {
	str, _ := pageID.(string)
	p, ok := r.registeredPages[str]
	if !ok {
		return
	}

	if r.currentPage != nil && p.ID() == r.currentPage.ID() {
		return
	}

	// execute the previous page OnSwitchFrom() method.
	r.currentPage.OnSwitchFrom()

	r.currentPage = p
	r.pageStack = append(r.pageStack, p.ID())

	// execute the current page OnSwitchTo() method.
	r.currentPage.OnSwitchTo()
}

// OnDisplayNew clears the current page stuck before pushing the passed page
// as the only one in the stack.
func (r *Router) OnDisplayNew(pageID interface{}) {
	// empty out the previous pages.
	r.pageStack = r.pageStack[:0]

	r.OnDisplay(pageID)
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
				strPageID := r.pageStack[count-1]
				r.currentPage = r.registeredPages[strPageID]
			}
		} else {
			r.navDrawer.Appear(time.Now())
		}
	}

	if r.navDrawer.NavDestinationChanged() {
		strPageID, _ := r.navDrawer.CurrentNavDestination().(string)
		p, ok := r.registeredPages[strPageID]
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
	// If a temporary page is set display it instead.
	if r.tempPage != nil {
		return r.tempPage.Layout(gtx, th)
	}

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
