"""
Device API Endpoints
Handles device heartbeats and stats endpoints
"""

from fastapi import APIRouter, HTTPException, status, Depends
from sqlalchemy.orm import Session

from app.database import DatabaseOperations, get_db
from app.schemas import (
    HeartbeatRequest,
    StatsRequest,
    StatsResponse
)

router = APIRouter(prefix="/api/v1/devices", tags=["devices"])


@router.post("/{device_id}/heartbeat", status_code=status.HTTP_204_NO_CONTENT)
async def receive_heartbeat(device_id: str, heartbeat: HeartbeatRequest, db: Session = Depends(get_db)):
    """Receive heartbeat telemetry from a device."""
    if not DatabaseOperations.device_exists(db, device_id):
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device {device_id} not found"
        )

    DatabaseOperations.add_heartbeat(db, device_id, heartbeat.sent_at)


@router.post("/{device_id}/stats", status_code=status.HTTP_204_NO_CONTENT)
async def receive_stats(device_id: str, stats: StatsRequest, db: Session = Depends(get_db)):
    """Receive stats telemetry from a device."""
    if not DatabaseOperations.device_exists(db, device_id):
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device {device_id} not found"
        )

    DatabaseOperations.add_stats(db, device_id, stats.sent_at, stats.upload_time)


@router.get("/{device_id}/stats", response_model=StatsResponse)
async def get_device_stats(device_id: str, db: Session = Depends(get_db)):
    """Get calculated statistics for a device."""
    if not DatabaseOperations.device_exists(db, device_id):
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Device {device_id} not found"
        )

    uptime = DatabaseOperations.calculate_uptime(db, device_id)
    avg_upload_time = DatabaseOperations.calculate_avg_upload_time(db, device_id)

    return StatsResponse(
        uptime=uptime,
        avg_upload_time=avg_upload_time
    )
