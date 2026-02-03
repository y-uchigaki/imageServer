'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import {
  getTagList,
  createTag,
  updateTag,
  deleteTag,
  getMediaByTag,
  type Tag,
  type Media,
} from '@/lib/api';
import { tagSchema } from '@/lib/validations';

// YouTube URLからIDを抽出する関数
function extractYouTubeId(url: string): string | null {
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/)([^&\n?#]+)/,
    /^([a-zA-Z0-9_-]{11})$/,
  ];

  for (const pattern of patterns) {
    const match = url.match(pattern);
    if (match) {
      return match[1];
    }
  }

  return null;
}

export default function TagsPage() {
  const router = useRouter();
  const [tagList, setTagList] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // 作成フォームの状態
  const [createName, setCreateName] = useState('');
  const [createType, setCreateType] = useState<'all' | 'image' | 'audio' | 'video'>('all');
  const [creating, setCreating] = useState(false);
  const [createErrors, setCreateErrors] = useState<Record<string, string>>({});

  // 編集フォームの状態
  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState('');
  const [editType, setEditType] = useState<'all' | 'image' | 'audio' | 'video'>('all');
  const [updating, setUpdating] = useState(false);
  const [editErrors, setEditErrors] = useState<Record<string, string>>({});

  // メディア表示の状態
  const [selectedTagId, setSelectedTagId] = useState<string | null>(null);
  const [tagMedia, setTagMedia] = useState<Media[]>([]);
  const [loadingMedia, setLoadingMedia] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const tags = await getTagList();
      setTagList(tags);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreateErrors({});
    setError(null);

    // バリデーション
    const result = tagSchema.safeParse({
      name: createName.trim(),
      type: createType,
    });

    if (!result.success) {
      const errors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        const path = issue.path.join('.');
        errors[path] = issue.message;
      });
      setCreateErrors(errors);
      return;
    }

    try {
      setCreating(true);
      setError(null);
      setSuccess(null);
      await createTag(createName.trim(), createType);
      setSuccess('タグの作成に成功しました');
      setCreateName('');
      setCreateType('all');
      setCreateErrors({});
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : '作成に失敗しました');
    } finally {
      setCreating(false);
    }
  };

  const handleStartEdit = (tag: Tag) => {
    setEditingId(tag.id);
    setEditName(tag.name);
    setEditType(tag.type);
  };

  const handleCancelEdit = () => {
    setEditingId(null);
    setEditName('');
    setEditType('all');
    setEditErrors({});
  };

  const handleUpdate = async (id: string) => {
    setEditErrors({});
    setError(null);

    // バリデーション
    const result = tagSchema.safeParse({
      name: editName.trim(),
      type: editType,
    });

    if (!result.success) {
      const errors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        const path = issue.path.join('.');
        errors[path] = issue.message;
      });
      setEditErrors(errors);
      return;
    }

    try {
      setUpdating(true);
      setError(null);
      setSuccess(null);
      await updateTag(id, editName.trim(), editType);
      setSuccess('タグの更新に成功しました');
      setEditingId(null);
      setEditName('');
      setEditType('all');
      setEditErrors({});
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新に失敗しました');
    } finally {
      setUpdating(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('このタグを削除しますか？')) {
      return;
    }

    try {
      setError(null);
      setSuccess(null);
      await deleteTag(id);
      setSuccess('タグの削除に成功しました');
      if (selectedTagId === id) {
        setSelectedTagId(null);
        setTagMedia([]);
      }
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : '削除に失敗しました');
    }
  };

  const handleViewMedia = async (tagId: string) => {
    try {
      setLoadingMedia(true);
      setError(null);
      const media = await getMediaByTag(tagId);
      setTagMedia(media);
      setSelectedTagId(tagId);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'メディアの取得に失敗しました');
    } finally {
      setLoadingMedia(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <button
            onClick={() => router.push('/')}
            className="text-blue-600 hover:text-blue-800 underline mb-4"
          >
            ← ホームに戻る
          </button>
          <h1 className="text-4xl font-bold text-gray-900">タグ管理</h1>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded mb-4">
            {success}
          </div>
        )}

        {/* 作成フォーム */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-2xl font-semibold mb-4 text-gray-800">
            タグを作成
          </h2>
          <form onSubmit={handleCreate} className="space-y-4">
            <div className="flex gap-4">
              <div className="flex-1">
                <input
                  type="text"
                  value={createName || ''}
                  onChange={(e) => {
                    setCreateName(e.target.value);
                    if (createErrors.name) {
                      setCreateErrors((prev) => ({ ...prev, name: '' }));
                    }
                  }}
                  placeholder="タグ名"
                  className={`w-full px-3 py-2 border rounded ${
                    createErrors.name ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {createErrors.name && (
                  <p className="mt-1 text-sm text-red-600">{createErrors.name}</p>
                )}
              </div>
              <div>
                <select
                  value={createType}
                  onChange={(e) => {
                    setCreateType(e.target.value as 'all' | 'image' | 'audio' | 'video');
                    if (createErrors.type) {
                      setCreateErrors((prev) => ({ ...prev, type: '' }));
                    }
                  }}
                  className={`px-3 py-2 border rounded ${
                    createErrors.type ? 'border-red-500' : 'border-gray-300'
                  }`}
                >
                  <option value="all">すべて</option>
                  <option value="image">画像のみ</option>
                  <option value="audio">音楽のみ</option>
                  <option value="video">YouTubeのみ</option>
                </select>
                {createErrors.type && (
                  <p className="mt-1 text-sm text-red-600">{createErrors.type}</p>
                )}
              </div>
              <button
                type="submit"
                disabled={creating}
                className="bg-green-600 text-white px-6 py-2 rounded hover:bg-green-700 disabled:opacity-50"
              >
                {creating ? '作成中...' : '作成'}
              </button>
            </div>
          </form>
        </div>

        {/* タグ一覧 */}
        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <h2 className="text-2xl font-semibold mb-4 text-gray-800">
            タグ一覧
          </h2>

          {loading ? (
            <p className="text-gray-600">読み込み中...</p>
          ) : tagList.length === 0 ? (
            <p className="text-gray-600">タグがありません</p>
          ) : (
            <div className="space-y-4">
              {tagList.map((tag) => (
                <div
                  key={tag.id}
                  className="border border-gray-200 rounded-lg p-4 flex items-center justify-between"
                >
                  {editingId === tag.id ? (
                    <div className="flex-1 flex gap-2">
                      <div className="flex-1">
                        <input
                          type="text"
                          value={editName || ''}
                          onChange={(e) => {
                            setEditName(e.target.value);
                            if (editErrors.name) {
                              setEditErrors((prev) => ({ ...prev, name: '' }));
                            }
                          }}
                          className={`w-full px-3 py-2 border rounded ${
                            editErrors.name ? 'border-red-500' : 'border-gray-300'
                          }`}
                        />
                        {editErrors.name && (
                          <p className="mt-1 text-sm text-red-600">{editErrors.name}</p>
                        )}
                      </div>
                      <div>
                        <select
                          value={editType}
                          onChange={(e) => {
                            setEditType(e.target.value as 'all' | 'image' | 'audio' | 'video');
                            if (editErrors.type) {
                              setEditErrors((prev) => ({ ...prev, type: '' }));
                            }
                          }}
                          className={`px-3 py-2 border rounded ${
                            editErrors.type ? 'border-red-500' : 'border-gray-300'
                          }`}
                        >
                          <option value="all">すべて</option>
                          <option value="image">画像のみ</option>
                          <option value="audio">音楽のみ</option>
                          <option value="video">YouTubeのみ</option>
                        </select>
                        {editErrors.type && (
                          <p className="mt-1 text-sm text-red-600">{editErrors.type}</p>
                        )}
                      </div>
                      <button
                        onClick={() => handleUpdate(tag.id)}
                        disabled={updating}
                        className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
                      >
                        {updating ? '更新中...' : '保存'}
                      </button>
                      <button
                        onClick={handleCancelEdit}
                        className="bg-gray-300 text-gray-700 px-4 py-2 rounded hover:bg-gray-400"
                      >
                        キャンセル
                      </button>
                    </div>
                  ) : (
                    <>
                      <div className="flex-1">
                        <h3 className="font-semibold text-lg text-gray-800">
                          {tag.name}
                        </h3>
                        <p className="text-sm text-gray-500">
                          適用範囲: {tag.type === 'all' ? 'すべて' : tag.type === 'image' ? '画像のみ' : tag.type === 'audio' ? '音楽のみ' : 'YouTubeのみ'} | 作成日: {new Date(tag.created_at).toLocaleString('ja-JP')}
                        </p>
                      </div>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleViewMedia(tag.id)}
                          className="bg-blue-600 text-white px-4 py-2 rounded text-sm hover:bg-blue-700"
                        >
                          関連メディア
                        </button>
                        <button
                          onClick={() => handleStartEdit(tag)}
                          className="bg-yellow-600 text-white px-4 py-2 rounded text-sm hover:bg-yellow-700"
                        >
                          編集
                        </button>
                        <button
                          onClick={() => handleDelete(tag.id)}
                          className="bg-red-600 text-white px-4 py-2 rounded text-sm hover:bg-red-700"
                        >
                          削除
                        </button>
                      </div>
                    </>
                  )}
                </div>
              ))}
            </div>
          )}
        </div>

        {/* 関連メディア表示 */}
        {selectedTagId && (
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-semibold mb-4 text-gray-800">
              関連メディア
            </h2>

            {loadingMedia ? (
              <p className="text-gray-600">読み込み中...</p>
            ) : tagMedia.length === 0 ? (
              <p className="text-gray-600">関連するメディアがありません</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {tagMedia.map((media) => (
                  <div
                    key={media.id}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <h3 className="font-semibold text-lg text-gray-800 mb-2">
                      {media.title}
                    </h3>
                    {media.description && (
                      <p className="text-sm text-gray-600 mb-2">
                        {media.description}
                      </p>
                    )}
                    {media.cloudfront_url && (
                      <img
                        src={media.cloudfront_url}
                        alt={media.title}
                        className="w-full h-48 object-cover rounded mb-2"
                      />
                    )}
                    {media.youtube_url && (() => {
                      const videoId = extractYouTubeId(media.youtube_url);
                      return videoId ? (
                        <div className="mb-2">
                          <div className="relative w-full" style={{ paddingBottom: '56.25%' }}>
                            <iframe
                              className="absolute top-0 left-0 w-full h-full rounded"
                              src={`https://www.youtube.com/embed/${videoId}`}
                              title={media.title}
                              allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                              allowFullScreen
                            />
                          </div>
                          <a
                            href={media.youtube_url}
                            target="_blank"
                            rel="noopener noreferrer"
                            className="text-blue-600 hover:text-blue-800 underline text-xs mt-1 block"
                          >
                            元の動画を開く
                          </a>
                        </div>
                      ) : (
                        <a
                          href={media.youtube_url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 hover:text-blue-800 underline text-sm"
                        >
                          YouTube動画を開く
                        </a>
                      );
                    })()}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
