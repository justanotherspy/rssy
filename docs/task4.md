# Task 4: Svelte Frontend Development

## Overview
Build the complete SvelteKit frontend application with a two-panel layout, feed management interface, post display cards, and all required modals for adding/editing feeds and managing settings.

## Goals
- Create responsive two-panel layout (sidebar + main content)
- Build feed list sidebar with active selection
- Implement post card components with proper formatting
- Create modals for adding/editing feeds
- Add Reddit quick-add functionality
- Implement settings modal
- Connect all components to backend API
- Style application for clean, readable RSS reader experience

## Prerequisites
- Task 1 completed (SvelteKit project initialized)
- Task 2 completed (data models defined)
- Task 3 completed (backend API running)
- Frontend dependencies installed (axios, date-fns, lucide-svelte)

## Application Structure

### Page Layout
```
┌────────────────────────────────────────────┐
│           RSSY Header                      │
├───────────┬────────────────────────────────┤
│  Sidebar  │      Main Content Area         │
│           │                                │
│  [+] [⚙]  │   ┌──────────────────────┐     │
│           │   │   Post Card          │     │
│  #feed1   │   │   - Image            │     │
│  #feed2*  │   │   - Title            │     │
│  #feed3   │   │   - Content          │     │
│           │   └──────────────────────┘     │
│           │   ┌──────────────────────┐     │
│           │   │   Post Card          │     │
│           │   └──────────────────────┘     │
└───────────┴────────────────────────────────┘

* = active/selected feed
```

## Detailed Steps

### 1. Set Up API Client

**Step 1.1: Create API service**

Create [frontend/src/lib/api.ts](frontend/src/lib/api.ts):

```typescript
import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Feed {
  id: number;
  name: string;
  url: string;
  category: string;
  site_url: string;
  description: string;
  is_active: boolean;
  last_fetched_at: string | null;
  error_count: number;
  last_error: string;
  created_at: string;
  updated_at: string;
}

export interface Post {
  id: number;
  feed_id: number;
  title: string;
  link: string;
  description: string;
  content: string;
  author: string;
  published_at: string | null;
  image_url: string;
  guid: string;
  is_read: boolean;
  created_at: string;
  updated_at: string;
  feed_name?: string;
}

export interface CreateFeedRequest {
  name: string;
  url: string;
  category?: string;
  site_url?: string;
  description?: string;
}

export interface UpdateFeedRequest {
  name?: string;
  url?: string;
  category?: string;
  site_url?: string;
  description?: string;
  is_active?: boolean;
}

// Feed API
export const feedsApi = {
  getAll: () => api.get<{ data: Feed[] }>('/api/feeds'),
  getById: (id: number) => api.get<{ data: Feed }>(`/api/feeds/${id}`),
  create: (feed: CreateFeedRequest) => api.post<{ data: Feed }>('/api/feeds', feed),
  update: (id: number, feed: UpdateFeedRequest) => api.put<{ data: Feed }>(`/api/feeds/${id}`, feed),
  delete: (id: number) => api.delete(`/api/feeds/${id}`),
  createReddit: (subreddit: string) => api.post<{ data: Feed }>('/api/feeds/reddit', { subreddit }),
  refresh: (id: number) => api.post(`/api/feeds/${id}/refresh`),
};

// Post API
export const postsApi = {
  getAll: (limit = 50, offset = 0) =>
    api.get<{ data: Post[] }>(`/api/posts?limit=${limit}&offset=${offset}`),
  getByFeed: (feedId: number, limit = 50, offset = 0) =>
    api.get<{ data: Post[] }>(`/api/posts/feed/${feedId}?limit=${limit}&offset=${offset}`),
  markRead: (id: number, isRead: boolean) =>
    api.patch(`/api/posts/${id}/read`, { is_read: isRead }),
  deleteAll: () => api.delete('/api/posts'),
};

export default api;
```

**Step 1.2: Create environment config**

Create [frontend/.env](frontend/.env):
```
VITE_API_URL=http://localhost:8080
```

### 2. Create Shared Stores

**Step 2.1: Create stores for state management**

Create [frontend/src/lib/stores.ts](frontend/src/lib/stores.ts):

