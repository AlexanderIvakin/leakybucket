package main

import (
	"errors"
	"time"
)

// ErrBudgetExceeded is returned when the server exceeds its error budget
var ErrBudgetExceeded = errors.New("error budget exceeded")

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
}

// NewServer makes a new Server
func NewServer(budget ErrorBudget) *Server {
	return &Server{budget: budget}
}

// Run executes the specified payload function
func (srv *Server) Run(f func() error) error {
	return f()
}
