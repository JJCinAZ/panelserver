package minicron

import (
	"sort"
	"sync"
	"time"
)

// Handle is an opaque identifier for a job in a schedule. It must be a monotonically increasing number with a Schedule.
type Handle uint64

type CallBack func(...interface{})

type Job struct {
	handle           Handle
	Interval         time.Duration
	lastRun, nextRun time.Time
	callback         CallBack
	cbParams         []interface{}
}

type Schedule struct {
	mutex      sync.Mutex
	basetime   time.Time
	jobs       []Job
	nextHandle uint64
}

func NewSchedule() Schedule {
	return Schedule{
		basetime: time.Now(),
	}
}

func (sched *Schedule) AddJob(duration time.Duration, initialRun bool, cb CallBack, cbParams ...interface{}) Handle {
	sched.mutex.Lock()
	j := Job{
		Interval: duration,
		callback: cb,
		cbParams: cbParams,
		lastRun:  sched.basetime,
	}
	now := time.Now()
	// A schedule always has a fixed basetime from which all job next execution times are derived.  This
	// ensures the natural order of jobs and eliminates run-time delays from affecting order.
	// We need to increment the base time to calculate the next run for this new job so we loop until
	// we get to a time in the future.
	for {
		next := j.lastRun.Add(duration)
		if next.After(now) {
			j.nextRun = next
			break
		}
		j.lastRun = next
	}
	if initialRun {
		j.nextRun = j.nextRun.Add(-duration)
		j.lastRun = j.lastRun.Add(-duration)
	}
	sched.nextHandle++
	j.handle = Handle(sched.nextHandle)
	sched.jobs = append(sched.jobs, j)
	sched.sortjobs()
	sched.mutex.Unlock()
	return j.handle
}

// Sort jobs by the next one which needs to run
// mutex MUST be locked when running this
func (sched *Schedule) sortjobs() {
	sort.Slice(sched.jobs, func(i int, j int) bool {
		if sched.jobs[i].nextRun.Equal(sched.jobs[j].nextRun) {
			return sched.jobs[i].handle < sched.jobs[j].handle
		}
		return sched.jobs[i].nextRun.Before(sched.jobs[j].nextRun)
	})
}

func (j *Job) calcNextRun() {
}

func (sched *Schedule) RemoveJob(handle Handle) {
	sched.mutex.Lock()
	for i := range sched.jobs {
		if sched.jobs[i].handle == handle {
			sched.jobs = append(sched.jobs[:i], sched.jobs[i+1:]...)
			break
		}
	}
	sched.mutex.Unlock()
}

func (sched *Schedule) ServiceNextJob() {
	if len(sched.jobs) == 0 {
		return
	}
	now := time.Now()
	sched.mutex.Lock()
	j := &sched.jobs[0]
	if !now.Before(j.nextRun) {
		sched.basetime = now
		j.lastRun = now
		j.nextRun = j.lastRun.Add(j.Interval)
		if j.callback != nil {
			j.callback(j.cbParams...)
		}
		sched.sortjobs()
	}
	sched.mutex.Unlock()
}
