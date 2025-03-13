package main

import (
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Configuration for the load test
type LoadTestConfig struct {
	Concurrent int
	Duration   time.Duration
	BaseURL    string
}

// Statistics for the load test
type Stats struct {
	RequestCount int
	SuccessCount int
	ErrorCount   int
	TotalTime    time.Duration
	mutex        sync.Mutex
}

var endpoints = []string{
	"/test/fast",
	"/test/slow",
	"/test/random",
	"/test/error",
	"/test/ok",
}

func (s *Stats) recordRequest(success bool, startTime time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.RequestCount++
	if success {
		s.SuccessCount++
	} else {
		s.ErrorCount++
	}
	s.TotalTime += time.Since(startTime)
}

func main() {
	concurrent := flag.Int("concurrent", 100, "Number of concurrent users")
	duration := flag.Duration("duration", 30*time.Second, "Duration of the test")
	baseURL := flag.String("url", "http://127.0.0.1:8080", "Base URL of the API")
	flag.Parse()

	config := LoadTestConfig{
		Concurrent: *concurrent,
		Duration:   *duration,
		BaseURL:    *baseURL,
	}

	fmt.Printf("Starting load test with %d concurrent users for %s\n", config.Concurrent, config.Duration)
	fmt.Printf("Target URL: %s\n\n", config.BaseURL)

	stats := Stats{}
	endTime := time.Now().Add(config.Duration)

	var wg sync.WaitGroup

	for i := range config.Concurrent {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			worker(id, &stats, endTime, config.BaseURL)
		}(i)
	}

	wg.Wait()

	fmt.Println("\nLoad Test Results:")
	fmt.Printf("Total Requests: %d\n", stats.RequestCount)
	fmt.Printf("Successful Requests: %d\n", stats.SuccessCount)
	fmt.Printf("Failed Requests: %d\n", stats.ErrorCount)

	if stats.RequestCount > 0 {
		fmt.Printf("Success Rate: %.2f%%\n", float64(stats.SuccessCount)/float64(stats.RequestCount)*100)
		avgTime := stats.TotalTime.Seconds() / float64(stats.RequestCount)
		fmt.Printf("Average Response Time: %.2f ms\n", avgTime*1000)
		fmt.Printf("Requests per second: %.2f\n", float64(stats.RequestCount)/config.Duration.Seconds())
	}
}

func worker(_ int, stats *Stats, endTime time.Time, baseURL string) {
	client := &http.Client{
		Timeout: 6 * time.Second,
	}

	for time.Now().Before(endTime) {
		for _, endpoint := range endpoints {
			url := baseURL + endpoint
			timeStart := time.Now()
			resp, err := client.Get(url)

			if err != nil {
				stats.recordRequest(false, timeStart)
				continue
			}

			resp.Body.Close()
			stats.recordRequest(resp.StatusCode >= 200 && resp.StatusCode < 400, timeStart)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
