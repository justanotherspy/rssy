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
