package hcservice

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

// checkTask is a unit of work sent to a scheduler worker.
type checkTask struct {
	target CheckTarget
	somark int
}

// Scheduler performs periodic health checks using a tick-and-sweep algorithm
// with a bounded worker pool.
type Scheduler struct {
	state          *State
	somarks        *SomarkAllocator
	checkers       map[string]Checker // hcType -> checker
	workChan       chan checkTask
	nextCheck      map[string]time.Time // "vipKey:realAddr" -> next scheduled time
	spreadInterval time.Duration
	tickInterval   time.Duration
	workerCount    int
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewScheduler creates a new Scheduler.
//
// Parameters:
//   - state: The shared health state store.
//   - somarks: The somark allocator for looking up somark values.
//   - checkers: A map of healthcheck type to Checker implementation.
//   - cfg: The scheduler configuration.
//
// Returns a new Scheduler instance.
func NewScheduler(state *State, somarks *SomarkAllocator, checkers map[string]Checker, cfg SchedulerConfig) *Scheduler {
	return &Scheduler{
		state:          state,
		somarks:        somarks,
		checkers:       checkers,
		workChan:       make(chan checkTask, cfg.WorkerCount*2),
		nextCheck:      make(map[string]time.Time),
		spreadInterval: time.Duration(cfg.SpreadIntervalMs) * time.Millisecond,
		tickInterval:   time.Duration(cfg.TickIntervalMs) * time.Millisecond,
		workerCount:    cfg.WorkerCount,
	}
}

// Start launches the scheduler's worker pool and tick loop.
func (s *Scheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel

	// Launch workers
	for i := 0; i < s.workerCount; i++ {
		s.wg.Add(1)
		go s.worker(ctx)
	}

	// Launch tick loop
	s.wg.Add(1)
	go s.tickLoop(ctx)
}

// Stop cancels the scheduler and waits for all workers to finish.
func (s *Scheduler) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
	s.wg.Wait()
}

// NotifyVIPRegistered sets up staggered initial check times for a newly registered VIP.
//
// Parameters:
//   - key: The VIP key.
//   - realAddrs: The real addresses registered under this VIP.
func (s *Scheduler) NotifyVIPRegistered(key VIPKey, realAddrs []string) {
	now := time.Now()
	n := len(realAddrs)
	for i, addr := range realAddrs {
		taskKey := checkTaskKey(key, addr)
		var offset time.Duration
		if n > 1 {
			offset = s.spreadInterval * time.Duration(i) / time.Duration(n)
		}
		s.nextCheck[taskKey] = now.Add(offset)
	}
}

// NotifyRealsAdded sets up staggered initial check times for newly added reals.
//
// Parameters:
//   - key: The VIP key.
//   - realAddrs: The new real addresses added.
func (s *Scheduler) NotifyRealsAdded(key VIPKey, realAddrs []string) {
	s.NotifyVIPRegistered(key, realAddrs)
}

// NotifyVIPDeregistered removes all scheduled checks for a deregistered VIP.
//
// Parameters:
//   - key: The VIP key.
//   - realAddrs: The real addresses that were under this VIP.
func (s *Scheduler) NotifyVIPDeregistered(key VIPKey, realAddrs []string) {
	for _, addr := range realAddrs {
		delete(s.nextCheck, checkTaskKey(key, addr))
	}
}

// NotifyRealsRemoved removes scheduled checks for removed reals.
//
// Parameters:
//   - key: The VIP key.
//   - realAddrs: The real addresses removed.
func (s *Scheduler) NotifyRealsRemoved(key VIPKey, realAddrs []string) {
	s.NotifyVIPDeregistered(key, realAddrs)
}

// NotifyVIPUpdated replaces scheduled checks when a VIP's reals are replaced.
//
// Parameters:
//   - key: The VIP key.
//   - oldReals: The old real addresses to clean up.
//   - newReals: The new real addresses to schedule.
func (s *Scheduler) NotifyVIPUpdated(key VIPKey, oldReals, newReals []string) {
	s.NotifyVIPDeregistered(key, oldReals)
	s.NotifyVIPRegistered(key, newReals)
}

// tickLoop runs the sweep at each tick interval.
func (s *Scheduler) tickLoop(ctx context.Context) {
	defer s.wg.Done()
	ticker := time.NewTicker(s.tickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.tick()
		}
	}
}

// tick performs a single sweep over all check targets.
func (s *Scheduler) tick() {
	now := time.Now()
	targets := s.state.GetAllCheckTargets()

	for _, t := range targets {
		taskKey := checkTaskKey(t.VIPKey, t.RealAddr)

		// If never seen, set initial check time with random spread
		nextTime, exists := s.nextCheck[taskKey]
		if !exists {
			offset := time.Duration(rand.Int63n(int64(s.spreadInterval)))
			s.nextCheck[taskKey] = now.Add(offset)
			continue
		}

		if now.Before(nextTime) {
			continue
		}

		somark, ok := s.somarks.GetSomark(t.RealAddr)
		if !ok {
			continue
		}

		task := checkTask{
			target: t,
			somark: int(somark),
		}

		// Non-blocking send; skip if workers are saturated
		select {
		case s.workChan <- task:
			s.nextCheck[taskKey] = now.Add(time.Duration(t.Config.IntervalMs) * time.Millisecond)
		default:
			// Channel full, will retry next tick
		}
	}
}

// worker processes check tasks from the work channel.
func (s *Scheduler) worker(ctx context.Context) {
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case task, ok := <-s.workChan:
			if !ok {
				return
			}
			s.executeCheck(ctx, task)
		}
	}
}

// executeCheck runs a single health check and updates state.
func (s *Scheduler) executeCheck(ctx context.Context, task checkTask) {
	checker, ok := s.checkers[task.target.Config.Type]
	if !ok {
		log.Printf("no checker for type %q", task.target.Config.Type)
		return
	}

	timeout := checkerTimeout(&task.target.Config)
	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := checker.Check(checkCtx, task.target.VIPAddr, task.target.CheckPort, task.somark, &task.target.Config)
	s.state.UpdateRealHealth(task.target.VIPKey, task.target.RealAddr, result.Success)
}

// checkTaskKey builds the map key for the nextCheck schedule.
func checkTaskKey(key VIPKey, realAddr string) string {
	return key.String() + ":" + realAddr
}
