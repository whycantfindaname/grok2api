package requeststatus

import (
	"sync"
	"time"
)

type State string

const (
	StateRunning   State = "running"
	StateCompleted State = "completed"
	StateFailed    State = "failed"
)

type Status struct {
	State      State
	StartedAt  time.Time
	FinishedAt *time.Time
	StatusCode int
}

type requestKey struct {
	clientKeyID uint64
	requestID   string
}

type entry struct {
	status Status
	active int
}

// Registry tracks request lifecycle in process memory and retains terminal states briefly for polling clients.
type Registry struct {
	mu        sync.Mutex
	retention time.Duration
	entries   map[requestKey]entry
	nextPurge time.Time
}

func NewRegistry(retention time.Duration) *Registry {
	return &Registry{retention: retention, entries: make(map[requestKey]entry)}
}

func (r *Registry) Start(clientKeyID uint64, requestID string, now time.Time) {
	if r == nil || clientKeyID == 0 || requestID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.retention > 0 && (r.nextPurge.IsZero() || !now.Before(r.nextPurge)) {
		r.purgeExpired(now)
		r.nextPurge = now.Add(r.retention)
	}
	key := requestKey{clientKeyID: clientKeyID, requestID: requestID}
	value := r.entries[key]
	if value.active == 0 {
		value.status = Status{State: StateRunning, StartedAt: now.UTC()}
	}
	value.active++
	r.entries[key] = value
}

func (r *Registry) Finish(clientKeyID uint64, requestID string, statusCode int, now time.Time) {
	if r == nil || clientKeyID == 0 || requestID == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := requestKey{clientKeyID: clientKeyID, requestID: requestID}
	value, ok := r.entries[key]
	if !ok || value.active == 0 {
		return
	}
	value.active--
	if value.active == 0 {
		finishedAt := now.UTC()
		value.status.FinishedAt = &finishedAt
		value.status.StatusCode = statusCode
		value.status.State = StateCompleted
		if statusCode >= 400 {
			value.status.State = StateFailed
		}
	}
	r.entries[key] = value
}

func (r *Registry) Get(clientKeyID uint64, requestID string, now time.Time) (Status, bool) {
	if r == nil || clientKeyID == 0 || requestID == "" {
		return Status{}, false
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	key := requestKey{clientKeyID: clientKeyID, requestID: requestID}
	value, ok := r.entries[key]
	if !ok {
		return Status{}, false
	}
	if r.expired(value.status, now) {
		delete(r.entries, key)
		return Status{}, false
	}
	return value.status, true
}

func (r *Registry) purgeExpired(now time.Time) {
	for key, value := range r.entries {
		if r.expired(value.status, now) {
			delete(r.entries, key)
		}
	}
}

func (r *Registry) expired(status Status, now time.Time) bool {
	return status.FinishedAt != nil && r.retention > 0 && now.Sub(*status.FinishedAt) > r.retention
}
