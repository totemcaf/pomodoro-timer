package main

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	app := NewPomodoroApp()

	tests := []struct {
		duration time.Duration
		expected string
	}{
		{25 * time.Minute, "25:00"},
		{5 * time.Minute, "05:00"},
		{1*time.Minute + 30*time.Second, "01:30"},
		{45 * time.Second, "00:45"},
		{0, "00:00"},
	}

	for _, test := range tests {
		result := app.formatTime(test.duration)
		if result != test.expected {
			t.Errorf("formatTime(%v) = %s; expected %s", test.duration, result, test.expected)
		}
	}
}

func TestPomodoroAppInitialization(t *testing.T) {
	app := NewPomodoroApp()

	// Test default values
	if app.workTime != 20*time.Minute {
		t.Errorf("Expected default work time to be 20 minutes, got %v", app.workTime)
	}

	if app.shortBreakTime != 5*time.Minute {
		t.Errorf("Expected default short break time to be 5 minutes, got %v", app.shortBreakTime)
	}

	if app.longBreakTime != 15*time.Minute {
		t.Errorf("Expected default long break time to be 15 minutes, got %v", app.longBreakTime)
	}

	if app.shortBreaksBeforeLong != 4 {
		t.Errorf("Expected 4 short breaks before long, got %d", app.shortBreaksBeforeLong)
	}

	// Test initial state
	if app.isRunning {
		t.Error("Expected timer to not be running initially")
	}

	if app.isPaused {
		t.Error("Expected timer to not be paused initially")
	}

	if !app.isWorkTime {
		t.Error("Expected initial state to be work time")
	}
}

func TestAboutWindowConstants(t *testing.T) {
	// Test that about window constants are properly set
	if AppVersion == "" {
		t.Error("AppVersion should not be empty")
	}
	
	if AppAuthor == "" {
		t.Error("AppAuthor should not be empty")
	}
	
	if AppEmail == "" {
		t.Error("AppEmail should not be empty")
	}
	
	if AppWebsite == "" {
		t.Error("AppWebsite should not be empty")
	}
	
	if AppLicense == "" {
		t.Error("AppLicense should not be empty")
	}
	
	// Test specific values
	expectedVersion := "1.0.0"
	if AppVersion != expectedVersion {
		t.Errorf("Expected version %s, got %s", expectedVersion, AppVersion)
	}
	
	expectedAuthor := "Charly Fau + AI"
	if AppAuthor != expectedAuthor {
		t.Errorf("Expected author %s, got %s", expectedAuthor, AppAuthor)
	}
}

func BenchmarkFormatTime(b *testing.B) {
	app := NewPomodoroApp()
	duration := 25*time.Minute + 30*time.Second
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		app.formatTime(duration)
	}
}
