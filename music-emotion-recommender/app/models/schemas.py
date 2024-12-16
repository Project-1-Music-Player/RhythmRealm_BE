from typing import List, Optional
from pydantic import BaseModel

class SpotifyTrackInfo(BaseModel):
    preview_url: Optional[str]
    external_url: Optional[str]
    album_name: Optional[str]
    album_image: Optional[str]
    duration_ms: Optional[int]
    popularity: Optional[int]

class SongResponse(BaseModel):
    track: str
    artist: str
    genre: Optional[str]
    spotify_id: Optional[str]
    tags: List[str]
    score: float
    valence: float
    arousal: float
    dominance: float
    spotify_info: Optional[SpotifyTrackInfo]