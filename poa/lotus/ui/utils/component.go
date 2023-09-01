// Copyright (c) 2023 Migwi Ndung'u
// See LICENSE for details.

package utils

import (
	"fmt"
	"image/color"
)

const (
	AppName = "LOTUS"

	FEEDBACK_PAGE_ID = "FEEDBACK_PAGE"
	HOME_PAGE_ID     = "HOME_PAGE"
	ABOUT_PAGE_ID    = "ABOUT_PAGE"
	ACCOUNT_PAGE_ID  = "ACCOUNT_PAGE"
	SIGNUP_PAGE_ID   = "SIGNUP_PAGE"
	SPLASH_PAGE_ID   = "SPLASH_PAGE"
)

// AppVersion defines the current semantic version of the lotus application.
var AppVersion = &version{
	Major: 0,
	Minor: 0,
	Patch: 1,
}

var (
	// Application colour scheme.
	PrimaryColor   = color.NRGBA{R: 98, G: 0, B: 238, A: 255}
	DarkPriColor   = color.NRGBA{R: 55, G: 0, B: 179, A: 255}
	SecondaryColor = color.NRGBA{R: 3, G: 218, B: 198, A: 255}
	DarkSecColor   = color.NRGBA{R: 1, G: 135, B: 134, A: 255}
	ErrorColor     = color.NRGBA{R: 176, G: 0, B: 32, A: 255}
	SurfaceColor   = color.White
	BlackColor     = color.Black
	HighlightColor = color.NRGBA{R: 187, G: 134, B: 252}
)

type version struct {
	Major uint8
	Minor uint8
	Patch uint8
}

// String implements the stringer interface for struct version.
func (v *version) String() string {
	return fmt.Sprintf("%s v%d.%d.%d", AppName, v.Major, v.Minor, v.Patch)
}
