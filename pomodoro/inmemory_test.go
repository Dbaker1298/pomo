package pomodoro_test

import (
	"testing"

	"github.com/Dbaker1298/pomo/pomodoro"
	"github.com/Dbaker1298/pomo/pomodoro/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}
