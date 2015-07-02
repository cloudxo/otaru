package scheduler

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
)

type Result interface {
	Err() error
}

type Task interface {
	Run(ctx context.Context) Result
}

type State int32

const (
	JobScheduled = iota
	JobStarted   = iota
	JobFinished  = iota
	JobAborted   = iota
)

func (st State) String() string {
	switch st {
	case JobScheduled:
		return "JobScheduled"
	case JobStarted:
		return "JobStarted"
	case JobFinished:
		return "JobFinished"
	case JobAborted:
		return "JobAborted"
	default:
		return "UnknownJobState"
	}
}

var ZeroTime time.Time

type job struct {
	ID
	State
	// Issuer string

	CreatedAt   time.Time
	ScheduledAt time.Time
	StartedAt   time.Time
	FinishedAt  time.Time

	Task
	Result
	DoneCallback

	mu         sync.Mutex
	cancelfn   context.CancelFunc
	scheduledC chan struct{}
}

type JobView struct {
	ID    `json:"id"`
	State `json:"state,string"`
	// Issuer string

	CreatedAt   time.Time `json:"created_at"`
	ScheduledAt time.Time `json:"scheduled_at"`
	StartedAt   time.Time `json:"started_at"`
	FinishedAt  time.Time `json:"finishd_at"`

	Result `json:"result"`
}

type DoneCallback func(*JobView)

func (j *job) ViewWithLock() *JobView {
	return &JobView{
		ID:          j.ID,
		State:       j.State,
		CreatedAt:   j.CreatedAt,
		ScheduledAt: j.ScheduledAt,
		StartedAt:   j.StartedAt,
		FinishedAt:  j.FinishedAt,
		Result:      j.Result,
	}
}

func (j *job) View() *JobView {
	j.mu.Lock()
	defer j.mu.Unlock()
	return j.ViewWithLock()
}

type jobQuery struct {
	ID
	resultC chan *JobView
}

type abortReq struct {
	ID
	doneC chan struct{}
}

type Scheduler struct {
	numRunners int

	idGen

	scheduleC     chan *job
	queryC        chan *jobQuery
	abortC        chan *abortReq
	runC          chan *job
	joinScheduleC chan struct{}
	joinRunnerC   chan struct{}

	numWaitJobs int
}

const schedulerTickDuration = 300 * time.Millisecond

func NewScheduler() *Scheduler {
	s := &Scheduler{
		numRunners:    4, // FIXME
		idGen:         idGen{0},
		scheduleC:     make(chan *job, 1),
		queryC:        make(chan *jobQuery, 1),
		abortC:        make(chan *abortReq, 1),
		runC:          make(chan *job, 8),
		joinScheduleC: make(chan struct{}),
		joinRunnerC:   make(chan struct{}),
	}

	go s.schedulerMain()
	for i := 0; i < s.numRunners; i++ {
		go s.runnerMain()
	}

	return s
}

func (s *Scheduler) RunAt(task Task, at time.Time, cb DoneCallback) ID {
	id := s.idGen.genID()
	j := &job{
		ID:           id,
		CreatedAt:    time.Now(),
		ScheduledAt:  at,
		Task:         task,
		DoneCallback: cb,
		scheduledC:   make(chan struct{}),
	}

	s.scheduleC <- j
	<-j.scheduledC

	return id
}

func (s *Scheduler) RunImmediately(task Task, cb DoneCallback) ID {
	return s.RunAt(task, ZeroTime, cb)
}

func (s *Scheduler) Query(id ID) *JobView {
	q := &jobQuery{
		ID:      id,
		resultC: make(chan *JobView),
	}
	s.queryC <- q
	return <-q.resultC
}

func (s *Scheduler) Abort(id ID) {
	req := &abortReq{
		ID:    id,
		doneC: make(chan struct{}),
	}
	s.abortC <- req
	<-req.doneC
}

func (s *Scheduler) RunAllAndStop() {
	close(s.scheduleC)
	close(s.queryC)

	<-s.joinScheduleC
	for i := 0; i < s.numRunners; i++ {
		<-s.joinRunnerC
	}
}

func (s *Scheduler) schedulerMain() {
	tick := time.NewTicker(schedulerTickDuration) // FIXME: This should actually wait until next scheduled task instead of using fixed duration ticker.
	waitJobs := make([]*job, 0)
	jobs := make(map[ID]*job)

	defer func() {
		close(s.runC)
		s.joinScheduleC <- struct{}{}
	}()

	scheduleC := s.scheduleC
	for {
		if scheduleC == nil && s.numWaitJobs == 0 {
			return
		}

		select {
		case j, more := <-scheduleC:
			if !more {
				// stop polling on scheduleC
				scheduleC = nil
				continue
			}

			if _, ok := jobs[j.ID]; ok {
				log.Printf("job ID %v is already taken. received duplicate: %v", j.ID, j)
				if j.scheduledC != nil {
					close(j.scheduledC)
				}
				continue
			}
			jobs[j.ID] = j

			if j.ScheduledAt.Before(time.Now()) {
				s.runC <- j
			} else {
				waitJobs = append(waitJobs, j)
				s.numWaitJobs = len(waitJobs)
			}
			if j.scheduledC != nil {
				close(j.scheduledC)
			}

		case q := <-s.queryC:
			if q == nil {
				continue
			}
			id := q.ID
			j, ok := jobs[id]
			if !ok {
				q.resultC <- nil
			} else {
				q.resultC <- j.View()
			}

		case req := <-s.abortC:
			j, ok := jobs[req.ID]
			if !ok {
				log.Printf("Abort target job ID %d doesn't exist.", req.ID)
				req.doneC <- struct{}{}
				continue
			}
			j.mu.Lock()
			switch j.State {
			case JobScheduled:
				j.State = JobAborted
				if j.DoneCallback != nil {
					go j.DoneCallback(j.ViewWithLock())
				}
			case JobStarted:
				j.cancelfn()
			case JobFinished:
				// Job has already finished. Too late.
			case JobAborted:
				// Job is already aborted. Nothing to do.
			}
			j.mu.Unlock()
			req.doneC <- struct{}{}

		case <-tick.C:
			stillWaitJobs := make([]*job, 0, len(waitJobs))
			now := time.Now()
			for _, j := range waitJobs {
				if j.ScheduledAt.Before(now) {
					s.runC <- j
				} else {
					stillWaitJobs = append(stillWaitJobs, j)
				}
			}
			waitJobs = stillWaitJobs
			s.numWaitJobs = len(waitJobs)
		}
	}
}

func (s *Scheduler) runnerMain() {
	for j := range s.runC {
		j.mu.Lock()
		if j.State == JobAborted {
			j.mu.Unlock()
			continue
		}
		if j.State != JobScheduled {
			log.Printf("Skipping job not in scheduled state: %v", j)

			j.mu.Unlock()
			continue
		}

		task := j.Task
		ctx, cancelfn := context.WithCancel(context.Background())
		j.cancelfn = cancelfn
		j.StartedAt = time.Now()
		j.State = JobStarted

		j.mu.Unlock()
		result := task.Run(ctx)
		finishedAt := time.Now()
		j.mu.Lock()

		j.Result = result
		j.FinishedAt = finishedAt
		if ctx.Err() != nil {
			j.State = JobAborted
		} else {
			j.State = JobFinished
		}
		if j.DoneCallback != nil {
			go j.DoneCallback(j.ViewWithLock())
		}
		j.mu.Unlock()
	}
	s.joinRunnerC <- struct{}{}
}
