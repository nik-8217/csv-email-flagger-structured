package jobs

import (
    "sync"
    "time"
)

type JobStatus string

const (
    StatusQueued     JobStatus = "QUEUED"
    StatusInProgress JobStatus = "IN_PROGRESS"
    StatusDone       JobStatus = "DONE"
    StatusFailed     JobStatus = "FAILED"
)

type Job struct {
    ID        string    `json:"id"`
    Status    JobStatus `json:"status"`
    InputPath string    `json:"-"`
    Output    string    `json:"-"`
    Error     string    `json:"error,omitempty"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    Mode      string    `json:"mode"`
}

type JobStore struct {
    mu   sync.RWMutex
    jobs map[string]*Job
}

var Jobs = &JobStore{jobs: make(map[string]*Job)}

func (s *JobStore) Create(j *Job) {
    s.mu.Lock()
    s.jobs[j.ID] = j
    s.mu.Unlock()
}
func (s *JobStore) Get(id string) (*Job, bool) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    j, ok := s.jobs[id]
    return j, ok
}
func (s *JobStore) SetStatus(id string, st JobStatus, err error) {
    s.mu.Lock()
    if j, ok := s.jobs[id]; ok {
        j.Status = st
        if err != nil {
            j.Error = err.Error()
        }
        j.UpdatedAt = time.Now()
    }
    s.mu.Unlock()
}
