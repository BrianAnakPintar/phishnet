package ui

import (
	"image"
	imgcolor "image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/briananakpintar/phishnet/syscalls"
)

// findRandomImageInDataImages searches common data/images locations (similar to PhishTank's locateDataFiles)
// and returns a random image file path (png/jpg/jpeg) found, or an empty string when none exist.
func findRandomImageInDataImages() string {
	candidates := []string{}
	if ex, err := os.Executable(); err == nil {
		exDir := filepath.Dir(ex)
		candidates = append(candidates, filepath.Join(exDir, "data", "images"))
	}
	candidates = append(candidates, "./data/images", "data/images", filepath.Join("..", "data", "images"))

	var imgs []string
	for _, d := range candidates {
		entries, err := os.ReadDir(d)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			name := strings.ToLower(e.Name())
			// Accept common static image formats (png, jpg, jpeg)
			if strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") {
				imgs = append(imgs, filepath.Join(d, e.Name()))
			}
		}
	}
	if len(imgs) == 0 {
		return ""
	}
	// ensure randomness
	rand.Seed(time.Now().UnixNano())
	return imgs[rand.Intn(len(imgs))]
}

// makeRoundedCorners loads an image from disk and returns a copy with rounded corners applied.
// The returned image keeps the original dimensions; Fyne will scale it according to the widget size.
func makeRoundedCorners(path string, radius int) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	src, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	b := src.Bounds()
	w, h := b.Dx(), b.Dy()
	if radius <= 0 {
		radius = int(math.Min(float64(w), float64(h)) / 8)
	}

	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	// copy source into dst
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)

	// apply rounded mask by adjusting alpha in corner pixels
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// determine which corner we're in (if any)
			var cx, cy int
			if x < radius && y < radius {
				// top-left
				cx, cy = radius-1, radius-1
			} else if x >= w-radius && y < radius {
				// top-right
				cx, cy = w-radius, radius-1
			} else if x < radius && y >= h-radius {
				// bottom-left
				cx, cy = radius-1, h-radius
			} else if x >= w-radius && y >= h-radius {
				// bottom-right
				cx, cy = w-radius, h-radius
			} else {
				// not in a corner area
				continue
			}
			dx := float64(x - cx)
			dy := float64(y - cy)
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > float64(radius) {
				// outside the rounded corner: make transparent
				idx := dst.PixOffset(x, y)
				// set alpha to 0
				dst.Pix[idx+3] = 0
			}
		}
	}

	return dst, nil
}

func Popup(logs string, url string) {
	a := app.New()
	// Try to load icon.png from common locations and set it as the app/window icon.
	var iconRes fyne.Resource
	candidates := []string{}
	// prefer executable-dir first
	if ex, err := os.Executable(); err == nil {
		exDir := filepath.Dir(ex)
		candidates = append(candidates, filepath.Join(exDir, "icon.png"))
	}
	// also check current working directory (useful when running from project root)
	if wd, err := os.Getwd(); err == nil {
		candidates = append(candidates, filepath.Join(wd, "icon.png"))
	}
	// fallback relative paths
	candidates = append(candidates, "./icon.png", "icon.png", filepath.Join("..", "icon.png"))
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if res, err := fyne.LoadResourceFromPath(p); err == nil {
				a.SetIcon(res)
				iconRes = res
				break
			}
		}
	}

	a.Settings().SetTheme(NordTheme{})
	w := a.NewWindow("PhishNet Warning")
	if iconRes != nil {
		w.SetIcon(iconRes)
	}
	w.Resize(fyne.NewSize(360, 140))
	w.SetFixedSize(true)

	// Left side: rounded image or emoji
	var leftContainer *fyne.Container
	imgPath := findRandomImageInDataImages()
	if imgPath != "" {
		processed, err := makeRoundedCorners(imgPath, 0)
		if err == nil {
			img := canvas.NewImageFromImage(processed)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(56, 56))
			centered := container.NewCenter(img)
			leftContainer = container.NewVBox(layout.NewSpacer(), centered, layout.NewSpacer())
		} else {
			// fallback to loading directly
			img := canvas.NewImageFromFile(imgPath)
			img.FillMode = canvas.ImageFillContain
			img.SetMinSize(fyne.NewSize(56, 56))
			centered := container.NewCenter(img)
			leftContainer = container.NewVBox(layout.NewSpacer(), centered, layout.NewSpacer())
		}
	} else {
		emoji := canvas.NewText("PHISH!", imgcolor.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF})
		emoji.TextSize = 56
		centered := container.NewCenter(emoji)
		leftContainer = container.NewVBox(layout.NewSpacer(), centered, layout.NewSpacer())
	}

	// Placeholder text (non-scrollable) shown in main popup
	placeholder := "Whoops, looks like the site you're trying to visit is unsafe"
	placeholderLabel := widget.NewLabel(placeholder)
	placeholderLabel.Wrapping = fyne.TextWrapWord

	// URL label shown above the placeholder
	urlLabel := widget.NewLabelWithStyle("The URL: "+url, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	cancel := widget.NewButton("Cancel", func() { w.Close() })

	// Buttons: View Logs opens a new window with the read-only scrollable logs
	viewBtn := widget.NewButton("View Logs", func() {
		logsWin := a.NewWindow("PhishNet Logs")
		logsWin.Resize(fyne.NewSize(560, 300))
		// Create read-only multiline entry with scrolling for logs
		logEntry := widget.NewMultiLineEntry()
		logEntry.SetText(logs)
		logEntry.Wrapping = fyne.TextWrapOff
		logEntry.Disable()
		scroll := container.NewScroll(logEntry)
		scroll.SetMinSize(fyne.NewSize(540, 260))
		logsWin.SetContent(container.NewBorder(nil, nil, nil, nil, scroll))
		logsWin.Show()
	})
	// make the button less prominent so it reads more like a link
	viewBtn.Importance = widget.LowImportance

	continueBtn := widget.NewButton("Open Anyway", func() {
		syscalls.OpenChrome(url)
		w.Close()
	})

	buttons := container.NewHBox(cancel, layout.NewSpacer(), viewBtn, layout.NewSpacer(), continueBtn)

	// Layout: left image, placeholder text on right, buttons bottom, URL on top
	// stack URL and placeholder vertically and add flexible spacers to provide padding
	center := container.NewVBox(urlLabel, layout.NewSpacer(), placeholderLabel, layout.NewSpacer())
	content := container.NewBorder(nil, buttons, leftContainer, nil, center)

	// add padding around the whole content
	padded := container.NewPadded(content)

	w.SetContent(padded)
	w.ShowAndRun()
}
