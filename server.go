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
	// RecoveryRate is the number of errors to recover per given time window
	RecoveryRate int
	// TimeWindow is the size of the time window for ErrorRate and RecoveryRate
	TimeWindow time.Duration
}

// Server is a fictitious server performing a Workload function
type Server struct {
	// budget is the error budget assigned for the server
	budget ErrorBudget
	// Workload is a payload-function performed by the server
	Workload func(interface{}, <-chan interface{}) (<-chan interface{}, error)
}

// NewServer makes a new Server
func NewServer(budget ErrorBudget) *Server {
	return &Server{budget: budget}
}
