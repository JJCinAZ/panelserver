package minicron

import (
	"testing"
	"time"
)

func TestSchedule_sortjobs(t *testing.T) {
	s := NewSchedule()
	s.AddJob(time.Minute, false, nil)
	s.AddJob(time.Second*5, false, nil)
	s.AddJob(time.Minute, false, nil)
	s.AddJob(time.Minute, true, nil)
	s.AddJob(time.Second*5, false, nil)
	t.Run("sort1", func(t *testing.T) {
		s.sortjobs()
	})
	correctOrder := []Handle{4, 2, 5, 1, 3}
	for i := range s.jobs {
		if s.jobs[i].handle != correctOrder[i] {
			t.Errorf("sortjobs failed")
		}
	}
}
