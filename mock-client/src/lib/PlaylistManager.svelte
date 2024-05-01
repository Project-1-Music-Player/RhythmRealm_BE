<script lang="ts">
    import { onMount } from 'svelte';
    import type { User } from 'firebase/auth';
    import { auth } from '../firebase'; // Ensure this path is correct
  
    let user: User | null = null;
    let playlists: any[] = [];
    let error: Error | null = null;
    let newPlaylistName: string = '';
    let newPlaylistDescription: string = '';
    let selectedPlaylistId: string = '';
    let songToAddId: string = '';
  
    auth.onAuthStateChanged((u) => {
      user = u;
      if (user) {
        fetchPlaylists();
        console.log("User playlist fetched")
      }
    });
  
    async function getIdToken(): Promise<string> {
      if (user) {
        return await user.getIdToken();
      } else {
        throw new Error('No user is signed in.');
      }
    }
  
    async function fetchPlaylists() {
      try {
        const idToken = await getIdToken();
        const response = await fetch('http://localhost:3000/playlists', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${idToken}`,
          },
        });
  
        if (response.ok) {
          playlists = await response.json();
          console.log(playlists);
        } else {
          console.error('Failed to fetch playlists');
        }
      } catch (err) {
        console.error('Failed to fetch playlists');
      }
    }
  
    async function addPlaylist() {
      try {
        const idToken = await getIdToken();
        const response = await fetch('http://localhost:3000/playlists', {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${idToken}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ name: newPlaylistName, description: newPlaylistDescription }),
        });
  
        if (response.ok) {
          newPlaylistName = '';
          newPlaylistDescription = '';
          fetchPlaylists();
        } else {
          throw new Error('Failed to create playlist');
        }
      } catch (err) {
        error = err;
      }
    }
  
    async function addSongToPlaylist() {
      try {
        const idToken = await getIdToken();
        const response = await fetch(`http://localhost:3000/playlists/${selectedPlaylistId}/songs`, {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${idToken}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ song_id: songToAddId }),
        });
  
        if (response.ok) {
          songToAddId = '';
          fetchPlaylists();
        } else {
          throw new Error('Failed to add song to playlist');
        }
      } catch (err) {
        error = err;
      }
    }
  
    async function removeSongFromPlaylist(playlistId: string, songId: string) {
      try {
        const idToken = await getIdToken();
        await fetch(`http://localhost:3000/playlists/${playlistId}/songs/${songId}`, {
          method: 'DELETE',
          headers: {
            'Authorization': `Bearer ${idToken}`,
          },
        });
  
        fetchPlaylists();
      } catch (err) {
        error = err;
      }
    }
  
    onMount(() => {
      if (user) {
        fetchPlaylists();
      }
    });
  </script>
  
  <main>
    {#if user}
      <h2>Your Playlists</h2>
      <div>
        <input type="text" bind:value={newPlaylistName} placeholder="New Playlist Name" />
        <input type="text" bind:value={newPlaylistDescription} placeholder="New Playlist Description" />
        <button on:click={addPlaylist}>Create New Playlist</button>
      </div>
  
      {#each playlists as playlist}
        <div class="playlist">
          <h3>{playlist.name}</h3>
          <p>{playlist.description}</p> <!-- TODO: Add insert and remove song from the playlist -->



          <!-- <select bind:value={selectedPlaylistId}>
            <option value="" disabled>Select a playlist</option>
            {#each playlists as p}
              <option value={p.playlist_id}>{p.name}</option>
            {/each}
          </select>
          <input type="text" bind:value={songToAddId} placeholder="Song ID to Add" />
          <button on:click={addSongToPlaylist} disabled={!selectedPlaylistId}>Add Song</button>
          <div>
            <h4>Songs in Playlist</h4>
            {#each playlist.songs as song}
              <div class="song">
                <p>{song.title}</p>
                <button on:click={() => removeSongFromPlaylist(playlist.playlist_id, song.song_id)}>Remove</button>
              </div>
            {/each}
          </div> -->
        </div>
      {/each}
    {:else}
      <p>Please sign in to manage your playlists.</p>
    {/if}
    {#if error}
      <p class="error">{error.message}</p>
    {/if}
  </main>
  
  <style>
    .error {
      color: red;
    }
    /* Add more styles as needed */
  </style>