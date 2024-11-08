# System Monitor with Server-Sent Events

A real-time system monitoring application that displays CPU and memory usage statistics using Server-Sent Events (SSE). This project consists of a Go backend server that streams system metrics and a simple HTML frontend to display the data.

## Features

- Real-time CPU usage monitoring (User, System, and Idle percentages)
- Real-time Memory usage monitoring (Total, Used, and Usage percentage)
- Server-Sent Events for efficient real-time updates
- Cross-origin resource sharing (CORS) enabled
- Simple and lightweight frontend

## Prerequisites

- Go 1.16 or higher
- [gopsutil](https://github.com/shirou/gopsutil) library

## Installation

1. Clone the repository:

```bash
git clone https://github.com/sagar-gavhane/go-server-side-events.git
cd go-server-side-events
```

2. Install dependencies:

```bash
go mod tidy
```

3. Run the server:

```bash
go run main.go
```

4. Open `index.html` in your browser to see the system monitor in action.

## Code Explanation

- `main.go` contains the server code that streams the system metrics to the client.
- `index.html` contains the frontend code that displays the system metrics.

This is a server program that monitors system resources (CPU and Memory) and sends real-time updates to web clients using Server-Sent Events (SSE).

### Code Breakdown:

1. **Package and Imports**
   - The program uses several Go packages including standard ones for HTTP handling and logging
   - It uses `gopsutil` library to get CPU and memory statistics

2. **Main Function**
   - Sets up a web server that listens on port 8080
   - Registers the `/events` endpoint to handle SSE connections
   - If the server fails to start, it logs the error and exits

3. **SSE Handler Function (`sseHandler`)**
   - This is the main function that handles client connections and sends system stats
   - It sets up the necessary headers for SSE:
     - Tells the client it's an event stream
     - Disables caching
     - Keeps the connection alive
     - Allows cross-origin requests (from any domain)

4. **Monitoring Loop**
   - Creates two timers that tick every second:
     - One for memory stats
     - One for CPU stats
   - Sets up a way to detect when the client disconnects
   - Enters an infinite loop that:

   **For Memory:**
   - Every second, gets the current memory statistics
   - Sends an event named "mem" with:
     - Total memory
     - Used memory
     - Usage percentage

   **For CPU:**
   - Every second, gets the current CPU statistics
   - Sends an event named "cpu" with:
     - User CPU usage
     - System CPU usage
     - Idle CPU percentage

5. **Error Handling**
   - If there's an error getting stats, it logs the error but continues running
   - If there's an error sending data to the client, it ends the connection
   - If the client disconnects, it detects this and cleans up

6. **Data Streaming**
   - Uses `Flush()` after sending each update to ensure the client receives the data immediately
   - Formats the data in the SSE format (event:type\ndata:value)

This server acts as a continuous data stream, sending system statistics to any connected web clients in real-time, which can then display this information in a web interface.