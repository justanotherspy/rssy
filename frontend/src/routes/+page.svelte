<script lang="ts">
  import { onMount } from 'svelte';
  import { feeds, posts, selectedFeedId, loading, error } from '$lib/stores';
  import { feedsApi, postsApi } from '$lib/api';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import PostCard from '$lib/components/PostCard.svelte';
  import AddFeedModal from '$lib/components/AddFeedModal.svelte';
  import SettingsModal from '$lib/components/SettingsModal.svelte';

  onMount(async () => {
    loading.set(true);
    error.set(null);
    try {
      // Load feeds
      const feedsResponse = await feedsApi.getAll();
      feeds.set(feedsResponse.data.data);

      // Load all posts initially
      const postsResponse = await postsApi.getAll();
      posts.set(postsResponse.data.data);

      // Set initial selected feed to all
      selectedFeedId.set(null);
    } catch (err: any) {
      console.error('Error loading data:', err);
      error.set(err.response?.data?.error || 'Failed to load data');
    } finally {
      loading.set(false);
    }
  });

  // Update posts when selected feed changes
  $: if ($selectedFeedId !== undefined && $selectedFeedId !== null) {
    loadPostsForFeed($selectedFeedId);
  }

  async function loadPostsForFeed(feedId: number | null) {
    loading.set(true);
    error.set(null);
    try {
      if (feedId === null) {
        const response = await postsApi.getAll();
        posts.set(response.data.data);
      } else {
        const response = await postsApi.getByFeed(feedId);
        posts.set(response.data.data);
      }
    } catch (err: any) {
      console.error('Error loading posts:', err);
      error.set(err.response?.data?.error || 'Failed to load posts');
    } finally {
      loading.set(false);
    }
  }
</script>

<svelte:head>
  <title>RSSY - RSS Reader</title>
</svelte:head>

<div class="app-container">
  <Sidebar />

  <main class="main-content">
    {#if $error}
      <div class="error-banner">{$error}</div>
    {/if}

    {#if $loading}
      <div class="loading">Loading...</div>
    {:else if $posts.length === 0}
      <div class="empty-state">
        <h2>No posts yet</h2>
        <p>Add some feeds to get started!</p>
      </div>
    {:else}
      <div class="posts-container">
        {#each $posts as post (post.id)}
          <PostCard {post} />
        {/each}
      </div>
    {/if}
  </main>

  <AddFeedModal />
  <SettingsModal />
</div>

<style>
  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen,
      Ubuntu, Cantarell, sans-serif;
    background: #f5f5f5;
  }

  :global(html, body, #svelte) {
    height: 100%;
  }

  .app-container {
    display: flex;
    height: 100vh;
    overflow: hidden;
  }

  .main-content {
    flex: 1;
    overflow-y: auto;
    padding: 2rem;
  }

  .posts-container {
    max-width: 800px;
    margin: 0 auto;
  }

  .loading,
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: #666;
  }

  .empty-state h2 {
    margin: 0 0 0.5rem 0;
    font-size: 2rem;
    color: #333;
  }

  .empty-state p {
    margin: 0;
    font-size: 1.125rem;
  }

  .error-banner {
    background: #fee;
    color: #c33;
    padding: 1rem;
    border-radius: 4px;
    border: 1px solid #fcc;
    margin-bottom: 1rem;
  }
</style>
