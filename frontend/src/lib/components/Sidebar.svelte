<script lang="ts">
  import { feeds, selectedFeedId, showAddFeedModal, showSettingsModal, posts } from '$lib/stores';
  import { Plus, Settings } from 'lucide-svelte';
  import { postsApi } from '$lib/api';

  let loading = false;

  async function selectFeed(feedId: number | null) {
    selectedFeedId.set(feedId);
    loading = true;

    try {
      if (feedId === null) {
        const response = await postsApi.getAll();
        posts.set(response.data.data);
      } else {
        const response = await postsApi.getByFeed(feedId);
        posts.set(response.data.data);
      }
    } catch (err) {
      console.error('Error fetching posts:', err);
    } finally {
      loading = false;
    }
  }

  function openAddModal() {
    showAddFeedModal.set(true);
  }

  function openSettingsModal() {
    showSettingsModal.set(true);
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <h1>RSSY</h1>
    <div class="actions">
      <button on:click={openAddModal} title="Add feed" class="action-btn">
        <Plus size={20} />
      </button>
      <button on:click={openSettingsModal} title="Settings" class="action-btn">
        <Settings size={20} />
      </button>
    </div>
  </div>

  <nav class="feed-list">
    <button
      class="feed-item"
      class:active={$selectedFeedId === null}
      on:click={() => selectFeed(null)}
    >
      <span class="feed-name">#all</span>
    </button>

    {#each $feeds as feed (feed.id)}
      <button
        class="feed-item"
        class:active={$selectedFeedId === feed.id}
        on:click={() => selectFeed(feed.id)}
      >
        <span class="feed-name">#{feed.name}</span>
      </button>
    {/each}
  </nav>
</aside>

<style>
  .sidebar {
    width: 250px;
    height: 100vh;
    background: #1a1a1a;
    color: #fff;
    display: flex;
    flex-direction: column;
    border-right: 1px solid #333;
  }

  .sidebar-header {
    padding: 1.5rem 1rem;
    border-bottom: 1px solid #333;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .sidebar-header h1 {
    font-size: 1.5rem;
    margin: 0;
    font-weight: 700;
    letter-spacing: -0.5px;
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .action-btn {
    background: transparent;
    border: none;
    color: #fff;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background 0.2s;
  }

  .action-btn:hover {
    background: #333;
  }

  .feed-list {
    flex: 1;
    overflow-y: auto;
    padding: 0.5rem 0;
  }

  .feed-item {
    width: 100%;
    padding: 0.75rem 1rem;
    background: transparent;
    border: none;
    color: #ccc;
    text-align: left;
    cursor: pointer;
    font-size: 1rem;
    transition: background 0.2s, color 0.2s;
  }

  .feed-item:hover {
    background: #252525;
    color: #fff;
  }

  .feed-item.active {
    background: #333;
    color: #fff;
    font-weight: 600;
  }

  .feed-name {
    display: block;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
