// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package feedback

import (
	"image/color"
	"unicode"

	// alo "gioui.org/example/component/applayout"
	"gioui.org/example/component/icon"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/router"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils/values"
)

// pageID defines the feedback page id.
const pageID = utils.FEEDBACK_PAGE_ID

type (
	C = layout.Context
	D = layout.Dimensions
)

// FeedbackPage holds the state for a page demonstrating the features of
// the TextField component.
type FeedbackPage struct {
	inputAlignment                                               text.Alignment
	inputAlignmentEnum                                           widget.Enum
	nameInput, addressInput, priceInput, tweetInput, numberInput component.TextField
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
	return pageID
}

func (p *FeedbackPage) Actions() []component.AppBarAction {
	return []component.AppBarAction{}
}

func (p *FeedbackPage) Overflow() []component.OverflowAction {
	return []component.OverflowAction{}
}

func (p *FeedbackPage) NavItem() component.NavItem {
	return component.NavItem{
		Tag:  p.ID(),
		Name: values.StrFeedback,
		Icon: icon.EditIcon,
	}
}

func (p *FeedbackPage) OnSwitchTo()   {}
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
				p.nameInput.Alignment = p.inputAlignment
				return p.nameInput.Layout(gtx, th, "Name")
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "Responds to hover events.").Layout)
			}),
			layout.Rigid(func(gtx C) D {
				p.addressInput.Alignment = p.inputAlignment
				return p.addressInput.Layout(gtx, th, "Address")
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "Label animates properly when you click to select the text field.").Layout)
			}),
			layout.Rigid(func(gtx C) D {
				p.priceInput.Prefix = func(gtx C) D {
					th := *th
					th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
					return material.Label(&th, th.TextSize, "$").Layout(gtx)
				}
				p.priceInput.Suffix = func(gtx C) D {
					th := *th
					th.Palette.Fg = color.NRGBA{R: 100, G: 100, B: 100, A: 255}
					return material.Label(&th, th.TextSize, ".00").Layout(gtx)
				}
				p.priceInput.SingleLine = true
				p.priceInput.Alignment = p.inputAlignment
				return p.priceInput.Layout(gtx, th, "Price")
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "Can have prefix and suffix elements.").Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if err := func() string {
					for _, r := range p.numberInput.Text() {
						if !unicode.IsDigit(r) {
							return "Must contain only digits"
						}
					}
					return ""
				}(); err != "" {
					p.numberInput.SetError(err)
				} else {
					p.numberInput.ClearError()
				}
				p.numberInput.SingleLine = true
				p.numberInput.Alignment = p.inputAlignment
				return p.numberInput.Layout(gtx, th, "Number")
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "Can be validated.").Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if p.tweetInput.TextTooLong() {
					p.tweetInput.SetError("Too many characters")
				} else {
					p.tweetInput.ClearError()
				}
				p.tweetInput.CharLimit = 128
				p.tweetInput.Helper = "Tweets have a limited character count"
				p.tweetInput.Alignment = p.inputAlignment
				return p.tweetInput.Layout(gtx, th, "Tweet")
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "Can have a character counter and help text.").Layout)
			}),
			layout.Rigid(func(gtx C) D {
				if p.inputAlignmentEnum.Changed() {
					switch p.inputAlignmentEnum.Value {
					case layout.Start.String():
						p.inputAlignment = text.Start
					case layout.Middle.String():
						p.inputAlignment = text.Middle
					case layout.End.String():
						p.inputAlignment = text.End
					default:
						p.inputAlignment = text.Start
					}
					op.InvalidateOp{}.Add(gtx.Ops)
				}
				return layout.Inset{}.Layout(
					gtx,
					func(gtx C) D {
						return layout.Flex{
							Axis: layout.Vertical,
						}.Layout(
							gtx,
							layout.Rigid(func(gtx C) D {
								return material.Body2(th, "Text Alignment").Layout(gtx)
							}),
							layout.Rigid(func(gtx C) D {
								return layout.Flex{
									Axis: layout.Vertical,
								}.Layout(
									gtx,
									layout.Rigid(func(gtx C) D {
										return material.RadioButton(
											th,
											&p.inputAlignmentEnum,
											layout.Start.String(),
											"Start",
										).Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return material.RadioButton(
											th,
											&p.inputAlignmentEnum,
											layout.Middle.String(),
											"Middle",
										).Layout(gtx)
									}),
									layout.Rigid(func(gtx C) D {
										return material.RadioButton(
											th,
											&p.inputAlignmentEnum,
											layout.End.String(),
											"End",
										).Layout(gtx)
									}),
								)
							}),
						)
					},
				)
			}),
			layout.Rigid(func(gtx C) D {
				return layout.Inset{}.Layout(gtx, material.Body2(th, "This text field implementation was contributed by Jack Mordaunt. Thanks Jack!").Layout)
			}),
		)
	})
}