```typescript
import { writable } from 'svelte/store';
import type { Feed, Post } from './api';

export const feeds = writable<Feed[]>([]);
export const posts = writable<Post[]>([]);
export const selectedFeedId = writable<number | null>(null);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

// Modal states
export const showAddFeedModal = writable<boolean>(false);
export const showEditFeedModal = writable<boolean>(false);
export const showSettingsModal = writable<boolean>(false);
export const editingFeed = writable<Feed | null>(null);
```

### 3. Create Component Structure

**Step 3.1: Create Sidebar component**

Create [frontend/src/lib/components/Sidebar.svelte](frontend/src/lib/components/Sidebar.svelte):

```svelte
<script lang="ts">
  import { feeds, selectedFeedId, showAddFeedModal, showEditFeedModal, showSettingsModal } from '$lib/stores';
  import { Plus, Settings } from 'lucide-svelte';
  import { feedsApi, postsApi } from '$lib/api';

  async function selectFeed(feedId: number | null) {
    selectedFeedId.set(feedId);

    try {
      if (feedId === null) {
        const response = await postsApi.getAll();
        // Update posts store (will implement later)
      } else {
        const response = await postsApi.getByFeed(feedId);
        // Update posts store
      }
    } catch (err) {
      console.error('Error fetching posts:', err);
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
      <button on:click={openAddModal} title="Add feed">
        <Plus size={20} />
      </button>
      <button on:click={openSettingsModal} title="Settings">
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
  }

  .actions {
    display: flex;
    gap: 0.5rem;
  }

  .actions button {
    background: transparent;
    border: none;
    color: #fff;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .actions button:hover {
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
```

**Step 3.2: Create PostCard component**

Create [frontend/src/lib/components/PostCard.svelte](frontend/src/lib/components/PostCard.svelte):

```svelte
<script lang="ts">
  import type { Post } from '$lib/api';
  import { formatDistanceToNow } from 'date-fns';
  import { ExternalLink } from 'lucide-svelte';

  export let post: Post;

  function formatDate(dateString: string | null): string {
    if (!dateString) return '';
    try {
      return formatDistanceToNow(new Date(dateString), { addSuffix: true });
    } catch {
      return '';
    }
  }

  function openLink() {
    window.open(post.link, '_blank');
  }
</script>

<article class="post-card" on:click={openLink}>
  {#if post.image_url}
    <div class="post-image">
      <img src={post.image_url} alt={post.title} />
    </div>
  {/if}

  <div class="post-content">
    <div class="post-meta">
      <span class="feed-name">{post.feed_name || 'Unknown Feed'}</span>
      {#if post.published_at}
        <span class="post-date">{formatDate(post.published_at)}</span>
      {/if}
    </div>

    <h2 class="post-title">
      {post.title}
      <ExternalLink size={16} class="external-icon" />
    </h2>

    {#if post.author}
      <p class="post-author">By {post.author}</p>
    {/if}

    {#if post.description}
      <div class="post-description">
        {@html post.description}
      </div>
    {/if}
  </div>
</article>

<style>
  .post-card {
    background: #fff;
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    overflow: hidden;
    cursor: pointer;
    transition: box-shadow 0.2s, transform 0.2s;
  }

  .post-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    transform: translateY(-2px);
  }

  .post-image {
    width: 100%;
    height: auto;
    max-height: 400px;
    overflow: hidden;
    background: #f5f5f5;
  }

  .post-image img {
    width: 100%;
    height: auto;
    display: block;
    object-fit: cover;
  }

  .post-content {
    padding: 1.5rem;
  }

  .post-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 0.75rem;
    font-size: 0.875rem;
    color: #666;
  }

  .feed-name {
    font-weight: 600;
    color: #0066cc;
  }

  .post-date {
    color: #999;
  }

  .post-title {
    font-size: 1.5rem;
    font-weight: 700;
    margin: 0 0 0.5rem 0;
    color: #1a1a1a;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .post-author {
    font-size: 0.9rem;
    color: #666;
    margin: 0 0 1rem 0;
    font-style: italic;
  }

  .post-description {
    color: #333;
    line-height: 1.6;
    text-align: justify;
  }

  .post-description :global(p) {
    margin: 0 0 1rem 0;
  }

  .post-description :global(img) {
    max-width: 100%;
    height: auto;
    margin: 1rem 0;
  }

  :global(.external-icon) {
    opacity: 0.5;
  }
</style>
```

**Step 3.3: Create AddFeedModal component**

