package assets

import (
	"bytes"
	"embed"
	"fmt"
	"image"

	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	AccountBalanceIcon, _ = widget.NewIcon(icons.ActionAccountBalance)
	AccountBoxIcon, _     = widget.NewIcon(icons.ActionAccountBox)
	BackIcon, _           = widget.NewIcon(icons.NavigationArrowBack)
	CartIcon, _           = widget.NewIcon(icons.ActionAddShoppingCart)
	EditIcon, _           = widget.NewIcon(icons.ContentCreate)
	HeartIcon, _          = widget.NewIcon(icons.ActionFavorite)
	HomeIcon, _           = widget.NewIcon(icons.ActionHome)
	MenuIcon, _           = widget.NewIcon(icons.NavigationMenu)
	OtherIcon, _          = widget.NewIcon(icons.ActionHelp)
	PlusIcon, _           = widget.NewIcon(icons.ContentAdd)
	SettingsIcon, _       = widget.NewIcon(icons.ActionSettings)
	VisibilityIcon, _     = widget.NewIcon(icons.ActionVisibility)

	SplashImage, _ = getImage("splash_image.png")
)

//go:embed images/*
var images embed.FS

// GetImage returns the image read from the path provided as parameter.
func getImage(path string) (image.Image, error) {
	data, err := images.ReadFile(fmt.Sprintf("images/%s", path))
	if err != nil {
		return nil, err
	}

	newImg, _, err := image.Decode(bytes.NewBuffer(data))

	return newImg, err
}

//go:embed fonts/*
var font embed.FS

func GetFont(path string) ([]byte, error) {
	data, err := font.ReadFile(fmt.Sprintf("fonts/%s", path))
	if err != nil {
		return nil, err
	}

	return data, err
}
