# Fleet Monitoring API

A Gin-based REST API for monitoring device fleet uptime and performance, built as part of the SafelyYou coding challenge. The system tracks heartbeat telemetry and performance statistics from IoT devices, calculating uptime percentages and average upload times.

## Features

- **Device Registration**: Automatic device loading from CSV configuration
- **Heartbeat Monitoring**: Track device connectivity via periodic heartbeat signals
- **Performance Statistics**: Monitor and calculate average upload times
- **SQLite Persistence**: Lightweight database storage with GORM

## Project Structure

```text
├── main.go                      # Main application entry point (Gin server)
├── go.mod                       # Go module dependencies
├── device-simulator-win-amd64.exe  # Device simulator for testing
├── devices.csv                  # Device configuration file
├── api/
│   └── handlers.go              # HTTP request handlers for API endpoints
├── database/
│   ├── config.go                # Database configuration and connection
│   └── database.go              # Database operations and GORM setup
├── models/
│   └── models.go                # GORM ORM models (Device, Heartbeat, etc.)
└── services/
    ├── csv_service.go           # CSV processing service for device loading
    └── metrics_service.go       # Metrics calculation and monitoring service
```

## Installation & Setup

### Prerequisites

- Go 1.19 or higher
- Windows (for the provided device simulator)

### Installation Steps

1. **Clone or download the project**

2. **Make sure you have go installed**

   **Windows:**
   - Download Go from https://golang.org/dl/
   - Run the installer (`.msi` file)
   - Add Go to your PATH (the installer usually does this automatically)
   - Open a new Command Prompt or PowerShell and verify installation:

     ```bash
     go version
     ```

   **Mac:**
   - **Option 1 - Using Homebrew (recommended):**

     ```bash
     brew install go
     ```

   - **Option 2 - Direct download:**
     - Download Go from https://golang.org/dl/
     - Run the installer (`.pkg` file)
     - Add Go to your PATH by adding to your shell profile (`~/.bash_profile`, `~/.zshrc`, etc.):

      ```bash
       export PATH=$PATH:/usr/local/go/bin
      ```

   - Verify installation:

     ```bash
     go version
     ```

   **Linux (Ubuntu/Debian):**

   ```bash
   sudo apt update
   sudo apt install golang-go
   go version
   ```

3. **Install dependencies**

   ```bash
   go mod tidy
   ```

## How to Run

### 1. Start the API Server

```bash
go run main.go
```

The server will start on `http://127.0.0.1:6733` and will:

- Create SQLite database tables automatically
- Load device definitions from `devices.csv`
- Be ready to receive telemetry data

### 2. Run the Device Simulator

To test the API with simulated device data:

```bash
.\device-simulator-win-amd64.exe
```

This will:

- Send heartbeat and stats telemetry to your API
- Generate test results comparing expected vs actual calculations
- Output results to `results.txt`

### 3. Access the API

- **Root Endpoint**: <http://127.0.0.1:6733/>

## API Endpoints

### Device Telemetry

- `POST /api/v1/devices/{device_id}/heartbeat` - Receive heartbeat signals
- `POST /api/v1/devices/{device_id}/stats` - Receive performance statistics
- `GET /api/v1/devices/{device_id}/stats` - Get calculated device statistics

### Application Info

- `GET /` - Basic API information and version

## Solution Write-up

### Time Spent & Challenges

**Development Time**: Approximately 4-5 hours total, spread across learning the documentation for gin and gorm.

**Most Difficult Parts**:

1. **Database design decisions** - Balancing simplicity with performance, particularly around reading and writing to the sqlite database and understanding various configurations to achieve best write speeds for the application.

2. **Go Learning curve** - Since I implemented the project in Python, I was able to quickly implement the Go version by following a similar structure to the python version with implementing ORM based models and HTTP API endpoints.

3. **Modularizing the Go code** - In my initial rewrite to Go, my `database.go` file was becoming to large with various different methods to handle and return data from the database. Once I realized that it's becoming too convoluted I then seperated the database operations to be held in `database.go` and moved the various different service functions and database config to their own files to maintain modularity and also readability of the code that I have written.

### Scalability & Data Model Extensions

**To support more kinds of metrics**, the current architecture could be extended to:

1. **Generic Metrics Table**:

   ```sql
   CREATE TABLE metrics (
       id INTEGER PRIMARY KEY,
       device_id TEXT,
       metric_type TEXT,  -- 'heartbeat', 'upload_time', 'cpu_usage', 'memory', etc.
       value REAL,        -- Numeric value
       metadata JSON,     -- Additional metric-specific data
       sent_at TIMESTAMP
   );
   ```

### Runtime Complexity Analysis

**Current Implementation Complexities**:

1. **Heartbeat Storage**: `O(1)` - Simple INSERT operation
2. **Stats Storage**: `O(1)` - Simple INSERT operation
3. **Uptime Calculation**: `O(n)` where n = number of heartbeats for the device
   - Queries all heartbeats for a device
   - Sorts by timestamp
   - Could be optimized with indexing and windowing functions
4. **Average Upload Time**: `O(n)` where n = number of stats records for the device
   - Uses SQL AVG() function which is generally optimized by the database

**Optimization Opportunities**:

1. **Composite Indexing** (Immediate improvement):

   ```go
   type HeartbeatModel struct {
       DeviceID string `gorm:"index:idx_heartbeats_device_time;foreignKey:DeviceID"`
       ID       uint   `gorm:"primaryKey;autoincrement;"`
       SentAt   time.Time `gorm:"index:idx_heartbeats_device_time"`
   }

   type StatsModel struct {
       DeviceID   string    `gorm:"index:idx_stats_device_time;foreignKey:DeviceID"`
       ID         uint      `gorm:"primaryKey;autoincrement;"`
       SentAt     time.Time `gorm:"column:sent_at;index:idx_stats_device_time"`
       UploadTime int64     `gorm:"column:upload_time"`
   }
   ```

   - Improves WHERE clause performance to O(log n)
   - Still O(n) for full aggregations but with better constants

2. **Pre-computed Aggregates** (Near O(1) performance):

   ```go
   type DeviceMetricsCache struct {
       DeviceID         string    `gorm:"primaryKey;column:device_id"`
       TotalHeartbeats  int       `gorm:"column:total_heartbeats"`
       FirstHeartbeat   time.Time `gorm:"column:first_heartbeat"`
       LastHeartbeat    time.Time `gorm:"column:last_heartbeat"`
       TotalUploadTime  int64     `gorm:"column:total_upload_time"`
       StatsCount       int       `gorm:"column:stats_count"`
       LastUpdated      time.Time `gorm:"column:last_updated"`
   }

   func (DeviceMetricsCache) TableName() string { return "device_metrics_cache" }
   ```

   - Update incrementally on each telemetry insert
   - Query complexity becomes O(1)

**Scalability Characteristics**:

- **Memory Usage**: `O(d × t)` where d = devices, t = time period of stored data
- **Query Performance**: Currently O(n) with data volume per device, can be improved with indexing and caching strategies
- **Storage Growth**: Linear with telemetry volume, manageable with data retention policies

The current solution prioritizes simplicity and correctness over optimization, making it suitable for moderate-scale deployments while providing a solid foundation for scaling optimizations.

## Testing Results

The solution successfully passes all test cases from the device simulator:

- **Uptime calculations**: All devices show 100% accuracy between expected and actual values
- **Average upload time**: Precision maintained to nanosecond level
- **All device IDs**: Successfully handled 5 test devices with varying telemetry patterns

See `results.txt` for detailed test output showing exact matches between expected and calculated values.