Create [frontend/src/lib/components/AddFeedModal.svelte](frontend/src/lib/components/AddFeedModal.svelte):

```svelte
<script lang="ts">
  import { showAddFeedModal, feeds } from '$lib/stores';
  import { feedsApi } from '$lib/api';
  import { X } from 'lucide-svelte';

  let name = '';
  let url = '';
  let category = '';
  let isReddit = false;
  let subreddit = '';
  let loading = false;
  let error = '';

  function close() {
    showAddFeedModal.set(false);
    reset();
  }

  function reset() {
    name = '';
    url = '';
    category = '';
    isReddit = false;
    subreddit = '';
    error = '';
  }

  async function handleSubmit() {
    error = '';
    loading = true;

    try {
      if (isReddit) {
        if (!subreddit.trim()) {
          error = 'Subreddit name is required';
          loading = false;
          return;
        }
        await feedsApi.createReddit(subreddit.trim());
      } else {
        if (!name.trim() || !url.trim()) {
          error = 'Name and URL are required';
          loading = false;
          return;
        }
        await feedsApi.create({
          name: name.trim(),
          url: url.trim(),
          category: category.trim(),
        });
      }

      // Refresh feeds list
      const response = await feedsApi.getAll();
      feeds.set(response.data.data);

      close();
    } catch (err: any) {
      error = err.response?.data?.error || 'Failed to add feed';
    } finally {
      loading = false;
    }
  }
</script>

{#if $showAddFeedModal}
  <div class="modal-overlay" on:click={close}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Add New Feed</h2>
        <button class="close-btn" on:click={close}>
          <X size={24} />
        </button>
      </div>

      <div class="modal-content">
        <div class="tab-buttons">
          <button
            class:active={!isReddit}
            on:click={() => (isReddit = false)}
          >
            RSS Feed
          </button>
          <button
            class:active={isReddit}
            on:click={() => (isReddit = true)}
          >
            Reddit
          </button>
        </div>

        {#if error}
          <div class="error">{error}</div>
        {/if}

        <form on:submit|preventDefault={handleSubmit}>
          {#if isReddit}
            <div class="form-group">
              <label for="subreddit">Subreddit Name</label>
              <input
                id="subreddit"
                type="text"
                bind:value={subreddit}
                placeholder="programming"
                required
              />
              <small>Enter subreddit name without /r/</small>
            </div>
          {:else}
            <div class="form-group">
              <label for="name">Feed Name</label>
              <input
                id="name"
                type="text"
                bind:value={name}
                placeholder="My Favorite Blog"
                required
              />
            </div>

            <div class="form-group">
              <label for="url">Feed URL</label>
              <input
                id="url"
                type="url"
                bind:value={url}
                placeholder="https://example.com/feed.xml"
                required
              />
            </div>

            <div class="form-group">
              <label for="category">Category (optional)</label>
              <input
                id="category"
                type="text"
                bind:value={category}
                placeholder="Tech, News, etc."
              />
            </div>
          {/if}

          <div class="form-actions">
            <button type="button" class="btn-secondary" on:click={close}>
              Cancel
            </button>
            <button type="submit" class="btn-primary" disabled={loading}>
              {loading ? 'Adding...' : 'Add Feed'}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: #fff;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
    max-height: 90vh;
    overflow-y: auto;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e0e0e0;
  }

  .modal-header h2 {
    margin: 0;
    font-size: 1.5rem;
  }

  .close-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .modal-content {
    padding: 1.5rem;
  }

  .tab-buttons {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1.5rem;
  }

  .tab-buttons button {
    flex: 1;
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    background: #f5f5f5;
    cursor: pointer;
    border-radius: 4px;
    font-weight: 500;
  }

  .tab-buttons button.active {
    background: #0066cc;
    color: #fff;
    border-color: #0066cc;
  }

  .form-group {
    margin-bottom: 1.5rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
    color: #333;
  }

  .form-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    font-size: 1rem;
  }

  .form-group small {
    display: block;
    margin-top: 0.25rem;
    color: #666;
    font-size: 0.875rem;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    font-weight: 600;
  }

  .btn-primary {
    background: #0066cc;
    color: #fff;
  }

  .btn-primary:hover:not(:disabled) {
    background: #0052a3;
  }

  .btn-primary:disabled {
    background: #ccc;
    cursor: not-allowed;
  }

  .btn-secondary {
    background: #f5f5f5;
    color: #333;
  }

  .btn-secondary:hover {
    background: #e0e0e0;
  }

  .error {
    background: #fee;
    color: #c33;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }
</style>
```

