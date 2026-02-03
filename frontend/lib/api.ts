// APIクライアント関数

const API_BASE_URL = 'http://localhost:8080/api/v1';

// 型定義
export interface Media {
  id: string;
  type: 'image' | 'video' | 'audio';
  title: string;
  description?: string;
  s3_key?: string;
  cloudfront_url?: string;
  youtube_url?: string;
  tags: Tag[];
  created_at: string;
  updated_at: string;
}

export interface Tag {
  id: string;
  name: string;
  type: 'all' | 'image' | 'audio' | 'video';
  created_at: string;
  updated_at: string;
}

export interface MediaListResponse {
  media: Media[];
  total?: number;
  offset?: number;
  limit?: number;
  has_more?: boolean;
}

export interface TagListResponse {
  tags: Tag[];
}

export interface ErrorResponse {
  error: string;
}

export interface MessageResponse {
  message: string;
}

// メディア関連API
export async function getMediaList(): Promise<Media[]> {
  const response = await fetch(`${API_BASE_URL}/media`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch media list');
  }
  const data: MediaListResponse = await response.json();
  return data.media;
}

export async function getMediaListWithPagination(
  offset: number = 0,
  limit: number = 20,
  title?: string,
  tagIds?: string[]
): Promise<MediaListResponse> {
  const params = new URLSearchParams();
  params.append('offset', offset.toString());
  params.append('limit', limit.toString());
  if (title) {
    params.append('title', title);
  }
  if (tagIds && tagIds.length > 0) {
    tagIds.forEach((tagId) => {
      params.append('tag_ids', tagId);
    });
  }

  const response = await fetch(
    `${API_BASE_URL}/media?${params.toString()}`
  );
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch media list');
  }
  return await response.json();
}

export async function getMedia(id: string): Promise<Media> {
  const response = await fetch(`${API_BASE_URL}/media/${id}`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch media');
  }
  return await response.json();
}

export async function uploadMedia(
  file: File,
  title: string,
  description?: string,
  tagIds?: string[]
): Promise<Media> {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('title', title);
  if (description) {
    formData.append('description', description);
  }
  if (tagIds && tagIds.length > 0) {
    tagIds.forEach((tagId) => {
      formData.append('tag_ids', tagId);
    });
  }

  const response = await fetch(`${API_BASE_URL}/media/upload`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to upload media');
  }
  return await response.json();
}

export async function createMediaWithYouTube(
  youtubeUrl: string,
  title: string,
  description?: string,
  tagIds?: string[]
): Promise<Media> {
  const response = await fetch(`${API_BASE_URL}/media/youtube`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      youtube_url: youtubeUrl,
      title,
      ...(description && description.trim() ? { description } : {}),
      tag_ids: tagIds || [],
    }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to create media with YouTube');
  }
  return await response.json();
}

export async function deleteMedia(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/media/${id}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to delete media');
  }
}

export async function associateTag(mediaId: string, tagId: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/media/${mediaId}/tags`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ tag_id: tagId }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to associate tag');
  }
}

export async function removeTag(mediaId: string, tagId: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/media/${mediaId}/tags/${tagId}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to remove tag');
  }
}

export async function getMediaByTag(tagId: string): Promise<Media[]> {
  const response = await fetch(`${API_BASE_URL}/tags/${tagId}/media`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch media by tag');
  }
  const data: MediaListResponse = await response.json();
  return data.media;
}

// タグ関連API
export async function getTagList(): Promise<Tag[]> {
  const response = await fetch(`${API_BASE_URL}/tags`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch tag list');
  }
  const data: TagListResponse = await response.json();
  return data.tags;
}

export async function getTag(id: string): Promise<Tag> {
  const response = await fetch(`${API_BASE_URL}/tags/${id}`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch tag');
  }
  return await response.json();
}

export async function createTag(name: string, type: 'all' | 'image' | 'audio' | 'video' = 'all'): Promise<Tag> {
  const response = await fetch(`${API_BASE_URL}/tags`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, type }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to create tag');
  }
  return await response.json();
}

export async function updateTag(id: string, name: string, type: 'all' | 'image' | 'audio' | 'video'): Promise<Tag> {
  const response = await fetch(`${API_BASE_URL}/tags/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ name, type }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to update tag');
  }
  return await response.json();
}

