<!-- Sample code to search for songs -->
<script>
    import { onDestroy, onMount } from 'svelte';
    import { throttle } from 'lodash-es'; // Throttle function calls
  
    /**
     * @type {string | any[]}
     */
    let songs = [];
    let query = '';
    let page = 1;
    let loading = false;
  
    async function searchSongs() {
      loading = true;
      const response = await fetch(`http://localhost:3000/music/search?q=${encodeURIComponent(query)}&page=${page}&limit=10`);
      if (response.ok) {
        const newSongs = await response.json();
        songs = [...songs, ...newSongs];
        page += 1;
      } else {
        console.error('Failed to fetch songs');
      }
      loading = false;
    }
  
    function onScroll() {
      const threshold = 100; 
      const position = window.innerHeight + window.scrollY;
      const height = document.body.offsetHeight;
      if (position >= height - threshold && !loading) {
        searchSongs();
      }
    }
  
    const throttledOnScroll = throttle(onScroll, 200);
  
    onMount(() => {
      window.addEventListener('scroll', throttledOnScroll);
    });
  
    onDestroy(() => {
      window.removeEventListener('scroll', throttledOnScroll);
    });
  </script>
  
<input type="text" bind:value={query} placeholder="Search songs..." on:input="{() => { page = 1; songs = []; searchSongs(); }}" />
  
  {#if songs.length === 0 && !loading}
    <p>No songs found.</p>
  {:else}
    {#each songs as song}
      <div class="song">
        <h2>{song.title}</h2>
      </div>
    {/each}
  {/if}
  
  {#if loading}
    <p>Loading...</p>
  {/if}