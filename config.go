package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (p *PomodoroApp) showConfig() {
	if p.configWindow != nil {
		// Refresh values to current configuration before showing
		p.refreshConfigValues()
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
	p.workTimeEntry = widget.NewEntry()
	p.workTimeEntry.SetText(fmt.Sprintf("%.0f", p.workTime.Minutes()))
	p.workTimeEntry.Validator = func(text string) error {
		if _, err := time.ParseDuration(text + "m"); err != nil {
			return fmt.Errorf("Ingrese un número válido de minutos")
		}
		return nil
	}

	p.shortBreakEntry = widget.NewEntry()
	p.shortBreakEntry.SetText(fmt.Sprintf("%.0f", p.shortBreakTime.Minutes()))
	p.shortBreakEntry.Validator = p.workTimeEntry.Validator

	p.longBreakEntry = widget.NewEntry()
	p.longBreakEntry.SetText(fmt.Sprintf("%.0f", p.longBreakTime.Minutes()))
	p.longBreakEntry.Validator = p.workTimeEntry.Validator

	p.shortBreaksCountEntry = widget.NewEntry()
	p.shortBreaksCountEntry.SetText(fmt.Sprintf("%d", p.shortBreaksBeforeLong))
	p.shortBreaksCountEntry.Validator = func(text string) error {
		if _, err := strconv.Atoi(text); err != nil {
			return fmt.Errorf("Ingrese un número válido")
		}
		return nil
	}

	// Create form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Tiempo de trabajo (minutos):", Widget: p.workTimeEntry},
			{Text: "Tiempo de descanso corto (minutos):", Widget: p.shortBreakEntry},
			{Text: "Tiempo de descanso largo (minutos):", Widget: p.longBreakEntry},
			{Text: "Descansos cortos antes del largo:", Widget: p.shortBreaksCountEntry},
		},
		OnSubmit: func() {
			// Parse and validate all entries
			if workTime, err := time.ParseDuration(p.workTimeEntry.Text + "m"); err == nil {
				p.workTime = workTime
			}
			if shortBreak, err := time.ParseDuration(p.shortBreakEntry.Text + "m"); err == nil {
				p.shortBreakTime = shortBreak
			}
			if longBreak, err := time.ParseDuration(p.longBreakEntry.Text + "m"); err == nil {
				p.longBreakTime = longBreak
			}
			if count, err := strconv.Atoi(p.shortBreaksCountEntry.Text); err == nil {
				p.shortBreaksBeforeLong = count
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

func (p *PomodoroApp) refreshConfigValues() {
	// Update form entries with current configuration values
	if p.workTimeEntry != nil {
		p.workTimeEntry.SetText(fmt.Sprintf("%.0f", p.workTime.Minutes()))
	}
	if p.shortBreakEntry != nil {
		p.shortBreakEntry.SetText(fmt.Sprintf("%.0f", p.shortBreakTime.Minutes()))
	}
	if p.longBreakEntry != nil {
		p.longBreakEntry.SetText(fmt.Sprintf("%.0f", p.longBreakTime.Minutes()))
	}
	if p.shortBreaksCountEntry != nil {
		p.shortBreaksCountEntry.SetText(fmt.Sprintf("%d", p.shortBreaksBeforeLong))
	}
}