export async function deleteTag(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/tags/${id}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to delete tag');
  }
}

// TODO関連API
export interface Todo {
  id: string;
  title: string;
  description?: string;
  start_date?: string;
  end_date?: string;
  due_date?: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
}

export interface TodoListResponse {
  todos: Todo[];
  total?: number;
  offset?: number;
  limit?: number;
  has_more?: boolean;
}

export async function createTodo(
  title: string,
  description?: string,
  startDate?: string,
  endDate?: string,
  dueDate?: string
): Promise<Todo> {
  // 日付をRFC3339形式に変換（YYYY-MM-DD形式の場合は時刻を追加）
  const formatDate = (dateStr?: string): string | undefined => {
    if (!dateStr) return undefined;
    // 既にRFC3339形式の場合はそのまま
    if (dateStr.includes('T')) return dateStr;
    // YYYY-MM-DD形式の場合は時刻を追加
    return `${dateStr}T00:00:00Z`;
  };

  const response = await fetch(`${API_BASE_URL}/todos`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      title,
      description: description || undefined,
      start_date: formatDate(startDate),
      end_date: formatDate(endDate),
      due_date: formatDate(dueDate),
    }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to create todo');
  }
  return await response.json();
}

export async function getTodo(id: string): Promise<Todo> {
  const response = await fetch(`${API_BASE_URL}/todos/${id}`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch todo');
  }
  return await response.json();
}

export async function getTodoList(): Promise<Todo[]> {
  const response = await fetch(`${API_BASE_URL}/todos`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch todo list');
  }
  const data: TodoListResponse = await response.json();
  return data.todos;
}

export async function getTodosByDateRange(
  startDate: string,
  endDate: string
): Promise<Todo[]> {
  const response = await fetch(
    `${API_BASE_URL}/todos/date-range?start_date=${startDate}&end_date=${endDate}`
  );
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch todos by date range');
  }
  const data: TodoListResponse = await response.json();
  return data.todos;
}

export async function getTodosByDate(date: string): Promise<Todo[]> {
  const response = await fetch(`${API_BASE_URL}/todos/date?date=${date}`);
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch todos by date');
  }
  const data: TodoListResponse = await response.json();
  return data.todos;
}

export async function getTodosWithoutDueDate(
  offset: number = 0,
  limit: number = 20
): Promise<TodoListResponse> {
  const response = await fetch(
    `${API_BASE_URL}/todos/without-due-date?offset=${offset}&limit=${limit}`
  );
  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to fetch todos without due date');
  }
  return await response.json();
}

export async function updateTodo(
  id: string,
  title: string,
  description?: string,
  startDate?: string,
  endDate?: string,
  dueDate?: string,
  completed: boolean = false
): Promise<Todo> {
  // 日付をRFC3339形式に変換（YYYY-MM-DD形式の場合は時刻を追加）
  const formatDate = (dateStr?: string): string | undefined => {
    if (!dateStr) return undefined;
    // 既にRFC3339形式の場合はそのまま
    if (dateStr.includes('T')) return dateStr;
    // YYYY-MM-DD形式の場合は時刻を追加
    return `${dateStr}T00:00:00Z`;
  };

  const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      title,
      description: description || undefined,
      start_date: formatDate(startDate),
      end_date: formatDate(endDate),
      due_date: formatDate(dueDate),
      completed,
    }),
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to update todo');
  }
  return await response.json();
}

export async function deleteTodo(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/todos/${id}`, {
    method: 'DELETE',
  });

  if (!response.ok) {
    const error: ErrorResponse = await response.json();
    throw new Error(error.error || 'Failed to delete todo');
  }
}
