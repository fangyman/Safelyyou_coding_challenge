"""
Debug API Endpoints
Provides debugging endpoints for development and troubleshooting
"""

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from app.database import DatabaseOperations, get_db
from app.schemas import (
    DeviceListResponse,
    DeviceDebugResponse
)

router = APIRouter(prefix="/debug", tags=["debug"])


@router.get("/devices", response_model=DeviceListResponse)
async def debug_devices(db: Session = Depends(get_db)):
    """Debug endpoint to list all devices."""
    devices = DatabaseOperations.get_all_devices(db)
    return DeviceListResponse(devices=devices)


@router.get("/stats/{device_id}", response_model=DeviceDebugResponse)
async def debug_device_data(device_id: str, db: Session = Depends(get_db)):
    """Debug endpoint to show raw data for a device."""
    heartbeat_count, stats_count = DatabaseOperations.get_device_counts(db, device_id)

    return DeviceDebugResponse(
        device_id=device_id,
        heartbeat_count=heartbeat_count,
        stats_count=stats_count
    )
