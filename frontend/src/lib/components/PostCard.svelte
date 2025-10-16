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
