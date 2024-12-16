from typing import List, Dict
import os
import spotipy
from spotipy.oauth2 import SpotifyClientCredentials
from dotenv import load_dotenv
from app.models.schemas import SpotifyTrackInfo

# Load environment variables
load_dotenv()

# Initialize Spotify client
spotify = spotipy.Spotify(
    client_credentials_manager=SpotifyClientCredentials(
        client_id=os.getenv('SPOTIFY_CLIENT_ID'),
        client_secret=os.getenv('SPOTIFY_CLIENT_SECRET')
    )
)

def get_spotify_track_infos(spotify_ids: List[str]) -> Dict[str, SpotifyTrackInfo]:
    """Fetch track information from Spotify in batches"""
    result = {}
    batch_size = 50  # Spotify API allows up to 50 tracks per request
    
    try:
        # Process spotify_ids in batches
        for i in range(0, len(spotify_ids), batch_size):
            batch_ids = spotify_ids[i:i + batch_size]
            batch_ids = [id for id in batch_ids if id]  # Filter out None/empty values
            
            if not batch_ids:
                continue
                
            tracks = spotify.tracks(batch_ids)['tracks']
            
            for track in tracks:
                if track:
                    result[track['id']] = SpotifyTrackInfo(
                        preview_url=track.get('preview_url'),
                        external_url=track.get('external_urls', {}).get('spotify'),
                        album_name=track.get('album', {}).get('name'),
                        album_image=track.get('album', {}).get('images', [{}])[0].get('url'),
                        duration_ms=track.get('duration_ms'),
                        popularity=track.get('popularity')
                    )
    except Exception as e:
        print(f"Error fetching Spotify track info: {str(e)}")
        
    return result