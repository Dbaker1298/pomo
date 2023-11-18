package cmd

import (
	"github.com/Dbaker1298/pomo/pomodoro"
	"github.com/Dbaker1298/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	return repository.NewInMemoryRepo(), nil
}
