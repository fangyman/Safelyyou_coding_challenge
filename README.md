# Fleet Monitoring API

A FastAPI-based REST API for monitoring device fleet uptime and performance, built as part of the SafelyYou coding challenge. The system tracks heartbeat telemetry and performance statistics from IoT devices, calculating uptime percentages and average upload times.

## Features

- **Device Registration**: Automatic device loading from CSV configuration
- **Heartbeat Monitoring**: Track device connectivity via periodic heartbeat signals
- **Performance Statistics**: Monitor and calculate average upload times
- **SQLite Persistence**: Lightweight database storage with SQLAlchemy ORM
- **Debug Endpoints**: Development and troubleshooting tools
- **OpenAPI Documentation**: Interactive API documentation at `/docs`

## Project Structure

```text
├── main.py                    # FastAPI application entry point
├── requirements.txt           # Python dependencies
├── fleet_monitoring.db        # SQLite database (created on first run)
├── device-simulator-win-amd64.exe  # Device simulator for testing
├── app/
│   ├── database.py           # Database operations and SQLAlchemy setup
│   ├── models.py             # SQLAlchemy ORM models
│   ├── schemas.py            # Pydantic request/response models
│   └── devices.csv           # Device configuration file
└── apis/
    ├── devices.py            # Device telemetry endpoints
    ├── debug.py              # Debug and monitoring endpoints
    └── root.py               # Basic application info endpoints
```

## Installation & Setup

### Prerequisites

- Python 3.8 or higher
- Windows (for the provided device simulator)

### Installation Steps

1. **Clone or download the project**

2. **Create and activate a virtual environment** (recommended)

   ```bash
   python -m venv venv
   venv\Scripts\activate  # On Windows
   ```

3. **Install dependencies**

   ```bash
   pip install -r requirements.txt
   ```

## How to Run

### 1. Start the API Server

```bash
python main.py
```

The server will start on `http://127.0.0.1:6733` and will:

- Create SQLite database tables automatically
- Load device definitions from `app/devices.csv`
- Be ready to receive telemetry data

### 2. Run the Device Simulator (Optional)

To test the API with simulated device data:

```bash
.\device-simulator-win-amd64.exe
```

This will:

- Send heartbeat and stats telemetry to your API
- Generate test results comparing expected vs actual calculations
- Output results to `results.txt`

### 3. Access the API

- **Interactive Documentation**: <http://127.0.0.1:6733/docs>
- **Root Endpoint**: <http://127.0.0.1:6733/>
- **Device Stats Example**: <http://127.0.0.1:6733/api/v1/devices/60-6b-44-84-dc-64/stats>
- **Debug Device List**: <http://127.0.0.1:6733/debug/devices>

## API Endpoints

### Device Telemetry

- `POST /api/v1/devices/{device_id}/heartbeat` - Receive heartbeat signals
- `POST /api/v1/devices/{device_id}/stats` - Receive performance statistics
- `GET /api/v1/devices/{device_id}/stats` - Get calculated device statistics

### Debug & Monitoring

- `GET /debug/devices` - List all registered devices
- `GET /debug/stats/{device_id}` - Get raw telemetry data counts

### Application Info

- `GET /` - Basic API information and version

## Solution Write-up

### Time Spent & Challenges

**Development Time**: Approximately 5-8 hours total, spread across design, implementation, and testing phases.

**Most Difficult Parts**:

1. **Duration formatting** - Implementing the specific duration string format (e.g., "3m21.858747766s") that matches the expected output format, especially handling nanosecond precision.

2. **Database design decisions** - Balancing simplicity with performance, particularly around indexing strategies and whether to store calculated values or compute them on-demand.

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

1. **Database Indexing** (Immediate improvement):

   ```sql
   CREATE INDEX idx_heartbeats_device_time ON heartbeats (device_id, sent_at);
   CREATE INDEX idx_stats_device_time ON stats (device_id, sent_at);
   ```

   - Improves WHERE clause performance to O(log n)
   - Still O(n) for full aggregations but with better constants

2. **Pre-computed Aggregates** (Near O(1) performance):

   ```sql
   CREATE TABLE device_metrics_cache (
       device_id TEXT PRIMARY KEY,
       total_heartbeats INTEGER,
       first_heartbeat TIMESTAMP,
       last_heartbeat TIMESTAMP,
       total_upload_time BIGINT,
       stats_count INTEGER,
       last_updated TIMESTAMP
   );
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
