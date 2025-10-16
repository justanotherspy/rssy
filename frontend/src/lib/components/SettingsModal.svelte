<script lang="ts">
  import { showSettingsModal, posts } from '$lib/stores';
  import { postsApi } from '$lib/api';
  import { X, Trash2 } from 'lucide-svelte';

  let refreshInterval = 10;
  let loading = false;
  let error = '';

  function close() {
    showSettingsModal.set(false);
    error = '';
  }

  async function handleDeleteAllPosts() {
    if (!confirm('Are you sure you want to delete all posts? This cannot be undone.')) {
      return;
    }

    loading = true;
    error = '';
    try {
      await postsApi.deleteAll();
      posts.set([]);
      alert('All posts deleted successfully');
      close();
    } catch (err: any) {
      error = err.response?.data?.error || 'Failed to delete posts';
    } finally {
      loading = false;
    }
  }

  function handleSaveSettings() {
    close();
  }
</script>

{#if $showSettingsModal}
  <div class="modal-overlay" on:click={close} on:keydown={(e) => e.key === 'Escape' && close()} role="presentation">
    <div class="modal" on:click|stopPropagation on:keydown|stopPropagation>
      <div class="modal-header">
        <h2>Settings</h2>
        <button class="close-btn" on:click={close}>
          <X size={24} />
        </button>
      </div>

      <div class="modal-content">
        {#if error}
          <div class="error">{error}</div>
        {/if}

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
            Done
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

  .form-group {
    margin-bottom: 2rem;
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
    font-size: 1rem;
  }

  .danger-zone p {
    margin: 0 0 1rem 0;
    color: #666;
    font-size: 0.9rem;
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
    transition: all 0.2s;
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
    transition: all 0.2s;
  }

  .btn-primary {
    background: #0066cc;
    color: #fff;
  }

  .btn-primary:hover {
    background: #0052a3;
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
