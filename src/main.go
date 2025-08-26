package main

import (
	"fmt"
	"time"

	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// Custom large time label using canvas.Text with background
type LargeTimeLabel struct {
	widget.BaseWidget
	text       *canvas.Text
	background *canvas.Rectangle
}

func NewLargeTimeLabel(timeText string) *LargeTimeLabel {
	label := &LargeTimeLabel{}

	// Create background rectangle
	label.background = canvas.NewRectangle(color.RGBA{0, 0, 0, 255}) // Start with black background

	// Create text
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

func (l *LargeTimeLabel) SetBackgroundColor(bgColor color.Color) {
	l.background.FillColor = bgColor
	l.background.Refresh()
}

func (l *LargeTimeLabel) CreateRenderer() fyne.WidgetRenderer {
	return &largeTimeRenderer{text: l.text, background: l.background}
}

type largeTimeRenderer struct {
	text       *canvas.Text
	background *canvas.Rectangle
}

func (r *largeTimeRenderer) Layout(size fyne.Size) {
	r.background.Resize(size)

	// Center the text within the background
	textSize := r.text.MinSize()
	x := (size.Width - textSize.Width) / 2
	y := (size.Height - textSize.Height) / 2

	r.text.Move(fyne.NewPos(x, y))
	r.text.Resize(textSize)
}

func (r *largeTimeRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 80) // Wider minimum size to fill window width
}

func (r *largeTimeRenderer) Refresh() {
	// Keep text white for better contrast on colored backgrounds
	r.text.Color = color.RGBA{255, 255, 255, 255}
	r.text.Refresh()
	r.background.Refresh()
}

func (r *largeTimeRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.background, r.text}
}

func (r *largeTimeRenderer) Destroy() {
}

type PomodoroApp struct {
	app           fyne.App
	mainWindow    fyne.Window
	configWindow  fyne.Window
	aboutWindow   fyne.Window
	timeLabel     *LargeTimeLabel
	startWorkBtn  *widget.Button
	suspendBtn    *widget.Button
	startBreakBtn *widget.Button
	configBtn     *widget.Button
	aboutBtn      *widget.Button

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

	// Set application icon
	myApp.SetIcon(GetAppIcon())

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
	p.mainWindow.Resize(fyne.NewSize(400, 140))
	p.mainWindow.CenterOnScreen()
	p.mainWindow.SetMaster() // Set as master window

	// Close config and about windows when main window closes
	p.mainWindow.SetCloseIntercept(func() {
		if p.configWindow != nil {
			p.configWindow.Close()
		}
		if p.aboutWindow != nil {
			p.aboutWindow.Close()
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
	p.aboutBtn = widget.NewButton("Acerca de", p.showAbout)

	// Initially suspend button is disabled
	p.suspendBtn.Disable()

	// Create layout
	buttonContainer := container.NewVBox(
		container.NewGridWithColumns(2,
			p.startWorkBtn,
			p.aboutBtn,
			p.startBreakBtn,
			p.configBtn,
			p.suspendBtn,
		),
	)

	// Use border container to make time label fill the width
	content := container.NewBorder(
		p.timeLabel,     // Top - time label fills width
		buttonContainer, // Bottom - buttons
		nil,             // Left
		nil,             // Right
		nil,             // Center
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

	// Set background to light green for work time
	p.timeLabel.SetBackgroundColor(color.RGBA{144, 238, 144, 255}) // Light green

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

	// Set background to turquoise for break time
	p.timeLabel.SetBackgroundColor(color.RGBA{64, 224, 208, 255}) // Turquoise

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

		// Set background to black when suspended/not running
		p.timeLabel.SetBackgroundColor(color.RGBA{0, 0, 0, 255}) // Black
	} else if p.isPaused {
		// Continue the timer
		p.isPaused = false
		p.isRunning = true
		p.suspendBtn.SetText("Suspender")
		p.startWorkBtn.Disable()
		p.startBreakBtn.Disable()

		// Restore appropriate color based on timer type
		if p.isWorkTime {
			p.timeLabel.SetBackgroundColor(color.RGBA{144, 238, 144, 255}) // Light green for work
		} else {
			p.timeLabel.SetBackgroundColor(color.RGBA{64, 224, 208, 255}) // Turquoise for break
		}

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
			if p.isPaused {
				continue
			}
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
