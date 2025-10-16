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
