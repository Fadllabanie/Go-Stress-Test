# Go Stress Test

Go Stress Test is a simple tool written in Go for stress testing HTTP endpoints. It sends a specified number of requests to a given URL and measures various performance metrics such as throughput, latency, and error rate.

## Features

- Stress test HTTP endpoints by sending a configurable number of requests.
- Measure throughput, average latency, minimum and maximum latency, and total errors.
- Display response time distribution histogram for better insight into request latencies.

## Installation

To use Go Stress Test, you need to have Go installed on your system.

```bash
go get github.com/YourRepo/go-stress-test
```


## Usage
To run Go Stress Test, execute the following command:
```bash
go-stress-test -url <URL> -requests <Number of Requests>
```

- url: URL of the endpoint to test.
- requests: Number of requests to send (default is 100).

## Example:
```bash
go-stress-test -url https://example.com/api -requests 1000
```
## Performance Metrics
After running the stress test, the tool provides the following performance metrics:

- Total Requests
- Throughput (requests/second)
- Average Latency
- Minimum Latency
- Maximum Latency
- Total Errors
Additionally, it generates and displays a response time distribution histogram to visualize request latencies.


## build
- macos build 
- windows build 
- linux build 

you can download build and work on it without clone project
just after run build he will ask you about endpoint - number of requests

## Credits
This tool utilizes the following third-party libraries:

- cheggaaa/pb: For displaying progress bars.
- fatih/color: For colorful console output.

## Author
This stress testing tool is developed by FADL-LABANIE.

## License
This project is licensed under the MIT License - see the LICENSE file for details.