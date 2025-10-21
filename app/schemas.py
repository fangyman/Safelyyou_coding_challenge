"""
Pydantic schemas for request/response models
"""

from datetime import datetime
from pydantic import BaseModel


class HeartbeatRequest(BaseModel):
    """Request model for heartbeat endpoint."""
    sent_at: datetime


class StatsRequest(BaseModel):
    """Request model for stats endpoint."""
    sent_at: datetime
    upload_time: int


class StatsResponse(BaseModel):
    """Response model for device stats endpoint."""
    uptime: float
    avg_upload_time: str  # Duration string format like "5m10s"


class DeviceListResponse(BaseModel):
    """Response model for device list debug endpoint."""
    devices: list[str]


class DeviceDebugResponse(BaseModel):
    """Response model for device debug endpoint."""
    device_id: str
    heartbeat_count: int
    stats_count: int
