from fastapi import APIRouter, HTTPException, Query
from typing import List, Optional
import json
from sklearn.preprocessing import MinMaxScaler
from app.models.schemas import SongResponse
from app.db.database import db_pool
from app.services.spotify import get_spotify_track_infos

router = APIRouter()

@router.get("/tags", response_model=List[str])
async def get_all_tags():
    """Get all unique tags from the dataset"""
    try:
        with db_pool.get_connection() as conn:
            cursor = conn.cursor()
            cursor.execute('SELECT DISTINCT tag FROM song_tags ORDER BY tag')
            tags = [row[0] for row in cursor.fetchall()]
            return tags
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/recommend/{tag}", response_model=List[SongResponse])
async def recommend_by_tag(
    tag: str,
    page: int = Query(1, ge=1),
    page_size: int = Query(10, ge=1, le=50),
    valence: Optional[float] = None,
    arousal: Optional[float] = None,
    dominance: Optional[float] = None,
    include_spotify_info: bool = Query(True, description="Whether to include additional Spotify information")
):
    """
    Recommend songs based on a specific tag and optional emotion parameters with pagination
    - page: Page number (starts from 1)
    - page_size: Number of items per page (1-50)
    - valence: pleasure-displeasure (1-7)
    - arousal: excitement-calm (1-7)
    - dominance: dominance-submissiveness (1-7)
    """
    try:
        with db_pool.get_connection() as conn:
            cursor = conn.cursor()

            # Base query with emotion scoring
            query = '''
            WITH matched_songs AS (
                SELECT DISTINCT 
                    s.*,
                    CASE 
                        WHEN ? IS NOT NULL 
                        THEN (1 - ABS(s.normalized_valence - ?)) * 0.3 +
                             (1 - ABS(s.normalized_arousal - ?)) * 0.3 +
                             (1 - ABS(s.normalized_dominance - ?)) * 0.3
                        ELSE 1
                    END as emotion_score
                FROM songs s
                JOIN song_tags st ON s.id = st.song_id
                WHERE st.tag = ?
            )
            SELECT *
            FROM matched_songs
            ORDER BY emotion_score DESC
            LIMIT ? OFFSET ?
            '''

            # Normalize emotion values if provided
            use_emotions = all(v is not None for v in [valence, arousal, dominance])
            if use_emotions:
                scaler = MinMaxScaler()
                normalized_emotions = scaler.fit_transform([[valence, arousal, dominance]])[0]
            else:
                normalized_emotions = [None, None, None]

            cursor.execute(query, (
                1 if use_emotions else None,
                *normalized_emotions,
                tag.lower(),
                page_size,
                (page - 1) * page_size
            ))

            results = []
            spotify_ids = []
            
            for row in cursor.fetchall():
                spotify_id = row[4]
                if spotify_id:
                    spotify_ids.append(spotify_id)
                
                results.append({
                    'track': row[1],
                    'artist': row[2],
                    'genre': row[3],
                    'spotify_id': spotify_id,
                    'tags': json.loads(row[5]),
                    'score': float(row[12]),
                    'valence': float(row[6]),
                    'arousal': float(row[7]),
                    'dominance': float(row[8]),
                    'spotify_info': None  # Will be populated later
                })

            if results and include_spotify_info:
                # Fetch Spotify information in batch
                spotify_info_map = get_spotify_track_infos(spotify_ids)
                
                # Add Spotify information to results
                for result in results:
                    if result['spotify_id']:
                        result['spotify_info'] = spotify_info_map.get(result['spotify_id'])

            return results

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/recommend/multiple/", response_model=List[SongResponse])
async def recommend_by_multiple_tags(
    tags: str,
    page: int = Query(1, ge=1),
    page_size: int = Query(10, ge=1, le=50),
    valence: Optional[float] = None,
    arousal: Optional[float] = None,
    dominance: Optional[float] = None
):
    """
    Recommend songs based on multiple tags and optional emotion parameters with pagination
    - tags: comma-separated list of tags
    - page: Page number (starts from 1)
    - page_size: Number of items per page (1-50)
    - valence: pleasure-displeasure (1-7)
    - arousal: excitement-calm (1-7)
    - dominance: dominance-submissiveness (1-7)
    """
    try:
        tag_list = [tag.strip().lower() for tag in tags.split(',')]
        placeholders = ','.join('?' * len(tag_list))

        with db_pool.get_connection() as conn:
            cursor = conn.cursor()

            query = f'''
            WITH matched_songs AS (
                SELECT DISTINCT 
                    s.*,
                    COUNT(DISTINCT st.tag) as tag_matches,
                    CASE 
                        WHEN ? IS NOT NULL 
                        THEN (1 - ABS(s.normalized_valence - ?)) * 0.3 +
                             (1 - ABS(s.normalized_arousal - ?)) * 0.3 +
                             (1 - ABS(s.normalized_dominance - ?)) * 0.3
                        ELSE 1
                    END as emotion_score
                FROM songs s
                JOIN song_tags st ON s.id = st.song_id
                WHERE st.tag IN ({placeholders})
                GROUP BY s.id
            )
            SELECT *,
                (tag_matches * 0.4 + emotion_score * 0.6) as final_score
            FROM matched_songs
            ORDER BY final_score DESC
            LIMIT ? OFFSET ?
            '''

            # Normalize emotion values if provided
            use_emotions = all(v is not None for v in [valence, arousal, dominance])
            if use_emotions:
                scaler = MinMaxScaler()
                normalized_emotions = scaler.fit_transform([[valence, arousal, dominance]])[0]
            else:
                normalized_emotions = [None, None, None]

            cursor.execute(query, (
                1 if use_emotions else None,
                *normalized_emotions,
                *tag_list,
                page_size,
                (page - 1) * page_size
            ))

            results = []
            for row in cursor.fetchall():
                results.append({
                    'track': row[1],
                    'artist': row[2],
                    'genre': row[3],
                    'spotify_id': row[4],
                    'tags': json.loads(row[5]),
                    'score': float(row[14]),  # final_score
                    'valence': float(row[6]),
                    'arousal': float(row[7]),
                    'dominance': float(row[8])
                })

            if not results:
                raise HTTPException(status_code=404, detail=f"No songs found with tags: {tags}")

            return results

    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))