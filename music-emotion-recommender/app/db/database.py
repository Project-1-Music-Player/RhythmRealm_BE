import sqlite3
from contextlib import contextmanager
import pandas as pd
import json
from sklearn.preprocessing import MinMaxScaler
import os

DB_PATH = "music_recommendations.db"

class DatabasePool:
    def __init__(self):
        self.db_path = DB_PATH

    @contextmanager
    def get_connection(self):
        conn = sqlite3.connect(self.db_path, timeout=30.0)
        try:
            yield conn
        finally:
            conn.close()

db_pool = DatabasePool()

def safe_eval_list(x):
    if isinstance(x, str):
        try:
            cleaned = x.strip('[]').replace("'", "").split(',')
            return [item.strip() for item in cleaned if item.strip()]
        except:
            return []
    return x if isinstance(x, list) else []

def init_db():
    # Check if database already exists and is initialized
    if os.path.exists(DB_PATH):
        try:
            with db_pool.get_connection() as conn:
                cursor = conn.cursor()
                cursor.execute("SELECT COUNT(*) FROM songs")
                count = cursor.fetchone()[0]
                if count > 0:
                    print("Database already initialized")
                    return
        except:
            pass

    try:
        # Load the dataset
        df = pd.read_csv("muse_v3.csv")
        
        # Drop rows with missing required values
        df = df.dropna(subset=['track', 'artist'])
        
        # Process tags
        df['seeds'] = df['seeds'].apply(safe_eval_list)
        
        # Fill NaN values in emotion columns
        emotion_columns = ['valence_tags', 'arousal_tags', 'dominance_tags']
        df[emotion_columns] = df[emotion_columns].fillna(df[emotion_columns].mean())
        
        # Normalize emotional features
        scaler = MinMaxScaler()
        emotion_features = scaler.fit_transform(df[emotion_columns])

        with db_pool.get_connection() as conn:
            cursor = conn.cursor()
            
            # Drop existing tables
            cursor.execute('DROP TABLE IF EXISTS song_tags')
            cursor.execute('DROP TABLE IF EXISTS songs')
            
            # Create tables and indices
            cursor.execute('''
            CREATE TABLE songs (
                id INTEGER PRIMARY KEY,
                track TEXT NOT NULL,
                artist TEXT NOT NULL,
                genre TEXT,
                spotify_id TEXT,
                tags TEXT NOT NULL,
                valence REAL,
                arousal REAL,
                dominance REAL,
                normalized_valence REAL,
                normalized_arousal REAL,
                normalized_dominance REAL
            )
            ''')

            cursor.execute('''
            CREATE TABLE song_tags (
                song_id INTEGER,
                tag TEXT NOT NULL,
                FOREIGN KEY (song_id) REFERENCES songs (id)
            )
            ''')

            cursor.execute('CREATE INDEX idx_song_tags_tag ON song_tags(tag)')
            cursor.execute('CREATE INDEX idx_songs_spotify_id ON songs(spotify_id)')

            # Insert data in batches
            batch_size = 1000
            for i in range(0, len(df), batch_size):
                batch = df.iloc[i:i+batch_size]
                
                for _, row in batch.iterrows():
                    try:
                        idx = row.name
                        cursor.execute('''
                        INSERT INTO songs (
                            track, artist, genre, spotify_id, tags,
                            valence, arousal, dominance,
                            normalized_valence, normalized_arousal, normalized_dominance
                        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
                        ''', (
                            str(row['track']), str(row['artist']),
                            str(row['genre']) if pd.notna(row['genre']) else None,
                            str(row['spotify_id']) if pd.notna(row['spotify_id']) else None,
                            json.dumps(row['seeds']),
                            float(row['valence_tags']), float(row['arousal_tags']), float(row['dominance_tags']),
                            float(emotion_features[idx][0]), float(emotion_features[idx][1]), float(emotion_features[idx][2])
                        ))
                        
                        song_id = cursor.lastrowid
                        for tag in row['seeds']:
                            if tag and isinstance(tag, str):
                                cursor.execute('INSERT INTO song_tags (song_id, tag) VALUES (?, ?)',
                                             (song_id, tag.lower().strip()))
                    except Exception as e:
                        print(f"Error inserting row {idx}: {str(e)}")
                        continue
                
                conn.commit()  # Commit after each batch

            print("Database initialized successfully")

    except Exception as e:
        print(f"Error during database initialization: {str(e)}")
        raise