package main

import (
	"errors"
	"time"
)

// ErrWorkload is the error returned by the Workload function
// in case of failure
var ErrWorkload = errors.New("workload function failed")

// ErrorBudget contains error budget information for the Server
type ErrorBudget struct {
	// ErrorRate is the number of errors per given time window
	ErrorRate int
	// RecoverRate is the number of errors to recover per given time window
	RecoverRate int
	// TimeWindow is the size of the time window for ErrorRate and RecoverRate
	TimeWindow time.Duration
}

// Server is a fictitious server performing a Workload function
type Server struct {
	budget ErrorBudget
}

// NewServer makes a new Server
func NewServer(budget ErrorBudget) *Server {
	return &Server{budget: budget}
}

// Workload is the Server's function
func (srv *Server) Workload(success bool) error {
	if success {
		return nil
	}
	return ErrWorkload
}
