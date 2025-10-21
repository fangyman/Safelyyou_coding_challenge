"""
Root API Endpoints
Provides basic application information endpoints
"""

from fastapi import APIRouter

router = APIRouter(tags=["root"])


@router.get("/")
async def root():
    """Root endpoint with basic API information."""
    return {
        "message": "Fleet Monitoring API with SQLite",
        "version": "1.0.0",
        "docs": "/docs",
        "database": "SQLite",
        "architecture": "Modular"
    }