**Step 3.4: Create SettingsModal component**

Create [frontend/src/lib/components/SettingsModal.svelte](frontend/src/lib/components/SettingsModal.svelte):

```svelte
<script lang="ts">
  import { showSettingsModal } from '$lib/stores';
  import { postsApi } from '$lib/api';
  import { X, Trash2 } from 'lucide-svelte';

  let refreshInterval = 10;
  let loading = false;

  function close() {
    showSettingsModal.set(false);
  }

  async function handleDeleteAllPosts() {
    if (!confirm('Are you sure you want to delete all posts? This cannot be undone.')) {
      return;
    }

    loading = true;
    try {
      await postsApi.deleteAll();
      alert('All posts deleted successfully');
      // Optionally refresh posts
    } catch (err) {
      alert('Failed to delete posts');
    } finally {
      loading = false;
    }
  }

  function handleSaveSettings() {
    // TODO: Implement settings save to backend
    alert('Settings saved (not yet implemented)');
    close();
  }
</script>

{#if $showSettingsModal}
  <div class="modal-overlay" on:click={close}>
    <div class="modal" on:click|stopPropagation>
      <div class="modal-header">
        <h2>Settings</h2>
        <button class="close-btn" on:click={close}>
          <X size={24} />
        </button>
      </div>

      <div class="modal-content">
        <div class="form-group">
          <label for="refresh-interval">Feed Refresh Interval (minutes)</label>
          <input
            id="refresh-interval"
            type="number"
            bind:value={refreshInterval}
            min="1"
            max="1440"
          />
          <small>How often to check feeds for new posts</small>
        </div>

        <div class="danger-zone">
          <h3>Danger Zone</h3>
          <p>Permanently delete all posts from the database.</p>
          <button
            class="btn-danger"
            on:click={handleDeleteAllPosts}
            disabled={loading}
          >
            <Trash2 size={16} />
            {loading ? 'Deleting...' : 'Delete All Posts'}
          </button>
        </div>

        <div class="form-actions">
          <button type="button" class="btn-secondary" on:click={close}>
            Cancel
          </button>
          <button type="button" class="btn-primary" on:click={handleSaveSettings}>
            Save Settings
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

<style>
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }

  .modal {
    background: #fff;
    border-radius: 8px;
    width: 90%;
    max-width: 500px;
  }

  .modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.5rem;
    border-bottom: 1px solid #e0e0e0;
  }

  .modal-header h2 {
    margin: 0;
  }

  .close-btn {
    background: transparent;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
  }

  .modal-content {
    padding: 1.5rem;
  }

  .form-group {
    margin-bottom: 2rem;
  }

  .form-group label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 600;
  }

  .form-group input {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    font-size: 1rem;
  }

  .form-group small {
    display: block;
    margin-top: 0.25rem;
    color: #666;
    font-size: 0.875rem;
  }

  .danger-zone {
    background: #fee;
    border: 1px solid #fcc;
    border-radius: 4px;
    padding: 1rem;
    margin-bottom: 2rem;
  }

  .danger-zone h3 {
    margin: 0 0 0.5rem 0;
    color: #c33;
  }

  .danger-zone p {
    margin: 0 0 1rem 0;
    color: #666;
  }

  .btn-danger {
    background: #c33;
    color: #fff;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 4px;
    cursor: pointer;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .btn-danger:hover:not(:disabled) {
    background: #a22;
  }

  .btn-danger:disabled {
    background: #ccc;
    cursor: not-allowed;
  }

  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.75rem 1.5rem;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
    font-weight: 600;
  }

  .btn-primary {
    background: #0066cc;
    color: #fff;
  }

  .btn-secondary {
    background: #f5f5f5;
    color: #333;
  }
</style>
```

### 4. Create Main Page

**Step 4.1: Update main route**

Update [frontend/src/routes/+page.svelte](frontend/src/routes/+page.svelte):

