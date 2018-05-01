package main

import (
    "flag"
    "fmt"
    "log"
    "math"
    "net/http"
    "sync"
    "time"
)

var mutex = &sync.Mutex{}
var requests = 0
var requestCPUDuration = flag.Duration("request-cpu-burn-duration", 200 * time.Millisecond, "Duration of heavy CPU burn done on each request")
var requestSleepDuration = flag.Duration("request-sleep-duration", 100 * time.Millisecond, "Duration of sleep on each request")
var requestDuration = 200 * time.Millisecond

func handler(w http.ResponseWriter, r *http.Request) {
	// Process one request at a time
	// so that latency grows when system is overloaded
	// (aka. pretend we're python)
	mutex.Lock()
	defer mutex.Unlock()

	requests += 1

	start := time.Now()

	for time.Now().Sub(start) < *requestCPUDuration {
		// Burn some CPU on pointless math.
		x := 0.0001
		for i := 0; i < 100000; i++ {
			x += math.Sqrt(x)
		}
	}
	time.Sleep(*requestSleepDuration)
	fmt.Fprintf(w, "OK")
}

func metrics(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "# HELP http_requests_total The amount of requests served by the server in total\n# TYPE http_requests_total counter\nhttp_requests_total %v\n", requests)
}

func main() {
    flag.Parse()
    http.HandleFunc("/app", handler)
    http.HandleFunc("/metrics", metrics)
    log.Fatal(http.ListenAndServe(":1234", nil))
}
