package about

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

type (
	C = layout.Context
	D = layout.Dimensions
)

// AboutPage holds the state for a page demonstrating the features of
// the AppBar component.
type AboutPage struct {
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
	return utils.ABOUT_PAGE_ID
}

func (p *AboutPage) NavItem() *component.NavItem {
	return &component.NavItem{
		Tag:  p.ID(),
		Name: values.StrAbout,
		Icon: assets.OtherIcon,
	}
}

func (p *AboutPage) OnSwitchTo() {}

func (p *AboutPage) OnSwitchFrom() {}

func (p *AboutPage) HandleEvents() {}

func (p *AboutPage) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Alignment: layout.Middle,
			Axis:      layout.Vertical,
		}.Layout(gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "THIS IS THE ABOUT PAGE").Layout)
			}),
		)
	})
}
