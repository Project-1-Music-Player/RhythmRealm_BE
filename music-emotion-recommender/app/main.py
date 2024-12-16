from fastapi import FastAPI
from app.api.routes import router
from app.db.database import init_db

app = FastAPI(title="Music Recommendation API")

# Include routers
app.include_router(router, prefix="/api")

# Initialize database at startup
@app.on_event("startup")
async def startup_event():
    init_db()