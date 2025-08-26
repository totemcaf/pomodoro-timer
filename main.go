package main

import (
	"fmt"
	"time"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Custom large time label using canvas.Text
type LargeTimeLabel struct {
	widget.BaseWidget
	text *canvas.Text
}

func NewLargeTimeLabel(timeText string) *LargeTimeLabel {
	label := &LargeTimeLabel{}
	label.text = canvas.NewText(timeText, color.RGBA{255, 255, 255, 255})
	label.text.TextSize = 72 // Large font size
	label.text.TextStyle = fyne.TextStyle{Bold: true, Monospace: true}
	label.text.Alignment = fyne.TextAlignCenter

	label.ExtendBaseWidget(label)
	return label
}

func (l *LargeTimeLabel) SetText(text string) {
	l.text.Text = text
	l.text.Refresh()
}

func (l *LargeTimeLabel) CreateRenderer() fyne.WidgetRenderer {
	return &largeTimeRenderer{text: l.text}
}

type largeTimeRenderer struct {
	text *canvas.Text
}

func (r *largeTimeRenderer) Layout(size fyne.Size) {
	r.text.Resize(size)
}

func (r *largeTimeRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 80)
}

func (r *largeTimeRenderer) Refresh() {
	// Update text color based on current theme
	r.text.Color = theme.ForegroundColor()
	r.text.Refresh()
}

func (r *largeTimeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.text}
}

func (r *largeTimeRenderer) Destroy() {
}

type PomodoroApp struct {
	app           fyne.App
	mainWindow    fyne.Window
	configWindow  fyne.Window
	timeLabel     *LargeTimeLabel
	startWorkBtn  *widget.Button
	suspendBtn    *widget.Button
	startBreakBtn *widget.Button
	configBtn     *widget.Button

	// Timer state
	isRunning     bool
	isPaused      bool
	timeRemaining time.Duration
	isWorkTime    bool

	// Configuration
	workTime              time.Duration
	shortBreakTime        time.Duration
	longBreakTime         time.Duration
	shortBreaksBeforeLong int
	currentBreakCount     int

	// Config form entries
	workTimeEntry         *widget.Entry
	shortBreakEntry       *widget.Entry
	longBreakEntry        *widget.Entry
	shortBreaksCountEntry *widget.Entry
}

func NewPomodoroApp() *PomodoroApp {
	myApp := app.New()

	return &PomodoroApp{
		app:                   myApp,
		workTime:              20 * time.Minute, // Default 20 minutes
		shortBreakTime:        5 * time.Minute,  // Default 5 minutes
		longBreakTime:         15 * time.Minute, // Default 15 minutes
		shortBreaksBeforeLong: 4,                // Default 4 short breaks before long
		timeRemaining:         20 * time.Minute,
		isWorkTime:            true,
	}
}

func (p *PomodoroApp) createMainWindow() {
	p.mainWindow = p.app.NewWindow("Pomodoro Timer")
	p.mainWindow.Resize(fyne.NewSize(400, 280))
	p.mainWindow.CenterOnScreen()
	p.mainWindow.SetMaster() // Set as master window

	// Close config window when main window closes
	p.mainWindow.SetCloseIntercept(func() {
		if p.configWindow != nil {
			p.configWindow.Close()
		}
		p.mainWindow.Close()
	})

	// Create time display label with large font
	p.timeLabel = NewLargeTimeLabel(p.formatTime(p.timeRemaining))

	// Create buttons
	p.startWorkBtn = widget.NewButton("Iniciar tiempo de trabajo", p.startWorkTime)
	p.suspendBtn = widget.NewButton("Suspender", p.suspend)
	p.startBreakBtn = widget.NewButton("Iniciar tiempo de descanso", p.startBreakTime)
	p.configBtn = widget.NewButton("Configuración", p.showConfig)

	// Initially suspend button is disabled
	p.suspendBtn.Disable()

	// Create layout
	content := container.NewVBox(
		widget.NewCard("", "", container.NewCenter(p.timeLabel)),
		container.NewGridWithColumns(2,
			p.startWorkBtn,
			p.suspendBtn,
		),
		container.NewGridWithColumns(2,
			p.startBreakBtn,
			p.configBtn,
		),
	)

	p.mainWindow.SetContent(content)
}

func (p *PomodoroApp) formatTime(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func (p *PomodoroApp) startWorkTime() {
	p.timeRemaining = p.workTime
	p.isWorkTime = true
	p.isRunning = true
	p.isPaused = false

	p.startWorkBtn.Disable()
	p.startBreakBtn.Disable()
	p.suspendBtn.Enable()
	p.suspendBtn.SetText("Suspender")

	p.updateTimeDisplay()
	go p.runTimer()
}

func (p *PomodoroApp) startBreakTime() {
	// Determine if it's time for a long break
	if p.currentBreakCount >= p.shortBreaksBeforeLong {
		p.timeRemaining = p.longBreakTime
		p.currentBreakCount = 0
	} else {
		p.timeRemaining = p.shortBreakTime
		p.currentBreakCount++
	}

	p.isWorkTime = false
	p.isRunning = true
	p.isPaused = false

	p.startWorkBtn.Disable()
	p.startBreakBtn.Disable()
	p.suspendBtn.Enable()
	p.suspendBtn.SetText("Suspender")

	p.updateTimeDisplay()
	go p.runTimer()
}

func (p *PomodoroApp) suspend() {
	if p.isRunning {
		// Suspend the timer
		p.isPaused = true
		p.isRunning = false
		p.suspendBtn.SetText("Continuar")
		p.startWorkBtn.Enable()
		p.startBreakBtn.Enable()
	} else if p.isPaused {
		// Continue the timer
		p.isPaused = false
		p.isRunning = true
		p.suspendBtn.SetText("Suspender")
		p.startWorkBtn.Disable()
		p.startBreakBtn.Disable()
		go p.runTimer()
	}
}

func (p *PomodoroApp) updateTimeDisplay() {
	fyne.Do(func() {
		p.timeLabel.SetText(p.formatTime(p.timeRemaining))
	})
}

func (p *PomodoroApp) runTimer() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for p.isRunning {
		select {
		case <-ticker.C:
			p.timeRemaining -= time.Second
			p.updateTimeDisplay()

			if p.timeRemaining <= 0 {
				p.timerFinished()
				return
			}
		}
	}
}

func (p *PomodoroApp) timerFinished() {
	p.isRunning = false
	p.isPaused = false

	fyne.Do(func() {
		// Bring window to front and show message
		p.mainWindow.RequestFocus()
		p.mainWindow.Show()

		if p.isWorkTime {
			// Work time finished, start break automatically
			popup := widget.NewModalPopUp(
				widget.NewLabel("Iniciar un ciclo de descanso"),
				p.mainWindow.Canvas(),
			)
			popup.Show()

			// Auto start break after showing popup
			go func() {
				time.Sleep(2 * time.Second)
				fyne.Do(func() {
					popup.Hide()
					p.startBreakTime()
				})
			}()
		} else {
			// Break time finished, start work automatically
			popup := widget.NewModalPopUp(
				widget.NewLabel("Continuar trabajando"),
				p.mainWindow.Canvas(),
			)
			popup.Show()

			// Auto start work after showing popup
			go func() {
				time.Sleep(2 * time.Second)
				fyne.Do(func() {
					popup.Hide()
					p.startWorkTime()
				})
			}()
		}
	})
}

func (p *PomodoroApp) Run() {
	p.createMainWindow()
	p.mainWindow.ShowAndRun()
}

func main() {
	app := NewPomodoroApp()
	app.Run()
}
