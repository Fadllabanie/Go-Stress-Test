package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/fatih/color"
)

var (
	totalRequests int
	totalLatency  time.Duration
	minLatency    time.Duration
	maxLatency    time.Duration = time.Duration(0)
	totalErrors   int
)


func draw() {
	fmt.Println(`
	

	ooooooooooooooooooooooooooooooooooooo
	8                                .d88
	8  oooooooooooooooooooooooooooood8888
	8  8888888888888888888888888P"   8888    oooooooooooooooo
	8  8888888888888888888888P"      8888    8              8
	8  8888888888888888888P"         8888    8             d8
	8  8888888888888888P"            8888    8            d88
	8  8888888888888P"               8888    8           d888
	8  8888888888P"                  8888    8          d8888
	8  8888888P"   Go Stress Test    8888    8         d88888
	8  8888P"       FADL-LABANIE     8888    8        d888888
	8  8888oooooooooooooooooooooocgmm8888    8       d8888888
	8 .od88888888888888888888888888888888    8      d88888888
	8888888888888888888888888888888888888    8     d888888888
						 8    d8888888888
	   ooooooooooooooooooooooooooooooo       8   d88888888888
	  d                       ...oood8b      8  d888888888888
	 d              ...oood888888888888b     8 d8888888888888
	d     ...oood88888888888888888888888b    8d88888888888888
  	dood8888888888888888888888888888888888b
`)

}

func sendRequest(url string, wg *sync.WaitGroup, ch chan<- time.Duration) {
	defer wg.Done()

	startTime := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error sending request:", err)
		totalErrors++
		return
	}
	defer resp.Body.Close()

	// Read and discard response body to simulate a real request
	_, err = io.Copy(ioutil.Discard, resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		totalErrors++
		return
	}

	// Calculate latency for the request
	latency := time.Since(startTime)

	// Update total latency
	totalLatency += latency

	// Update min and max latency
	if latency < minLatency || minLatency == 0 {
		minLatency = latency
	}
	if latency > maxLatency {
		maxLatency = latency
	}

	// Increment total requests counter
	totalRequests++

	// Send latency to channel for histogram calculation
	ch <- latency
}

func main() {
	draw()
	var url string
	var numRequests int

	// Prompt user to input URL
	fmt.Print("Enter URL: ")
	fmt.Scanln(&url)

	// Prompt user to input number of requests
	fmt.Print("Enter number of requests: ")
	fmt.Scanln(&numRequests)

	if numRequests <= 0 {
		fmt.Println("Number of requests must be greater than zero")
		return
	}

	// Display the number of requests
	fmt.Println("Number of Requests:", numRequests)

	var wg sync.WaitGroup
	wg.Add(numRequests)

	// Channel to collect latencies for response time distribution
	latencyCh := make(chan time.Duration, numRequests)

	// Create progress bar
	bar := pb.StartNew(numRequests)

	// Start time for progress bar update
	progressStartTime := time.Now()

	// Send concurrent requests
	for i := 0; i < numRequests; i++ {
		go func() {
			sendRequest(url, &wg, latencyCh)
			// Increment progress bar
			bar.Increment()
		}()
	}

	// Wait for all requests to finish
	wg.Wait()
	close(latencyCh)

	// Stop progress bar
	bar.Finish()

	elapsedTime := time.Since(progressStartTime)

	// Calculate throughput
	throughput := float64(totalRequests) / elapsedTime.Seconds()

	// Calculate average latency
	avgLatency := totalLatency / time.Duration(totalRequests)

	// Print performance metrics with color and formatting
	fmt.Println("Performance Metrics:")
	color.Green("Total Requests: %d\n", totalRequests)
	color.Blue("=====================================================")

	color.Cyan("Throughput: %.2f requests/second\n", throughput)
	color.Blue("=====================================================")

	color.Yellow("Average Latency: %s\n", avgLatency)
	color.Blue("=====================================================")

	color.Magenta("Min Latency: %s\n", minLatency)
	color.Blue("=====================================================")

	color.Red("Max Latency: %s\n", maxLatency)
	color.Blue("=====================================================")

	color.Red("Total Errors: %d\n", totalErrors)
	color.Blue("=====================================================")

	// Generate and print response time distribution histogram
	printResponseTimeDistribution(latencyCh)
}

func printResponseTimeDistribution(latencyCh <-chan time.Duration) {
	// Histogram buckets for response time distribution (in milliseconds)
	buckets := [...]int{1, 5, 10, 20, 50, 100, 200, 500, 1000}

	// Initialize histogram counters
	histogram := make(map[int]int)
	for _, bucket := range buckets {
		histogram[bucket] = 0
	}

	// Count latencies in histogram buckets
	for latency := range latencyCh {
		latencyMS := int(latency.Milliseconds())
		for _, bucket := range buckets {
			if latencyMS <= bucket {
				histogram[bucket]++
				break
			}
		}
	}

	// Print histogram
	color.Cyan("Response Time Distribution (Milliseconds):")
	for i, bucket := range buckets {
		fmt.Printf("%d-%d ms: %d\n", 0, bucket, histogram[bucket])
		color.Blue("=====================================================")

		if i < len(buckets)-1 {
			fmt.Printf("%d-%d ms: %d\n", bucket+1, buckets[i+1], histogram[buckets[i+1]])
		}
	}
}
