"""
Database operations and setup for Fleet Monitoring API
"""

import csv
from datetime import datetime
from pathlib import Path
from sqlalchemy import create_engine, func
from sqlalchemy.orm import sessionmaker, Session
from .models import Base, DeviceModel, HeartbeatModel, StatsModel

# Database configuration
DATABASE_URL = "sqlite:///./fleet_monitoring.db"
engine = create_engine(DATABASE_URL, connect_args={"check_same_thread": False})
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def get_db():
    """Database dependency for FastAPI."""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


class DatabaseOperations:
    """Centralized database operations for device management."""

    @staticmethod
    def create_tables():
        """Create all database tables."""
        Base.metadata.create_all(bind=engine)

    @staticmethod
    def load_devices_from_csv(db: Session):
        """Load device definitions from devices.csv file."""
        csv_path = Path("app/devices.csv")
        if not csv_path.exists():
            raise FileNotFoundError("devices.csv file not found")

        # Read all device IDs from CSV
        with open(csv_path, 'r', newline='') as csvfile:
            device_ids = {row['device_id'].strip() for row in csv.DictReader(csvfile) if row['device_id'].strip()}

        # Get existing device IDs
        existing_ids = {d[0] for d in db.query(DeviceModel.device_id).all()}

        # Insert new devices in bulk
        new_devices = [DeviceModel(device_id=device_id) for device_id in device_ids - existing_ids]
        db.add_all(new_devices)
        db.commit()

    @staticmethod
    def device_exists(db: Session, device_id: str) -> bool:
        """Check if a device exists in the database."""
        return db.query(DeviceModel).filter(DeviceModel.device_id == device_id).first() is not None

    @staticmethod
    def add_heartbeat(db: Session, device_id: str, sent_at: datetime):
        """Add a heartbeat record."""
        heartbeat = HeartbeatModel(device_id=device_id, sent_at=sent_at)
        db.add(heartbeat)
        db.commit()

    @staticmethod
    def add_stats(db: Session, device_id: str, sent_at: datetime, upload_time: int):
        """Add a stats record."""
        stats = StatsModel(device_id=device_id, sent_at=sent_at, upload_time=upload_time)
        db.add(stats)
        db.commit()

    @staticmethod
    def calculate_uptime(db: Session, device_id: str) -> float:
        """Calculate uptime percentage based on heartbeat count vs expected maximum."""
        # Get count of heartbeats for the device
        heartbeat_count = db.query(HeartbeatModel.sent_at).filter(
            HeartbeatModel.device_id == device_id
            ).distinct().order_by(
                HeartbeatModel.sent_at
            ).all()

        if len(heartbeat_count) < 2:
            return 0.0

        first_heartbeat = heartbeat_count[0][0]
        last_heartbeat = heartbeat_count[-1][0]

        total_minutes = (last_heartbeat - first_heartbeat).total_seconds() / 60

        uptime = (len(heartbeat_count) / total_minutes) * 100
        return min(uptime, 100.0)

    @staticmethod
    def calculate_avg_upload_time(db: Session, device_id: str) -> str:
        """Calculate average upload time and return as duration string."""
        # Get average upload time for the device
        result = db.query(func.avg(StatsModel.upload_time)).filter(
            StatsModel.device_id == device_id
        ).scalar()

        if result is None:
            return "0s"

        avg_nanoseconds = float(result)
        return DatabaseOperations._format_duration(avg_nanoseconds)

    @staticmethod
    def _format_duration(nanoseconds: float) -> str:
        """Convert nanoseconds to duration string format like '5m10.123456789s'."""
        total_seconds = nanoseconds / 1_000_000_000

        # Build duration parts
        hours, remainder = divmod(int(total_seconds), 3600)
        minutes, _ = divmod(remainder, 60)
        seconds = total_seconds % 60

        parts = []
        if hours:
            parts.append(f"{hours}h")
        if minutes:
            parts.append(f"{minutes}m")
        if seconds > 0 or not parts:
            parts.append(f"{seconds:.9f}".rstrip('0').rstrip('.') + "s")

        return "".join(parts)

    @staticmethod
    def get_all_devices(db: Session) -> list[str]:
        """Get all device IDs from the database."""
        devices = db.query(DeviceModel.device_id).all()
        return [device[0] for device in devices]

    @staticmethod
    def get_device_counts(db: Session, device_id: str) -> tuple[int, int]:
        """Get heartbeat and stats counts for a device."""
        heartbeat_count = db.query(HeartbeatModel).filter(HeartbeatModel.device_id == device_id).count()
        stats_count = db.query(StatsModel).filter(StatsModel.device_id == device_id).count()
        return heartbeat_count, stats_count
