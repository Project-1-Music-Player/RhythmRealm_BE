from fastapi import FastAPI
from app.api.routes import router
from app.db.database import init_db
import os
from fastapi.middleware.cors import CORSMiddleware


app = FastAPI(title="Music Recommendation API")
origins = [
    "http://localhost:3001",
    "http://localhost:8000",  
]
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"], 
    allow_headers=["*"], 
)

# Include routers
app.include_router(router, prefix="/api")

# Initialize database at startup
@app.on_event("startup")
async def startup_event():
    init_db()