package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const (
	AppVersion     = "1.0.0"
	AppAuthor      = "Charly Fau + AI"
	AppEmail       = "totemcaf@gmail.com"
	AppWebsite     = "https://github.com/totemcaf"
	AppLicense     = "Open Source"
	AppCopyright   = "© 2025"
	AppDescription = "Una aplicación simple y elegante de Temporizador Pomodoro para Linux construida con Go y Fyne."
	AppCredits     = "Gracias a Cursor IDE Agent"
)

func (p *PomodoroApp) showAbout() {
	if p.aboutWindow != nil {
		p.aboutWindow.Show()
		return
	}

	p.createAboutWindow()
	p.aboutWindow.Show()
}

func (p *PomodoroApp) createAboutWindow() {
	p.aboutWindow = p.app.NewWindow("Acerca de Pomodoro Timer")
	p.aboutWindow.Resize(fyne.NewSize(500, 400))
	p.aboutWindow.CenterOnScreen()
	p.aboutWindow.SetFixedSize(true)

	// Set close intercept to hide instead of close
	p.aboutWindow.SetCloseIntercept(func() {
		p.aboutWindow.Hide()
	})

	// Create app icon using the embedded SVG resource
	iconResource := GetAppIcon()
	iconImage := canvas.NewImageFromResource(iconResource)
	iconImage.FillMode = canvas.ImageFillContain
	iconImage.SetMinSize(fyne.NewSize(64, 64))
	iconWidget := iconImage

	// App name and version
	titleLabel := widget.NewLabel("Pomodoro Timer")
	titleLabel.TextStyle = fyne.TextStyle{Bold: true}
	titleLabel.Alignment = fyne.TextAlignCenter

	versionLabel := widget.NewLabel(fmt.Sprintf("Versión %s", AppVersion))
	versionLabel.Alignment = fyne.TextAlignCenter

	// Description
	descriptionLabel := widget.NewLabel(AppDescription)
	descriptionLabel.Wrapping = fyne.TextWrapWord
	descriptionLabel.Alignment = fyne.TextAlignCenter

	// Author information
	authorLabel := widget.NewLabel(fmt.Sprintf("Autor: %s", AppAuthor))
	authorLabel.Alignment = fyne.TextAlignCenter

	emailLabel := widget.NewLabel(fmt.Sprintf("Email: %s", AppEmail))
	emailLabel.Alignment = fyne.TextAlignCenter

	websiteLabel := widget.NewLabel(fmt.Sprintf("Sitio web: %s", AppWebsite))
	websiteLabel.Alignment = fyne.TextAlignCenter

	// Copyright and license
	copyrightLabel := widget.NewLabel(fmt.Sprintf("%s %s", AppCopyright, AppAuthor))
	copyrightLabel.Alignment = fyne.TextAlignCenter

	licenseLabel := widget.NewLabel(fmt.Sprintf("Licencia: %s", AppLicense))
	licenseLabel.Alignment = fyne.TextAlignCenter

	// Credits
	creditsLabel := widget.NewLabel(AppCredits)
	creditsLabel.Alignment = fyne.TextAlignCenter
	creditsLabel.TextStyle = fyne.TextStyle{Italic: true}

	// Features list
	featuresTitle := widget.NewLabel("Características:")
	featuresTitle.TextStyle = fyne.TextStyle{Bold: true}

	features := widget.NewLabel(
		"• Pantalla visual del temporizador con fondos de colores\n" +
			"• Tiempos configurables para trabajo y descansos\n" +
			"• Transiciones automáticas entre períodos\n" +
			"• Funcionalidad de pausa/reanudar\n" +
			"• Gestión adecuada de ventanas\n" +
			"• Interfaz en español")
	features.Wrapping = fyne.TextWrapWord

	// Close button
	closeButton := widget.NewButton("Cerrar", func() {
		p.aboutWindow.Hide()
	})

	// Create sections
	headerSection := container.NewVBox(
		container.NewCenter(iconWidget),
		titleLabel,
		versionLabel,
		widget.NewSeparator(),
	)

	infoSection := container.NewVBox(
		descriptionLabel,
		widget.NewSeparator(),
		authorLabel,
		emailLabel,
		websiteLabel,
		widget.NewSeparator(),
		copyrightLabel,
		licenseLabel,
		widget.NewSeparator(),
		creditsLabel,
	)

	featuresSection := container.NewVBox(
		featuresTitle,
		features,
	)

	// Main content with scrolling
	content := container.NewVBox(
		headerSection,
		infoSection,
		featuresSection,
		widget.NewSeparator(),
		container.NewCenter(closeButton),
	)

	// Add some padding and make it scrollable
	scrollContent := container.NewScroll(content)
	scrollContent.SetMinSize(fyne.NewSize(480, 380))

	paddedContent := container.NewPadded(scrollContent)

	p.aboutWindow.SetContent(paddedContent)
}
