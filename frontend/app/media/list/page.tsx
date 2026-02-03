'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import {
  getMediaListWithPagination,
  getTagList,
  deleteMedia,
  type Media,
  type Tag,
} from '@/lib/api';
import MediaCard from '@/share/component/MediaCard';

export default function MediaListPage() {
  const router = useRouter();
  const [mediaList, setMediaList] = useState<Media[]>([]);
  const [tagList, setTagList] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [offset, setOffset] = useState(0);
  const limit = 20;
  const observerTarget = useRef<HTMLDivElement>(null);

  // フィルターの状態
  const [titleFilter, setTitleFilter] = useState('');
  const [selectedTagIds, setSelectedTagIds] = useState<string[]>([]);

  useEffect(() => {
    loadTags();
    loadInitialData();
  }, []);

  // フィルターが変更されたときにデータを再読み込み
  useEffect(() => {
    loadInitialData();
  }, [titleFilter, selectedTagIds]);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loadingMore && !loading) {
          loadMore();
        }
      },
      { threshold: 0.1 }
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [hasMore, loadingMore, loading, titleFilter, selectedTagIds]);

  const loadTags = async () => {
    try {
      const tags = await getTagList();
      setTagList(tags);
    } catch (err) {
      console.error('Failed to load tags:', err);
    }
  };

  const loadInitialData = async () => {
    try {
      setLoading(true);
      setError(null);
      setOffset(0);
      const response = await getMediaListWithPagination(
        0,
        limit,
        titleFilter || undefined,
        selectedTagIds.length > 0 ? selectedTagIds : undefined
      );
      setMediaList(response.media || []);
      setHasMore(response.has_more ?? false);
      setOffset(response.media?.length || 0);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const loadMore = async () => {
    if (loadingMore || !hasMore) return;

    try {
      setLoadingMore(true);
      const response = await getMediaListWithPagination(
        offset,
        limit,
        titleFilter || undefined,
        selectedTagIds.length > 0 ? selectedTagIds : undefined
      );
      setMediaList((prev) => [...prev, ...(response.media || [])]);
      setHasMore(response.has_more ?? false);
      setOffset((prev) => prev + (response.media?.length || 0));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoadingMore(false);
    }
  };

  const toggleTag = (tagId: string) => {
    setSelectedTagIds((prev) =>
      prev.includes(tagId)
        ? prev.filter((id) => id !== tagId)
        : [...prev, tagId]
    );
  };

  const clearFilters = () => {
    setTitleFilter('');
    setSelectedTagIds([]);
  };

  const handleDelete = async (id: string) => {
    if (!confirm('このメディアを削除しますか？')) {
      return;
    }

    try {
      setError(null);
      setSuccess(null);
      await deleteMedia(id);
      setSuccess('メディアの削除に成功しました');
      // 削除後、リストから該当メディアを削除
      setMediaList((prev) => prev.filter((media) => media.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : '削除に失敗しました');
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
          <div className="flex items-center justify-between mb-4">
            <h1 className="text-4xl font-bold text-gray-900">メディア一覧</h1>
            <Link
              href="/media/upload"
              className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
            >
              メディアをアップロード
            </Link>
          </div>
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

        {/* フィルターセクション */}
        <div className="bg-white rounded-lg shadow p-6 mb-6">
          <h2 className="text-xl font-semibold mb-4 text-gray-800">絞り込み</h2>
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                タイトル検索
              </label>
              <input
                type="text"
                value={titleFilter}
                onChange={(e) => setTitleFilter(e.target.value)}
                placeholder="タイトルで検索..."
                className="w-full px-3 py-2 border border-gray-300 rounded"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                タグで絞り込み
              </label>
              <div className="flex flex-wrap gap-2">
                {tagList.map((tag) => (
                  <button
                    key={tag.id}
                    type="button"
                    onClick={() => toggleTag(tag.id)}
                    className={`px-3 py-1 rounded text-sm ${
                      selectedTagIds.includes(tag.id)
                        ? 'bg-blue-600 text-white'
                        : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                    }`}
                  >
                    {tag.name}
                  </button>
                ))}
              </div>
            </div>
            {(titleFilter || selectedTagIds.length > 0) && (
              <button
                onClick={clearFilters}
                className="text-sm text-blue-600 hover:text-blue-800 underline"
              >
                フィルターをクリア
              </button>
            )}
          </div>
        </div>

        {/* メディア一覧 */}
        <div className="bg-white rounded-lg shadow p-6">
          {loading ? (
            <p className="text-gray-600">読み込み中...</p>
          ) : mediaList.length === 0 ? (
            <p className="text-gray-600">メディアがありません</p>
          ) : (
            <>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {mediaList.map((media) => (
                  <MediaCard
                    key={media.id}
                    media={media}
                    onDelete={handleDelete}
                  />
                ))}
              </div>
              {/* 無限スクロール用の監視要素 */}
              <div ref={observerTarget} className="h-10 flex items-center justify-center">
                {loadingMore && (
                  <p className="text-gray-600">読み込み中...</p>
                )}
                {!hasMore && mediaList.length > 0 && (
                  <p className="text-gray-500 text-sm">すべてのメディアを表示しました</p>
                )}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
