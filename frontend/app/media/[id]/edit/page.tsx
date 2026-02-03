'use client';

import { useRouter, useParams } from 'next/navigation';
import { useState, useEffect } from 'react';
import { getMedia, getTagList, associateTag, removeTag, deleteMedia, type Media, type Tag } from '@/lib/api';

export default function MediaEditPage() {
  const router = useRouter();
  const params = useParams();
  const mediaId = params.id as string;

  const [media, setMedia] = useState<Media | null>(null);
  const [tagList, setTagList] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  // 編集フォームの状態
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [selectedTagIds, setSelectedTagIds] = useState<string[]>([]);

  useEffect(() => {
    loadData();
  }, [mediaId]);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const [mediaData, tagsData] = await Promise.all([
        getMedia(mediaId),
        getTagList(),
      ]);
      setMedia(mediaData);
      setTitle(mediaData.title);
      setDescription(mediaData.description || '');
      setSelectedTagIds(mediaData.tags.map((tag) => tag.id));
      
      // メディアタイプに応じてタグをフィルタリング
      const filteredTags = tagsData.filter((tag) => {
        if (tag.type === 'all') return true;
        if (mediaData.type === 'image' && tag.type === 'image') return true;
        if (mediaData.type === 'audio' && tag.type === 'audio') return true;
        if (mediaData.type === 'video' && tag.type === 'video') return true;
        return false;
      });
      setTagList(filteredTags);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの読み込みに失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleAddTag = async (tagId: string) => {
    if (selectedTagIds.includes(tagId)) {
      return;
    }
    try {
      setError(null);
      setSuccess(null);
      await associateTag(mediaId, tagId);
      setSelectedTagIds([...selectedTagIds, tagId]);
      setSuccess('タグを追加しました');
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'タグの追加に失敗しました');
    }
  };

  const handleRemoveTag = async (tagId: string) => {
    try {
      setError(null);
      setSuccess(null);
      await removeTag(mediaId, tagId);
      setSelectedTagIds(selectedTagIds.filter((id) => id !== tagId));
      setSuccess('タグを削除しました');
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'タグの削除に失敗しました');
    }
  };

  const handleDelete = async () => {
    if (!confirm('本当にこのメディアを削除しますか？')) {
      return;
    }
    try {
      setError(null);
      await deleteMedia(mediaId);
      router.push('/media/list');
    } catch (err) {
      setError(err instanceof Error ? err.message : '削除に失敗しました');
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p className="text-gray-600">読み込み中...</p>
      </div>
    );
  }

  if (!media) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p className="text-red-600">メディアが見つかりません</p>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="mb-4">
        <button
          onClick={() => router.back()}
          className="text-blue-600 hover:text-blue-800 underline"
        >
          ← 戻る
        </button>
      </div>

      <h1 className="text-3xl font-bold mb-6 text-gray-800">メディア編集</h1>

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

      <div className="bg-white rounded-lg shadow p-6 mb-6">
        <h2 className="text-2xl font-semibold mb-4 text-gray-800">
          {media.title}
        </h2>

        {/* メディア表示 */}
        {media.cloudfront_url && media.type === 'image' && (
          <div className="mb-4">
            <img
              src={media.cloudfront_url}
              alt={media.title}
              className="w-full max-w-md h-auto rounded"
            />
          </div>
        )}

        {media.cloudfront_url && media.type === 'audio' && (
          <div className="mb-4">
            <audio controls src={media.cloudfront_url} className="w-full max-w-md">
              お使いのブラウザは音声再生に対応していません。
            </audio>
          </div>
        )}

        {media.youtube_url && (
          <div className="mb-4">
            <div className="relative w-full max-w-2xl" style={{ paddingBottom: '56.25%' }}>
              <iframe
                className="absolute top-0 left-0 w-full h-full rounded"
                src={`https://www.youtube.com/embed/${media.youtube_url.match(/(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/|youtube\.com\/shorts\/)([^&\n?#\/]+)/)?.[1] || ''}`}
                title={media.title}
                allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
                allowFullScreen
              />
            </div>
          </div>
        )}

        {/* 現在のタグ */}
        {media.tags.length > 0 && (
          <div className="mb-4">
            <h3 className="font-semibold mb-2">現在のタグ</h3>
            <div className="flex flex-wrap gap-2">
              {media.tags.map((tag) => (
                <span
                  key={tag.id}
                  className="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm flex items-center gap-2"
                >
                  {tag.name}
                  <button
                    onClick={() => handleRemoveTag(tag.id)}
                    className="text-red-600 hover:text-red-800 font-bold"
                  >
                    ×
                  </button>
                </span>
              ))}
            </div>
          </div>
        )}

        {/* タグ追加 */}
        <div className="mb-4">
          <h3 className="font-semibold mb-2">タグを追加</h3>
          <div className="flex flex-wrap gap-2">
            {tagList
              .filter((tag) => !selectedTagIds.includes(tag.id))
              .map((tag) => (
                <button
                  key={tag.id}
                  onClick={() => handleAddTag(tag.id)}
                  className="px-3 py-1 bg-gray-100 text-gray-700 rounded-full text-sm hover:bg-gray-200"
                >
                  + {tag.name}
                </button>
              ))}
          </div>
        </div>

        {/* 説明 */}
        {media.description && (
          <div className="mb-4">
            <p className="text-gray-600">{media.description}</p>
          </div>
        )}

        {/* 削除ボタン */}
        <div className="mt-6">
          <button
            onClick={handleDelete}
            className="bg-red-600 text-white px-4 py-2 rounded hover:bg-red-700"
          >
            メディアを削除
          </button>
        </div>
      </div>
    </div>
  );
}
