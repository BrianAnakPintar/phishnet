package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CartoonTheme defines a playful cartoon-like theme
type CartoonTheme struct{}

// Colors: bright and bold
func (c CartoonTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.RGBA{R: 255, G: 246, B: 207, A: 255} // light yellow
	case theme.ColorNameButton:
		return color.RGBA{R: 255, G: 183, B: 3, A: 255} // cartoon orange
	case theme.ColorNameDisabled:
		return color.RGBA{R: 180, G: 180, B: 180, A: 255}
	case theme.ColorNameInputBackground:
		return color.RGBA{R: 255, G: 236, B: 179, A: 255}
	case theme.ColorNameForeground:
		return color.Black
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Use bigger, cartoonish fonts
func (c CartoonTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style) // You could swap in a comic font here if available
}

// Larger padding, rounder look
func (c CartoonTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 12
	case theme.SizeNamePadding:
		return 8
	case theme.SizeNameInlineIcon:
		return 14
	default:
		return theme.DefaultTheme().Size(name)
	}
}

func (c CartoonTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Nord color palette
var (
	nord0  = color.NRGBA{46, 52, 64, 255} // dark background
	nord1  = color.NRGBA{59, 66, 82, 255}
	nord2  = color.NRGBA{67, 76, 94, 255}
	nord3  = color.NRGBA{76, 86, 106, 255}
	nord4  = color.NRGBA{216, 222, 233, 255} // light text
	nord5  = color.NRGBA{229, 233, 240, 255}
	nord6  = color.NRGBA{236, 239, 244, 255}
	nord7  = color.NRGBA{143, 188, 187, 255} // cyan
	nord8  = color.NRGBA{136, 192, 208, 255}
	nord9  = color.NRGBA{129, 161, 193, 255} // blue
	nord10 = color.NRGBA{94, 129, 172, 255}
	nord11 = color.NRGBA{191, 97, 106, 255}  // red
	nord12 = color.NRGBA{208, 135, 112, 255} // orange
	nord13 = color.NRGBA{235, 203, 139, 255} // yellow
	nord14 = color.NRGBA{163, 190, 140, 255} // green
	nord15 = color.NRGBA{180, 142, 173, 255} // purple
)

// NordTheme implements fyne.Theme
type NordTheme struct{}

func (NordTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return nord0
	case theme.ColorNameButton:
		return nord9
	case theme.ColorNameDisabled:
		return nord3
	case theme.ColorNameDisabledButton:
		return nord2
	case theme.ColorNameHover:
		return nord2
	case theme.ColorNamePlaceHolder:
		return nord3
	case theme.ColorNamePrimary:
		return nord8
	case theme.ColorNameScrollBar:
		return nord2
	case theme.ColorNameShadow:
		return nord1
	case theme.ColorNameForeground:
		if variant == theme.VariantDark {
			return nord6
		}
		return nord0
	default:
		return nord4
	}
}

func (NordTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (NordTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (NordTheme) Size(name fyne.ThemeSizeName) float32 {
	switch name {
	case theme.SizeNameText:
		return 14
	case theme.SizeNameHeadingText:
		return 20
	case theme.SizeNameInlineIcon:
		return 18
	default:
		return theme.DefaultTheme().Size(name)
	}
}
