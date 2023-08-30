// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package router

import (
	"log"
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

	*component.ModalNavDrawer
	NavAnim component.VisibilityAnimation
	menuBar *component.AppBar
	*component.ModalLayer
	NonModalDrawer, BottomBar bool
}

func NewRouter() *Router {
	modal := component.NewModal()

	nav := component.NewNav(utils.AppName, values.StrAppDescription)
	modalNav := component.ModalNavFrom(&nav, modal)

	bar := component.NewAppBar(modal)
	bar.NavigationIcon = assets.MenuIcon

	na := component.VisibilityAnimation{
		State:    component.Invisible,
		Duration: time.Millisecond * 250,
	}

	return &Router{
		ModalLayer:     modal,
		ModalNavDrawer: modalNav,
		menuBar:        bar,
		NavAnim:        na,
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
			r.menuBar.Title = p.NavItem().Name
		}
		r.registeredPages[p.ID()] = p
		r.ModalNavDrawer.AddNavItem(p.NavItem())
	}
}

// func (r *Router) SwitchTo(tag string) {
// 	p, ok := r.pages[tag]
// 	if !ok {
// 		return
// 	}
// 	navItem := p.NavItem()
// 	r.current = tag
// 	r.AppBar.Title = navItem.Name
// 	r.AppBar.SetActions(p.Actions(), p.Overflow())
// }

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

	if r.ModalNavDrawer.NavDestinationChanged() {
		p, ok := r.registeredPages[r.ModalNavDrawer.CurrentNavDestination()]
		if ok {
			r.currentPage = p
			r.menuBar.Title = p.NavItem().Name
		}
	}

	r.currentPage.HandleEvents()
}

func (r *Router) AddBackButton() {
	// r.AppBar.NavigationIcon = assets.BackButton
}

// Layout handles ploting the componnets of the current page by calling the
// actual page Layout method.
func (r *Router) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	for _, event := range r.menuBar.Events(gtx) {
		switch event := event.(type) {
		case component.AppBarNavigationClicked:
			if r.NonModalDrawer {
				r.NavAnim.ToggleVisibility(gtx.Now)
			} else {
				r.ModalNavDrawer.Appear(gtx.Now)
				r.NavAnim.Disappear(gtx.Now)
			}
		case component.AppBarContextMenuDismissed:
			log.Printf("Context menu dismissed: %v", event)
		case component.AppBarOverflowActionClicked:
			log.Printf("Overflow action selected: %v", event)
		}
	}

	paint.Fill(gtx.Ops, th.Palette.Bg)
	content := layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Max.X /= 2
				return r.NavDrawer.Layout(gtx, th, &r.NavAnim)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if r.currentPage == nil {
					return layout.Dimensions{}
				}
				return r.currentPage.Layout(gtx, th)
			}),
		)
	})

	bar := layout.Rigid(func(gtx layout.Context) layout.Dimensions {
		return r.menuBar.Layout(gtx, th, "Menu", "Actions")
	})

	flex := layout.Flex{Axis: layout.Vertical}
	if r.BottomBar {
		flex.Layout(gtx, content, bar)
	} else {
		flex.Layout(gtx, bar, content)
	}
	r.ModalLayer.Layout(gtx, th)
	return layout.Dimensions{Size: gtx.Constraints.Max}
}
