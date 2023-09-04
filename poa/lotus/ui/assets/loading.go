// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package assets

import (
	"math"
	"sync"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/dmigwi/dhamana-protocol/poa/lotus/ui/utils"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"
)

// LoadingIcon stores all the fields required to animate the loading icon.
type LoadingIcon struct {
	icon  *widget.Icon
	anime *gween.Sequence
	stop  bool

	mu sync.RWMutex
}

// NewLoadingIcon returns a loading icon component.
func NewLoadingIcon() *LoadingIcon {
	// loadingIcon is an IconVg  generated from a SVG file.
	loadingIconVG := []byte{
		0x89, 0x49, 0x56, 0x47, 0x02, 0x0a, 0x00, 0x50, 0x50, 0xb0, 0xb0, 0xc0, 0x19, 0x87, 0x19, 0x73,
		0x00, 0x7c, 0x54, 0xe9, 0x25, 0x86, 0xa0, 0x21, 0x76, 0x21, 0x71, 0x60, 0xd9, 0x77, 0x60, 0x80,
		0x90, 0x19, 0x86, 0xe1, 0x8e, 0x9c, 0xdd, 0x8f, 0xe9, 0xf5, 0x7b, 0xb0, 0x51, 0x7a, 0x0d, 0x7f,
		0x6c, 0x21, 0x7a, 0x6c, 0x31, 0x74, 0x90, 0x51, 0x84, 0x25, 0x75, 0x94, 0x31, 0x74, 0xe8, 0x78,
		0x20, 0x19, 0x89, 0x19, 0x77, 0xe2, 0xdd, 0x8f, 0x7c, 0xb0, 0xa9, 0x7f, 0x39, 0x7d, 0x91, 0x7e,
		0x8d, 0x7a, 0xc5, 0x7c, 0x39, 0x78, 0x20, 0x29, 0x7d, 0xd9, 0x82, 0xb0, 0x15, 0x81, 0x81, 0x81,
		0xc5, 0x81, 0x35, 0x83, 0x0d, 0x82, 0xf1, 0x84, 0xe7, 0x0d, 0x84, 0xe2, 0x84, 0xcd, 0x8b, 0xe9,
		0x0d, 0x84, 0xb0, 0xc9, 0x82, 0xa9, 0x7f, 0x7d, 0x85, 0x95, 0x7e, 0xcd, 0x87, 0xc9, 0x7c, 0x20,
		0x21, 0x7d, 0x21, 0x7d, 0xb0, 0x81, 0x7e, 0x15, 0x81, 0xd1, 0x7c, 0xc9, 0x81, 0x15, 0x7b, 0x11,
		0x82, 0xe3, 0xc9, 0x87, 0x29, 0x7b, 0x20, 0xd9, 0x82, 0xd1, 0x82, 0xb0, 0xcd, 0x81, 0xb1, 0x7d,
		0xe9, 0x82, 0x76, 0x3d, 0x83, 0x39, 0x78, 0xe7, 0xf5, 0x7b, 0xb0, 0xb9, 0x7f, 0xbd, 0x81, 0x0d,
		0x7f, 0x71, 0x83, 0xf5, 0x7d, 0xf5, 0x84, 0xe1,
	}

	loading, _ := widget.NewIcon(loadingIconVG)

	return &LoadingIcon{
		icon:  loading,
		anime: gween.NewSequence(gween.New(5, 0, .05, ease.Linear)),
	}
}

// Start starts the sequence animation.
func (l *LoadingIcon) Start() {
	l.mu.Lock()
	l.stop = false
	if l.anime != nil {
		l.anime.SetLoop(-1)
	}
	l.mu.Unlock()
}

// Stops the animation sequence from running.
func (l *LoadingIcon) Stop() {
	l.mu.Lock()
	l.stop = true
	if l.anime != nil {
		l.anime.SetLoop(0)
	}
	l.mu.Unlock()
}

func (l *LoadingIcon) Layout(gtx layout.Context) layout.Dimensions {
	l.mu.Lock()
	defer l.mu.Unlock()

	r := op.Record(gtx.Ops)
	dims := l.icon.Layout(gtx, utils.SecondaryColor)
	c := r.Stop()

	gtx.Constraints.Min = dims.Size

	dt := float32(gtx.Now.Second())
	val, _, _ := l.anime.Update(dt)

	if !l.stop {
		op.InvalidateOp{}.Add(gtx.Ops)
		defer rotate(gtx, val).Push(gtx.Ops).Pop()
	}

	c.Add(gtx.Ops)

	return dims
}

func rotate(gtx layout.Context, value float32) op.TransformOp {
	pt := gtx.Constraints.Min.Div(2)
	origin := f32.Pt(float32(pt.X), float32(pt.Y))
	trans := f32.Affine2D{}.Rotate(origin, -value*2*math.Pi)
	return op.Affine(trans)
}