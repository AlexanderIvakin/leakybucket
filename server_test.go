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
		t.Error("Failed to create a server")
	}
}

func TestBoringWorkload(t *testing.T) {
	t.Parallel()

	errBudget := getErrorBudget(3, 2, time.Second)

	srv := NewServer(errBudget)

	// Specify "boring" workload - output 42 until cancelled
	srv.Workload = func(arg interface{}, term <-chan interface{}) (<-chan interface{}, error) {
		resChan := make(chan interface{})
		go func() {
			defer close(resChan)
			for {
				select {
				case <-term:
					return
				case resChan <- 42:
				}
			}
		}()
		return resChan, nil
	}

	term := make(chan interface{})
	defer close(term)
	resChan, err := srv.Workload(nil, term)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for i := 0; i < 1000; i++ {
		res := <-resChan
		intRes, _ := res.(int)
		if intRes != 42 {
			t.Errorf("Expected to get %d, got - %d", 42, intRes)
		}
	}
}

func TestInterestingWorkload(t *testing.T) {
	t.Parallel()

	errBudget := getErrorBudget(3, 2, 10*time.Millisecond)

	srv := NewServer(errBudget)

	// Specify "interesting" workload that will fail arbitrarily
	// 0-10 times in a 10 millisecond window
	srv.Workload = func(arg interface{}, term <-chan interface{}) (<-chan interface{}, error) {
		resChan := make(chan interface{})
		go func() {
			defer close(resChan)
			const maxErrors int = 10
			windowTicker := time.NewTicker(10 * time.Millisecond)
			rnd := rand.New(rand.NewSource(1337))
			totalErrors := rnd.Intn(maxErrors)
			errorsCounter := 0
			defer windowTicker.Stop()
			for {
				select {
				case <-term:
					return
				case <-windowTicker.C:
					totalErrors = rnd.Intn(maxErrors)
					errorsCounter = 0
				default:
					if errorsCounter < totalErrors {
						coinflip := rnd.Intn(2)
						if coinflip == 0 {
							// throw an error
							select {
							case <-term:
								return
							case resChan <- errors.New("Boom"):
							}
							errorsCounter++
						} else {
							// normal "boring" computation
							select {
							case <-term:
								return
							case resChan <- 42:
							}
						}
					} else {
						// normal "boring" computation
						select {
						case <-term:
							return
						case resChan <- 42:
						}
					}
				}
			}
		}()
		return resChan, nil
	}

	term := make(chan interface{})
	defer close(term)
	resChan, err := srv.Workload(nil, term)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	for i := 0; i < 1000; i++ {
		res := <-resChan
		intRes, _ := res.(int)
		if intRes != 42 {
			t.Errorf("Expected to get %d, got - %d", 42, intRes)
		}
	}
}
