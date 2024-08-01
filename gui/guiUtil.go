package gui

import "fyne.io/fyne/v2"
import "github.com/kbinani/screenshot"

func WindowSize() fyne.Size {
	if screenshot.NumActiveDisplays() > 0 {
		bounds := screenshot.GetDisplayBounds(0)
		return fyne.NewSize(float32(bounds.Dx())*0.5, float32(bounds.Dy())*0.5)
	}
	return fyne.NewSize(700, 450)
}

func WindowWidth() float32 {
	canvasWidth := cubusWindow.Canvas().Size().Width
	if canvasWidth == 0 {
		return WindowSize().Width
	} else {
		return canvasWidth
	}
}

func WindowHeight() float32 {
	canvasHeight := cubusWindow.Canvas().Size().Height
	if canvasHeight == 0 {
		return WindowSize().Height
	} else {
		return canvasHeight
	}
}
