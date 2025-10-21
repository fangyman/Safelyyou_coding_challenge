"""
SQLAlchemy ORM Models for Fleet Monitoring API
"""

from sqlalchemy import Column, String, DateTime, Integer
from sqlalchemy.ext.declarative import declarative_base

Base = declarative_base()


class DeviceModel(Base):
    """Database model for devices."""
    __tablename__ = "devices"

    device_id = Column(String, primary_key=True, index=True)


class HeartbeatModel(Base):
    """Database model for heartbeat records."""
    __tablename__ = "heartbeats"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    device_id = Column(String, index=True)
    sent_at = Column(DateTime)


class StatsModel(Base):
    """Database model for stats records."""
    __tablename__ = "stats"

    id = Column(Integer, primary_key=True, index=True, autoincrement=True)
    device_id = Column(String, index=True)
    sent_at = Column(DateTime)
    upload_time = Column(Integer)  # nanoseconds
