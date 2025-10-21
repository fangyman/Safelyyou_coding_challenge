"""
Fleet Monitoring API
SafelyYou Coding Challenge

A FastAPI-based REST API for monitoring device fleet uptime and performance.
Uses SQLite database for data persistence with modular architecture.
"""

from contextlib import asynccontextmanager
import uvicorn
from fastapi import FastAPI

from app.database import DatabaseOperations, SessionLocal
from apis.devices import router as devices_router
from apis.debug import router as debug_router
from apis.root import router as root_router

@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan - handles startup and shutdown events."""
    # Startup
    DatabaseOperations.create_tables()

    # Load devices from CSV
    db = SessionLocal()
    try:
        DatabaseOperations.load_devices_from_csv(db)
    finally:
        db.close()

    yield


app = FastAPI(
    title="Fleet Monitoring API",
    description="Device monitoring API for SafelyYou fleet management with SQLite persistence",
    version="1.0.0",
    lifespan=lifespan
)

# Include API routers
app.include_router(root_router)
app.include_router(devices_router)
app.include_router(debug_router)


if __name__ == "__main__":
    # Run the server on port 6733 as expected by the device simulator
    uvicorn.run(app, host="127.0.0.1", port=6733)