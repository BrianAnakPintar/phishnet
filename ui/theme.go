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
