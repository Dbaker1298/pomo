package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Category constants
const (
	CategoryPomodoro = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak = "LongBreak"
)

// State constants
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Struct to represent a Pomodoro interval
type Interval struct {
	ID int64
	StartTime time.Time
	PlannedDuration time.Duration
	ActualDuration time.Duration
	Category string
	State int
}

// "Repository Patter" - Abstract the data source by defining the interface here
// This interface defines the methods Create, Update, ByID, Last, and Breaks
type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ByID(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int) ([]Interval, error) 
}

// Error Values that this may return
var (
	ErrNoIntervals = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted = errors.New("Interval is completed or cancelled")
	ErrInvalidState = errors.New("Invalid State")
	ErrInvalidID = errors.New("Invalid ID")
)

// Configuration required to instantiate an interval
type IntervalConfig struct {
	repo Repository
	PomodoroDuration time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration time.Duration
}