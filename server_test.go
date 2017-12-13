package main

import (
	"testing"
	"time"
)

func getErrorBudget(errorRate, recoverRate int, timeWindow time.Duration) ErrorBudget {
	return ErrorBudget{
		ErrorRate:   errorRate,
		RecoverRate: recoverRate,
		TimeWindow:  timeWindow,
	}
}

func TestNewServer(t *testing.T) {
	errBudget := getErrorBudget(3, 2, time.Second)
	srv := NewServer(errBudget)
	if srv == nil {
		t.Error("Failed to create a server")
	}
}

func TestWorkloadSuccess(t *testing.T) {
	errBudget := getErrorBudget(3, 2, time.Second)
	srv := NewServer(errBudget)
	err := srv.Workload(true)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestWorkloadFailure(t *testing.T) {
	errBudget := getErrorBudget(3, 2, time.Second)
	srv := NewServer(errBudget)
	err := srv.Workload(false)
	if err != ErrWorkload {
		t.Errorf("Expected to get %v, but got %v", ErrWorkload, err)
	}
}
