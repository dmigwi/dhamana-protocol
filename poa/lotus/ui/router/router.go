// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package router

import (
	"sync"
	"sync/atomic"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"

	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"
)

// layoutRunner prevents execution of the Router.Layout() method if it is currently
// executing. The previous instance need to complete before starting on new onces.
var layoutRunner uint32

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

	// navDrawer defines the component holding the navigation bar title and fields.
	navDrawer *component.ModalNavDrawer
	// menuBar implements the top bar as defined in material UI guidelines.
	// https://m2.material.io/components/app-bars-top
	// It enables attaching the navigation bar to the current page.
	menuBar *component.AppBar

	// modalLayer is drawn on top of the normal UI allowing it to display extra
	// components like the navigation bar UI or Modals.
	modalLayer *component.ModalLayer

	// mu controls the read and write access to the router fields preventing
	// multipe write that could eventually lead to app crash.
	mu sync.RWMutex
}

type (
	C = layout.Context
	D = layout.Dimensions
)

// NewRouter returns a new instance of the router pointer.
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
		r.mu.Lock()
		r.tempPage = p
		r.mu.Unlock()
	}
}

// DisableTemporaryPage disables the currently set temporary page enabling the
// regular pages already registered to be displayed.
func (r *Router) DisableTemporaryPage() {
	r.mu.Lock()
	r.tempPage = nil
	r.mu.Unlock()
}

// Register receives all the pages that will be regulary accessed.
func (r *Router) Register(pages ...Page) {
	r.mu.Lock()
	defer r.mu.Unlock()

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
	r.mu.Lock()
	r.menuBar.NavigationIcon = assets.BackIcon
	r.mu.Unlock()
}

// OnDisplay appends the passed page onto of the current page stack. It helps
// maintain page access history for future backtracking if need be.
func (r *Router) OnDisplay(pageID interface{}) {
	r.mu.RLock()
	str, _ := pageID.(string)
	p, ok := r.registeredPages[str]
	if !ok {
		r.mu.RUnlock()
		return
	}

	if r.currentPage != nil && p.ID() == r.currentPage.ID() {
		r.mu.RUnlock()
		return
	}

	previousPage := r.currentPage
	currentPage := p
	r.mu.RUnlock()

	r.mu.Lock()
	r.currentPage = p
	r.pageStack = append(r.pageStack, p.ID())
	r.mu.Unlock()

	// execute the previous page OnSwitchFrom() method without mutex protection.
	previousPage.OnSwitchFrom()

	// execute the current page OnSwitchTo() method without mutex protection.
	currentPage.OnSwitchTo()
}

// OnDisplayNew clears the current page stuck before pushing the passed page
// as the only one in the stack.
func (r *Router) OnDisplayNew(pageID interface{}) {
	r.mu.Lock()
	// empty out the previous pages.
	r.pageStack = r.pageStack[:0]
	r.mu.Unlock()

	r.OnDisplay(pageID)
}

func (r *Router) ProcessEvents() {
	r.mu.Lock()
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
	currentPage := r.currentPage
	r.mu.Unlock()

	// execute the HandleEvents() method without mutex protection.
	currentPage.HandleEvents()
}

// AddBackButton set the bar navigation button to back arrow enabling pages
// backtracking.
func (r *Router) AddBackButton() {
	r.mu.Lock()
	r.menuBar.NavigationIcon = assets.BackIcon
	r.mu.Unlock()
}

// Layout handles ploting the componnets of the current page by calling the
// actual page Layout method.
func (r *Router) Layout(gtx C, th *material.Theme) D {
	// Before we enter into a mutex lock state, check that all the previous
	// Layout() instance are not running. If affirmative return an empty page.
	if !atomic.CompareAndSwapUint32(&layoutRunner, 0, 1) {
		return D{}
	}

	r.mu.Lock()
	defer func() {
		r.mu.Unlock()

		// Reset the layout runner allowing other layout updates to proceed.
		atomic.StoreUint32(&layoutRunner, 0)
	}()

	// If a temporary page is set display it instead.
	if r.tempPage != nil {
		return r.tempPage.Layout(gtx, th)
	}

	return layout.Flex{
		Axis:    layout.Vertical,
		Spacing: layout.SpaceSides,
	}.Layout(gtx,
		// Top bar
		layout.Rigid(func(gtx C) D {
			navItem := r.currentPage.NavItem()
			if navItem != nil {
				r.menuBar.Title = navItem.Name
			}
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
