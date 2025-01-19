package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func handleDivide(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	tags := map[string]string{}

	// add dummy process to make sure we can get latency more than 1 millisecond
	time.Sleep(time.Duration(rand.Intn(1000)+10) * time.Millisecond)

	defer func() {
		pushHistogram(r.Context(), time.Since(start).Milliseconds(), tags)
	}()

	a, b, err := parseQueryParams(r)
	if err != nil {
		tags["error"] = err.Error()
		tags["status"] = "error"
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if b == 0 {
		tags["error"] = "division by zero"
		tags["status"] = "error"
		http.Error(w, "Cannot divide by zero", http.StatusBadRequest)
		return
	}

	result := float64(a) / float64(b)
	_, _ = fmt.Fprintf(w, "Division of %d by %d is %.2f", a, b, result)

	tags["status"] = "success"

}

func parseQueryParams(r *http.Request) (int, int, error) {
	query := r.URL.Query()
	aStr := query.Get("a")
	bStr := query.Get("b")

	if aStr == "" || bStr == "" {
		return 0, 0, fmt.Errorf("missing query parameters 'a' or 'b'")
	}

	a, err := strconv.Atoi(aStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value for 'a': %s", aStr)
	}

	b, err := strconv.Atoi(bStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid value for 'b': %s", bStr)
	}

	return a, b, nil
}
