package scheduler

import (
	"fmt"

	"github.com/go-co-op/gocron/v2"
)

func Create(cexpr string) (gocron.Scheduler, error) {
	s, _ := gocron.NewScheduler()
	cdef := gocron.CronJob(cexpr, false)

	if _, err := s.NewJob(cdef, gocron.NewTask(func() {
		fmt.Println("here")
	})); err != nil {
		return nil, err
	}

	return s, nil
}
