package main

import (
	"errors"
	"math/rand"
	"testing"
	"time"
)

func getErrorBudget(errorRate, recoveryRate int, timeWindow time.Duration) ErrorBudget {
	return ErrorBudget{
		ErrorRate:    errorRate,
		RecoveryRate: recoveryRate,
		TimeWindow:   timeWindow,
	}
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	errBudget := getErrorBudget(3, 2, time.Second)
	srv := NewServer(errBudget)
	if srv == nil {
		t.Fatal("Failed to create a server")
	}
}

func TestErrorRateLimiting(t *testing.T) {
	t.Parallel()

	errorRate := 3
	recoveryRate := 3 // same recovery rate as error rate for simplicity
	timeWindow := 10 * time.Millisecond

	errBudget := getErrorBudget(errorRate, recoveryRate, timeWindow)
	srv := NewServer(errBudget)

	rnd := rand.New(rand.NewSource(1337))
	totalErrors := rnd.Intn(errorRate * 2)
	errorCount := 0

	failFunc := func() error { return errors.New("fail") }
	passFunc := func() error { return nil }

	totalTimeWindows := 10
	windowCount := 0
	ticker := time.NewTicker(timeWindow)
	defer ticker.Stop()
LOOP:
	for {
		select {
		case <-ticker.C:
			if windowCount == totalTimeWindows {
				break LOOP
			}
			totalErrors = rnd.Intn(errorRate * 2)
			errorCount = 0
			windowCount++
		default:
			var err error
			if errorCount < totalErrors {
				// fail
				err = srv.Run(failFunc)
				errorCount++
			} else {
				// pass
				err = srv.Run(passFunc)
			}

			if errorCount >= errorRate {
				if err != ErrBudgetExceeded {
					t.Errorf("Expected to get %v, but got - %v",
						ErrBudgetExceeded, err)
				}
			}
		}
	}
}
