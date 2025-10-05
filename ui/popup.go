package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/briananakpintar/phishnet/syscalls"
)

func Popup(logs string, url string) {
	a := app.New()
	a.Settings().SetTheme(CartoonTheme{})
	w := a.NewWindow("PhishNet Alert")
	w.Resize(fyne.NewSize(320, 125))
	w.SetFixedSize(true)

	// Cartoonish left side: colored circle with emoji
	circle := canvas.NewCircle(&color.NRGBA{R: 0x2A, G: 0xB8, B: 0xD9, A: 0xFF})
	emoji := canvas.NewText("🐳", color.White)
	emoji.TextSize = 56
	stack := container.NewStack(circle, emoji)
	centered := container.NewCenter(stack)
	circleContainer := container.NewVBox(layout.NewSpacer(), centered, layout.NewSpacer())

	// Scrollable log area: use a disabled MultiLineEntry so it can scroll both axes
	logEntry := widget.NewMultiLineEntry()
	logEntry.SetText(logs)
	// Allow horizontal scrolling by disabling wrapping
	logEntry.Wrapping = fyne.TextWrapOff
	scroll := container.NewScroll(logEntry)
	scroll.SetMinSize(fyne.NewSize(200, 80))

	// Buttons
	cancel := widget.NewButton("Cancel", func() {
		w.Close()
	})
	continueBtn := widget.NewButton("Continue Anyway", func() {
		syscalls.OpenChrome(url)
		w.Close()
	})
	buttons := container.NewHBox(layout.NewSpacer(), cancel, continueBtn)

	// Layout: cartoon left, logs right, buttons bottom
	content := container.NewBorder(nil, buttons, circleContainer, nil, scroll)

	w.SetContent(content)
	w.ShowAndRun()
}
