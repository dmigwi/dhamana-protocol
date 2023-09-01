// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package feedback

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

// FeedbackPage holds the state for a page demonstrating the features of
// the TextField component.
type FeedbackPage struct {
	widget.List
	*router.Router
}

// New constructs a FeedbackPage with the provided router.
func New(router *router.Router) *FeedbackPage {
	return &FeedbackPage{
		Router: router,
	}
}

func (p *FeedbackPage) ID() string {
	return utils.FEEDBACK_PAGE_ID
}

func (p *FeedbackPage) NavItem() *component.NavItem {
	return &component.NavItem{
		Tag:  p.ID(),
		Name: values.StrFeedback,
		Icon: assets.EditIcon,
	}
}

func (p *FeedbackPage) OnSwitchTo() {}

func (p *FeedbackPage) OnSwitchFrom() {}

func (p *FeedbackPage) HandleEvents() {}

func (p *FeedbackPage) Layout(gtx C, th *material.Theme) D {
	p.List.Axis = layout.Vertical
	return material.List(th, &p.List).Layout(gtx, 1, func(gtx C, _ int) D {
		return layout.Flex{
			Axis: layout.Vertical,
		}.Layout(
			gtx,
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "THIS IS THE FEEDBACK PAGE").Layout)
			}),
		)
	})
}
