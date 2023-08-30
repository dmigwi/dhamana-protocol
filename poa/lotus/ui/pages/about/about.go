package about

import (
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/assets"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"

	alo "gioui.org/example/component/applayout"
)

// pageID defines the about page id.
const pageID = utils.ABOUT_PAGE_ID

type (
	C = layout.Context
	D = layout.Dimensions
)

// AboutPage holds the state for a page demonstrating the features of
// the AppBar component.
type AboutPage struct {
	eliasCopyButton, chrisCopyButtonGH, chrisCopyButtonLP widget.Clickable
	widget.List
	*router.Router
}

// New constructs an AboutPage with the provided router.
func New(router *router.Router) *AboutPage {
	return &AboutPage{
		Router: router,
	}
}

func (p *AboutPage) ID() string {
	return pageID
}

func (p *AboutPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *AboutPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *AboutPage) NavItem() component.NavItem {
	return component.NavItem{
		Tag:  p.ID(),
		Name: values.StrAbout,
		Icon: assets.OtherIcon,
	}
}

func (p *AboutPage) OnSwitchTo()   {}
func (p *AboutPage) OnSwitchFrom() {}

func (p *AboutPage) HandleEvents() {}

const (
	sponsorEliasURL          = "https://github.com/sponsors/eliasnaur"
	sponsorChrisURLGitHub    = "https://github.com/sponsors/whereswaldon"
	sponsorChrisURLLiberapay = "https://liberapay.com/whereswaldon/"
)

func (p *AboutPage) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DefaultInset.Layout(gtx, material.Body1(th, `This library implements material design components from https://material.io using https://gioui.org.

If you like this library and work like it, please consider sponsoring Elias and/or Chris!`).Layout)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Elias Naur can be sponsored on GitHub at "+sponsorEliasURL).Layout,
					func(gtx C) D {
						if p.eliasCopyButton.Clicked() {
							clipboard.WriteOp{
								Text: sponsorEliasURL,
							}.Add(gtx.Ops)
						}
						return material.Button(th, &p.eliasCopyButton, "Copy Sponsorship URL").Layout(gtx)
					})
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return alo.DetailRow{}.Layout(gtx,
					material.Body1(th, "Chris Waldon can be sponsored on GitHub at "+sponsorChrisURLGitHub+" and on Liberapay at "+sponsorChrisURLLiberapay).Layout,

					func(gtx C) D {
						if p.chrisCopyButtonGH.Clicked() {
							clipboard.WriteOp{Text: sponsorChrisURLGitHub}.Add(gtx.Ops)
						}
						if p.chrisCopyButtonLP.Clicked() {
							clipboard.WriteOp{Text: sponsorChrisURLLiberapay}.Add(gtx.Ops)
						}
						return alo.DefaultInset.Layout(gtx, func(gtx C) D {
							return layout.Flex{}.Layout(gtx,
								layout.Flexed(.5, material.Button(th, &p.chrisCopyButtonGH, "Copy GitHub URL").Layout),
								layout.Flexed(.5, material.Button(th, &p.chrisCopyButtonLP, "Copy Liberapay URL").Layout),
							)
						})
					})
			}),
		)
	})
}
