package pomodoro_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Dbaker1298/pomo/pomodoro"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name   string
		input  [3]time.Duration
		expect pomodoro.IntervalConfig
	}{
		{
			name: "Default",
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   25 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "SingleInput",
			input: [3]time.Duration{
				20 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "MultiInput",
			input: [3]time.Duration{
				20 * time.Minute,
				10 * time.Minute,
				12 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 10 * time.Minute,
				LongBreakDuration:  12 * time.Minute,
			},
		},
	}

	// Execute tests for NewConfig
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var repo pomodoro.Repository
			config := pomodoro.NewConfig(
				repo,
				tc.input[0],
				tc.input[1],
				tc.input[2],
			)

			if config.PomodoroDuration != tc.expect.PomodoroDuration {
				t.Errorf("Expected Pomodoro Duration %q, got %q instead\n", tc.expect.PomodoroDuration, config.PomodoroDuration)
			}

			if config.ShortBreakDuration != tc.expect.ShortBreakDuration {
				t.Errorf("Expected Short Break Duration %q, got %q instead\n", tc.expect.ShortBreakDuration, config.ShortBreakDuration)
			}

			if config.LongBreakDuration != tc.expect.LongBreakDuration {
				t.Errorf("Expected Long Break Duration %q, got %q instead\n", tc.expect.LongBreakDuration, config.LongBreakDuration)
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	repo, cleanup := getRepo(t)
	defer cleanup()

	// For this test, need to start and complete each interval to allow GetInterval
	// to obtain the next category
	const duration = 1 * time.Millisecond
	config := pomodoro.NewConfig(repo, 3*duration, duration, 2*duration)

	// Start loop to execute the test 16 times
	for i := 0; i < 16; i++ {
		var (
			expCategory string
			expDuration time.Duration
		)

		switch {
		case i%2 != 0:
			expCategory = pomodoro.CategoryPomodoro
			expDuration = 3 * duration
		case i%8 == 0:
			expCategory = pomodoro.CategoryLongBreak
			expDuration = 2 * duration
		case i%2 == 0:
			expCategory = pomodoro.CategoryShortBreak
			expDuration = duration
		}

		testName := fmt.Sprintf("%s%d", expCategory, i)
		t.Run(testName, func(t *testing.T) {
			res, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Errorf("Expected no error, got %q instead\n", err)
			}

			noop := func(pomodoro.Interval) {}

			if err := res.Start(context.Background(), config, noop, noop, noop); err != nil {
				t.Fatal(err)
			}

			if res.Category != expCategory {
				t.Errorf("Expected category %q, got %q instead\n", expCategory, res.Category)
			}

			if res.PlannedDuration != expDuration {
				t.Errorf("Expected duration %q, got %q instead\n", expDuration, res.PlannedDuration)
			}

			if res.State != pomodoro.StateNotStarted {
				t.Errorf("Expected state = %q, got %q instead\n", pomodoro.StateNotStarted, res.State)
			}

			ui, err := repo.ByID(res.ID)
			if err != nil {
				t.Errorf("Expected no error, got %q instead\n", err)
			}

			if ui.State != pomodoro.StateDone {
				t.Errorf("Expected state = %q, got %q instead\n", pomodoro.StateDone, res.State)
			}
		})

	}
}
