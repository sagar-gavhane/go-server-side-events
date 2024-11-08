package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// main initializes the HTTP server and registers the SSE handler
func main() {
	// Register the /events endpoint for SSE connections
	http.HandleFunc("/events", sseHandler)

	// Start the server on port 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

// sseHandler handles Server-Sent Events connections and streams system metrics
func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Set headers required for SSE
	w.Header().Set("Content-Type", "text/event-stream") // Indicates this is an event stream
	w.Header().Set("Cache-Control", "no-cache")         // Prevents caching of events
	w.Header().Set("Connection", "keep-alive")          // Keeps the connection open
	w.Header().Set("Access-Control-Allow-Origin", "*")  // Allows cross-origin requests

	// Create tickers for periodic updates (1 second intervals)
	memT := time.NewTicker(time.Second)
	defer memT.Stop() // Ensure ticker is stopped when function exits

	cpuT := time.NewTicker(time.Second)
	defer cpuT.Stop() // Ensure ticker is stopped when function exits

	// Channel that will be closed when client disconnects
	clientGone := r.Context().Done()

	// Create response controller for flushing data
	rc := http.NewResponseController(w)

	// Infinite loop to continuously send events
	for {
		select {
		// Handle client disconnection
		case <-clientGone:
			fmt.Println("Client gone")

		// Handle memory metrics (triggers every second)
		case <-memT.C:
			// Get virtual memory statistics
			m, err := mem.VirtualMemory()

			if err != nil {
				log.Println("Failed to get memory info: ", err)
				continue
			}

			// Send memory metrics in SSE format
			// Format: event:mem\ndata:actual_data\n\n
			if _, err := fmt.Fprintf(w, "event:mem\ndata:Total: %d, Used: %d, Perc: %.2f%%\n\n",
				m.Total, m.Used, m.UsedPercent); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}

			// Ensure data is sent immediately
			rc.Flush()

		// Handle CPU metrics (triggers every second)
		case <-cpuT.C:
			// Get CPU times (false means not per CPU)
			c, err := cpu.Times(false)

			if err != nil {
				log.Println("Failed to get cpu info: ", err)
				continue
			}

			// Send CPU metrics in SSE format
			// Format: event:cpu\ndata:actual_data\n\n
			if _, err := fmt.Fprintf(w, "event:cpu\ndata:User: %.2f, Sys: %.2f, Idle: %.2f\n\n",
				c[0].User, c[0].System, c[0].Idle); err != nil {
				log.Printf("unable to write: %s", err.Error())
				return
			}

			// Ensure data is sent immediately
			rc.Flush()
		}
	}
}
