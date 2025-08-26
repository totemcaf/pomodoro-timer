package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (p *PomodoroApp) showConfig() {
	if p.configWindow != nil {
		p.configWindow.Show()
		return
	}

	p.createConfigWindow()
	p.configWindow.Show()
}

func (p *PomodoroApp) createConfigWindow() {
	p.configWindow = p.app.NewWindow("Configuración")
	p.configWindow.Resize(fyne.NewSize(400, 300))
	p.configWindow.CenterOnScreen()
	
	// Set close intercept to properly handle window closing
	p.configWindow.SetCloseIntercept(func() {
		p.configWindow.Hide()
	})

	// Create form entries
	workTimeEntry := widget.NewEntry()
	workTimeEntry.SetText(fmt.Sprintf("%.0f", p.workTime.Minutes()))
	workTimeEntry.Validator = func(text string) error {
		if _, err := time.ParseDuration(text + "m"); err != nil {
			return fmt.Errorf("Ingrese un número válido de minutos")
		}
		return nil
	}

	shortBreakEntry := widget.NewEntry()
	shortBreakEntry.SetText(fmt.Sprintf("%.0f", p.shortBreakTime.Minutes()))
	shortBreakEntry.Validator = workTimeEntry.Validator

	longBreakEntry := widget.NewEntry()
	longBreakEntry.SetText(fmt.Sprintf("%.0f", p.longBreakTime.Minutes()))
	longBreakEntry.Validator = workTimeEntry.Validator

	shortBreaksCountEntry := widget.NewEntry()
	shortBreaksCountEntry.SetText(fmt.Sprintf("%d", p.shortBreaksBeforeLong))
	shortBreaksCountEntry.Validator = func(text string) error {
		if _, err := fmt.Sscanf(text, "%d"); err != nil {
			return fmt.Errorf("Ingrese un número válido")
		}
		return nil
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Tiempo de trabajo (minutos):", Widget: workTimeEntry},
			{Text: "Tiempo de descanso corto (minutos):", Widget: shortBreakEntry},
			{Text: "Tiempo de descanso largo (minutos):", Widget: longBreakEntry},
			{Text: "Descansos cortos antes del largo:", Widget: shortBreaksCountEntry},
		},
		OnSubmit: func() {
			// Parse and validate all entries
			if workTime, err := time.ParseDuration(workTimeEntry.Text + "m"); err == nil {
				p.workTime = workTime
			}
			if shortBreak, err := time.ParseDuration(shortBreakEntry.Text + "m"); err == nil {
				p.shortBreakTime = shortBreak
			}
			if longBreak, err := time.ParseDuration(longBreakEntry.Text + "m"); err == nil {
				p.longBreakTime = longBreak
			}
			if count, err := fmt.Sscanf(shortBreaksCountEntry.Text, "%d", &p.shortBreaksBeforeLong); err == nil {
				_ = count // Suppress unused variable warning
			}

			// Update current timer if not running
			if !p.isRunning {
				if p.isWorkTime {
					p.timeRemaining = p.workTime
				} else {
					if p.currentBreakCount >= p.shortBreaksBeforeLong {
						p.timeRemaining = p.longBreakTime
					} else {
						p.timeRemaining = p.shortBreakTime
					}
				}
				p.updateTimeDisplay()
			}

			p.configWindow.Hide()
		},
		OnCancel: func() {
			p.configWindow.Hide()
		},
	}

	// Set submit and cancel button text in Spanish
	form.SubmitText = "Guardar"
	form.CancelText = "Cancelar"

	content := container.NewVBox(
		widget.NewLabel("Configuración del Pomodoro"),
		form,
	)

	p.configWindow.SetContent(content)
}
