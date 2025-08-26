package main

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed pomodoro-timer.svg
var iconSVGData []byte

// GetAppIcon returns the embedded app icon as a Fyne resource
func GetAppIcon() fyne.Resource {
	return fyne.NewStaticResource("pomodoro-timer.svg", iconSVGData)
}
