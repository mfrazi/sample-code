package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

// Make sure to run main function before running this function
func TestPopulateMetricHistogram(t *testing.T) {
	url := "http://localhost:8080/divide?a=%d&b=%d"
	minNum, maxNum := -5, 5
	maxRunningTime := time.Duration(10)

	start := time.Now()
	for time.Since(start) < maxRunningTime*time.Minute {
		a := rand.Intn(maxNum-minNum+1) + minNum
		b := rand.Intn(maxNum-minNum+1) + minNum
		finalURL := fmt.Sprintf(url, a, b)

		resp, err := http.Get(finalURL)
		if err == nil {
			_ = resp.Body.Close()
		}

		// Wait for a random interval between 1 and 5 seconds
		time.Sleep(time.Duration(rand.Intn(5)+1) * time.Second)
	}
}
