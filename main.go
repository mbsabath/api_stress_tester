package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var rate = flag.Int("rate", 10, "rate of requests per second")
var duration = flag.Int("duration", 10, "duration of the test in seconds")
var url = flag.String("url", "http://localhost:8080", "url to test")
var method = flag.String("method", "GET", "HTTP method")
var body = flag.String("body", "", "HTTP body")
var resultsFile = flag.String("results", "", "file to write results to")

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func (r Result) String() string {
	return fmt.Sprintf("%d, %v", r.StatusCode, r.Duration.Milliseconds())
}

func main() {
	flag.Parse()
	Client := http.Client{}
	defer Client.CloseIdleConnections()
	// Create a new test
	num_requests := *rate * *duration
	if *resultsFile == "" {
		*resultsFile = fmt.Sprintf("results_%drate_%dduration_%d.csv", *rate, *duration, time.Now().Unix())
	}
	requestWg := sync.WaitGroup{}
	requestWg.Add(num_requests)
	result_channel := make(chan Result, num_requests)
	go func() {
		for i := 0; i < num_requests; i++ {
			go makeRequest(Client, *url, *method, *body, &requestWg, result_channel)
			time.Sleep(time.Second / time.Duration(*rate))
		}
	}()

	go func() {
		requestWg.Wait()
		close(result_channel)
	}()

	result_file, err := os.Create(*resultsFile)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer result_file.Close()
	// Write the header to the file
	result_file.WriteString("Status Code, Duration\n")
	for result := range result_channel {
		// Write the result to the file
		result_file.WriteString(result.String() + "\n")
	}
}

func makeRequest(client http.Client, url string, method string, body string, wg *sync.WaitGroup, result_channel chan Result) {
	defer wg.Done()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making request: %v", err)
	}
	duration := time.Since(start)
	defer resp.Body.Close()
	// Read the response body
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}
	result := Result{resp.StatusCode, duration}
	result_channel <- result
}
