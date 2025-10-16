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
          category: category.trim() || undefined,
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

  function toggleTab(isRed: boolean) {
    isReddit = isRed;
    error = '';
  }
</script>

{#if $showAddFeedModal}
  <div class="modal-overlay" on:click={close} on:keydown={(e) => e.key === 'Escape' && close()} role="presentation">
    <div class="modal" on:click|stopPropagation on:keydown|stopPropagation>
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
            on:click={() => toggleTab(false)}
          >
            RSS Feed
          </button>
          <button
            class:active={isReddit}
            on:click={() => toggleTab(true)}
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
    transition: background 0.2s;
  }

  .close-btn:hover {
    background: #f5f5f5;
    border-radius: 4px;
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
    transition: all 0.2s;
  }

  .tab-buttons button:hover {
    background: #eee;
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
    transition: border-color 0.2s;
  }

  .form-group input:focus {
    outline: none;
    border-color: #0066cc;
    box-shadow: 0 0 0 3px rgba(0, 102, 204, 0.1);
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
    transition: all 0.2s;
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
    border: 1px solid #fcc;
  }
</style>