```svelte
<script lang="ts">
  import { onMount } from 'svelte';
  import { feeds, posts, selectedFeedId, loading } from '$lib/stores';
  import { feedsApi, postsApi } from '$lib/api';
  import Sidebar from '$lib/components/Sidebar.svelte';
  import PostCard from '$lib/components/PostCard.svelte';
  import AddFeedModal from '$lib/components/AddFeedModal.svelte';
  import SettingsModal from '$lib/components/SettingsModal.svelte';

  onMount(async () => {
    loading.set(true);
    try {
      // Load feeds
      const feedsResponse = await feedsApi.getAll();
      feeds.set(feedsResponse.data.data);

      // Load all posts initially
      const postsResponse = await postsApi.getAll();
      posts.set(postsResponse.data.data);
    } catch (err) {
      console.error('Error loading data:', err);
    } finally {
      loading.set(false);
    }
  });

  // Update posts when selected feed changes
  $: if ($selectedFeedId !== null) {
    loadPostsForFeed($selectedFeedId);
  }

  async function loadPostsForFeed(feedId: number | null) {
    loading.set(true);
    try {
      if (feedId === null) {
        const response = await postsApi.getAll();
        posts.set(response.data.data);
      } else {
        const response = await postsApi.getByFeed(feedId);
        posts.set(response.data.data);
      }
    } catch (err) {
      console.error('Error loading posts:', err);
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
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
    background: #f5f5f5;
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
</style>
```

### 5. Configure SvelteKit

**Step 5.1: Update svelte.config.js**

Verify [frontend/svelte.config.js](frontend/svelte.config.js) is using the static adapter or auto adapter.

**Step 5.2: Create app.html**

Ensure [frontend/src/app.html](frontend/src/app.html) exists with proper structure.

### 6. Testing the Frontend

**Step 6.1: Start development servers**

Terminal 1 - Backend:
```bash
cd backend
go run ./cmd/api
```

Terminal 2 - Frontend:
```bash
cd frontend
npm run dev
```

**Step 6.2: Test functionality**
- Visit http://localhost:5173
- Verify feeds load in sidebar
- Click on different feeds
- Add new feed via + button
- Add Reddit feed
- Open settings modal
- Verify posts display correctly
- Test responsive layout

### 7. Build for Production

**Step 7.1: Build frontend**
```bash
cd frontend
npm run build
npm run preview  # Test production build
```

**Step 7.2: Build backend**
```bash
cd backend
make backend-build
./bin/api
```

## Success Criteria

- Application loads without errors
- Sidebar displays all feeds
- Clicking feeds updates main content area
- Add feed modal works for both RSS and Reddit
- Settings modal opens and functions
- Post cards display properly with images
- Links open in new tabs
- Responsive layout works on different screen sizes
- All API calls succeed
- Error handling displays appropriate messages

## Next Steps

After completing Task 4, you have a working MVP! Consider:
- Adding authentication
- Implementing feed categorization
- Adding search functionality
- Implementing infinite scroll
- Adding keyboard shortcuts
- Implementing read/unread tracking
- Adding dark mode
- Deploying to production

## Notes

- Use Svelte stores for state management
- All API calls should handle errors gracefully
- Images in posts should be responsive
- Consider lazy loading for post images
- The layout should be mobile-friendly
- Use semantic HTML for accessibility

## Agent & Hook Recommendations for This Task

### Recommended Agent Usage
- **Direct tool usage** for creating component files
- **General-purpose agent** for complex component logic
- Test in browser frequently during development

### File Save Hook
The file save hook for formatting will be helpful:
```bash
if [[ "$FILE_PATH" == *.svelte ]] || [[ "$FILE_PATH" == *.ts ]]; then
  cd frontend && npx prettier --write "$FILE_PATH"
fi
```

### Testing Strategy
- Manual testing in browser
- Test all user flows
- Verify API integration
- Check responsive design
- Test error states

## Troubleshooting

**CORS errors:**
- Ensure backend allows frontend origin
- Check browser console for specific errors
- Verify CORS middleware in backend

**Posts not loading:**
- Check browser network tab
- Verify API endpoints are correct
- Check backend logs for errors

**Styling issues:**
- Verify Svelte scoped styles
- Check for CSS conflicts
- Use browser dev tools

**Build errors:**
- Check TypeScript types
- Verify all imports
- Run `npm run check`

**API connection issues:**
- Verify VITE_API_URL in .env
- Ensure backend is running
- Check network requests in browser dev tools
