<script lang="ts">
  import svelteLogo from './assets/svelte.svg'
  import viteLogo from '/vite.svg'
  import Counter from './lib/Counter.svelte'
  import type { User } from "firebase/auth";
  import { auth, googleProvider, signInAnonymously, signInWithPopup, signOut } from './firebase.js'; // Ensure this path is correct
  import { onMount } from 'svelte';
  
  // Use the code from these follow imports to implement the front end.
  import Search from './lib/Search.svelte';
  import PlaylistManager from './lib/PlaylistManager.svelte';
  let title: string | Blob, album: string | Blob, releaseDate: string | Blob, genre: string | Blob, songFile: string | Blob, thumbnailFile: string | Blob;
  let user: User | null = null;
  let error: Error | null = null;
  auth.onAuthStateChanged(console.log)
  let songs: any[] = [];
  auth.onAuthStateChanged((u) => {
    user = u;
    if (user) {
      console.log('User is signed in:', user.uid);
      fetchSongs(); // Fetch songs after user logs in
    } else {
      console.log('No user is signed in.');
    }
  });

  async function getIdToken(): Promise<string> {
    if (user) {
      try {
        const token = await user.getIdToken();
        return token;
      } catch (error) {
        console.error('Failed to get ID token:', error);
        throw error;
      }
    } else {
      throw new Error('No user is signed in.');
    }
  }

  // New function to fetch songs
  async function fetchSongs() {
    try {
      const idToken = await getIdToken();
      const response = await fetch('http://localhost:3000/music', {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${idToken}`,
        }
      });

      if (response.ok) {
        songs = await response.json();
      } else {
        console.error('Failed to fetch songs');
      }
    } catch (error) {
      console.error('Error fetching songs:', error);
    }
  }
  function thumbnailUrl(songID: string) {
    return `http://localhost:3000/music/thumbnail/${songID}`;
  }
  function streamUrl(songID: string) {
    return `http://localhost:3000/music/stream/${songID}`;
  }
  async function handleUpload(event: Event) {
    event.preventDefault();
    const formData = new FormData();
    formData.append('title', title);
    formData.append('album', album);
    formData.append('releaseDate', releaseDate);
    formData.append('genre', genre);
    formData.append('song', songFile);
    formData.append('thumbnail', thumbnailFile);  

    const idToken = await getIdToken(); // Function to get ID token from Firebase

    const response = await fetch('http://localhost:3000/music/upload', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${idToken}`,
      },
      body: formData
    });

    if (response.ok) {
      console.log('Music uploaded successfully');
    } else {
      console.error('Upload failed');
    }
  }

  function handleSongFileChange(event: any) {
    songFile = event.target.files[0];
  }

  function handleThumbnailFileChange(event: any) {
    thumbnailFile = event.target.files[0];
  }

  function handleAnonymousSignIn() {
    signInAnonymously(auth)
      .then((result) => {
        user = result.user;
        console.log('Signed in anonymously as:', user.uid);
      })
      .catch((err) => {
        error = err;
        console.error('Anonymous sign-in failed:', error);
      });
  }

  async function handleGoogleSignIn() {
    try {
      const result = await signInWithPopup(auth, googleProvider);
      user = result.user;
      console.log('Signed in with Google as:', user.uid);

      const idToken = await user.getIdToken();

      const response = await fetch('http://localhost:3000/auth/google', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${idToken}`,
        },
        body: JSON.stringify({
          username: user.displayName || '',
          email: user.email || '',
          role: 'listener', // Default role
        }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      console.log('User upserted successfully');
    } catch (err) {
      error = err as Error | null;
      console.error('Google sign-in failed:', error);
    }
  }
  function handleSignOut() {
    auth.signOut()
      .then(() => {
        user = null;
        location.reload();
        console.log('User signed out.');
      })
      .catch((err) => {
        error = err;
        console.error('Sign out failed:', error);
      });
  }
</script>

<main>

  <div class="card">
    <Counter />
  </div>
  <div>
    <button on:click={handleGoogleSignIn}>Sign in with Google</button>
    <button on:click={handleAnonymousSignIn}>Sign in Anonymously</button>
  </div>
  <p>
    Check out <a href="https://github.com/sveltejs/kit#readme" target="_blank" rel="noreferrer">SvelteKit</a>, the official Svelte app framework powered by Vite!
  </p>
  <h1>Firebase Auth with Svelte</h1>
  {#if user}
    <p>Signed in as: {user.displayName} ({user.uid})</p>
    <button on:click={handleSignOut}>Sign Out</button>
    <form on:submit|preventDefault={handleUpload}>
      <input type="text" bind:value={title} placeholder="Title" required />
      <input type="text" bind:value={album} placeholder="Album" required />
      <input type="date" bind:value={releaseDate} placeholder="Release Date" required />
      <input type="text" bind:value={genre} placeholder="Genre" required />
      <input type="file" on:change={handleSongFileChange} required />
      <input type="file" on:change={handleThumbnailFileChange} required />
      <button type="submit">Upload</button>
    </form>
    <div>
      <h2>Your Songs</h2>
      {#each songs as song}
        <div class="song">
          <h3>{song.title}</h3>
          <img src={thumbnailUrl(song.song_id)} alt={song.title} />
          <p>Album: {song.album}</p>
          <p>Release Date: {song.releaseDate}</p>
          <p>Genre: {song.genre}</p>
          <audio controls>
            <source src={streamUrl(song.song_id)} type="audio/mpeg" />
            Your browser does not support the audio element.
          </audio>
        </div>
      {/each}
    </div>
  {:else}
    <p>No user is signed in.</p>
  {/if}
  {#if error}
    <p style="color: red">Error: {error.message}</p>
  {/if}
    <Search />
    <PlaylistManager />
</main>

<style>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.svelte:hover {
    filter: drop-shadow(0 0 2em #ff3e00aa);
  }
</style>
