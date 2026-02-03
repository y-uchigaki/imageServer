'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import {
  uploadMedia,
  createMediaWithYouTube,
  getTagList,
  type Tag,
} from '@/lib/api';
import { useTags } from '@/share/globalUseState/useTags';
import { useLode } from '@/share/globalUseState/useLode';
import { fileUploadSchema, youtubeUploadSchema } from '@/lib/validations';

export default function MediaUploadPage() {
  const router = useRouter();
  const { allTags, tagList, setAllTags, setTagList } = useTags();
  const {  error, success, setLoading, setError, setSuccess } = useLode();

  // アップロードフォームの状態
  const [uploadMode, setUploadMode] = useState<'file' | 'youtube'>('file');
  const [uploadTitle, setUploadTitle] = useState('');
  const [uploadDescription, setUploadDescription] = useState('');
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploadYouTubeUrl, setUploadYouTubeUrl] = useState('');
  const [selectedTagIds, setSelectedTagIds] = useState<string[]>([]);
  const [uploading, setUploading] = useState(false);
  
  // バリデーションエラー
  const [fileErrors, setFileErrors] = useState<Record<string, string>>({});
  const [youtubeErrors, setYoutubeErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    loadTags();
  }, []);

  const loadTags = async () => {
    try {
      setLoading(true);
      setError(null);
      const tags = await getTagList();
      setAllTags(tags);
      filterTagsByMediaType(tags, null);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'タグの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const filterTagsByMediaType = (tags: Tag[], mediaType: 'image' | 'audio' | 'video' | null) => {
    if (mediaType === null) {
      // ファイルが選択されていない場合はすべてのタグを表示
      setTagList(tags.filter((tag) => tag.type === 'all'));
      return;
    }
    const filtered = tags.filter((tag) => {
      if (tag.type === 'all') return true;
      return tag.type === mediaType;
    });
    setTagList(filtered);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setUploadFile(file);
    setSelectedTagIds([]); // ファイル変更時に選択タグをリセット
    
    if (file) {
      // ファイルタイプを判定
      const isImage = file.type.startsWith('image/');
      const isAudio = file.type.startsWith('audio/');
      
      if (isImage) {
        filterTagsByMediaType(allTags, 'image');
      } else if (isAudio) {
        filterTagsByMediaType(allTags, 'audio');
      } else {
        filterTagsByMediaType(allTags, null);
      }
    } else {
      filterTagsByMediaType(allTags, null);
    }
  };

  const handleFileUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    setFileErrors({});
    setError(null);

    // バリデーション
    const result = fileUploadSchema.safeParse({
      file: uploadFile,
      title: uploadTitle,
      description: uploadDescription || undefined,
      tagIds: selectedTagIds.length > 0 ? selectedTagIds : undefined,
    });

    if (!result.success) {
      const errors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        const path = issue.path.join('.');
        errors[path] = issue.message;
      });
      setFileErrors(errors);
      return;
    }

    try {
      setUploading(true);
      setError(null);
      setSuccess(null);
      await uploadMedia(
        uploadFile!,
        uploadTitle,
        uploadDescription || undefined,
        selectedTagIds.length > 0 ? selectedTagIds : undefined
      );
      setSuccess('メディアのアップロードに成功しました');
      setUploadTitle('');
      setUploadDescription('');
      setUploadFile(null);
      setSelectedTagIds([]);
      setFileErrors({});
      // アップロード成功後、メディア一覧ページに遷移
      setTimeout(() => {
        router.push('/media/list');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'アップロードに失敗しました');
    } finally {
      setUploading(false);
    }
  };

  const handleYouTubeUpload = async (e: React.FormEvent) => {
    e.preventDefault();
    setYoutubeErrors({});
    setError(null);

    // バリデーション
    const result = youtubeUploadSchema.safeParse({
      youtubeUrl: uploadYouTubeUrl,
      title: uploadTitle,
      description: uploadDescription || undefined,
      tagIds: selectedTagIds.length > 0 ? selectedTagIds : undefined,
    });

    if (!result.success) {
      const errors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        const path = issue.path.join('.');
        errors[path] = issue.message;
      });
      setYoutubeErrors(errors);
      return;
    }

    try {
      setUploading(true);
      setError(null);
      setSuccess(null);
      await createMediaWithYouTube(
        uploadYouTubeUrl,
        uploadTitle,
        uploadDescription || undefined,
        selectedTagIds.length > 0 ? selectedTagIds : undefined
      );
      setSuccess('YouTubeメディアの作成に成功しました');
      setUploadTitle('');
      setUploadDescription('');
      setUploadYouTubeUrl('');
      setSelectedTagIds([]);
      setYoutubeErrors({});
      // アップロード成功後、メディア一覧ページに遷移
      setTimeout(() => {
        router.push('/media/list');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : '作成に失敗しました');
    } finally {
      setUploading(false);
    }
  };

  const toggleTag = (tagId: string) => {
    setSelectedTagIds((prev) =>
      prev.includes(tagId)
        ? prev.filter((id) => id !== tagId)
        : [...prev, tagId]
    );
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
            <h1 className="text-4xl font-bold text-gray-900">メディアをアップロード</h1>
            <Link
              href="/media/list"
              className="text-blue-600 hover:text-blue-800 underline"
            >
              メディア一覧を見る
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

        {/* アップロードフォーム */}
        <div className="bg-white rounded-lg shadow p-6">
          <div className="mb-4">
            <button
              onClick={() => {
                setUploadMode('file');
                setUploadTitle('');
                setUploadDescription('');
                setUploadFile(null);
                setUploadYouTubeUrl('');
                setSelectedTagIds([]);
                setFileErrors({});
                setYoutubeErrors({});
                filterTagsByMediaType(allTags, null);
              }}
              className={`px-4 py-2 rounded mr-2 ${
                uploadMode === 'file'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700'
              }`}
            >
              ファイルアップロード
            </button>
            <button
              onClick={() => {
                setUploadMode('youtube');
                setUploadTitle('');
                setUploadDescription('');
                setUploadFile(null);
                setUploadYouTubeUrl('');
                setSelectedTagIds([]);
                setFileErrors({});
                setYoutubeErrors({});
                filterTagsByMediaType(allTags, 'video');
              }}
              className={`px-4 py-2 rounded ${
                uploadMode === 'youtube'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700'
              }`}
            >
              YouTube URL
            </button>
          </div>

          {uploadMode === 'file' ? (
            <form key="file-upload" onSubmit={handleFileUpload} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  ファイル *
                </label>
                <input
                  type="file"
                  accept="image/*,audio/mpeg,audio/mp3,audio/wav,audio/wave,audio/*,.mp3,.wav,.wave"
                  onChange={handleFileChange}
                  className={`w-full px-3 py-2 border rounded ${
                    fileErrors.file ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {fileErrors.file && (
                  <p className="mt-1 text-sm text-red-600">{fileErrors.file}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  タイトル *
                </label>
                <input
                  type="text"
                  value={uploadTitle ?? ''}
                  onChange={(e) => {
                    setUploadTitle(e.target.value);
                    if (fileErrors.title) {
                      setFileErrors((prev) => ({ ...prev, title: '' }));
                    }
                  }}
                  className={`w-full px-3 py-2 border rounded ${
                    fileErrors.title ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {fileErrors.title && (
                  <p className="mt-1 text-sm text-red-600">{fileErrors.title}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  説明
                </label>
                <textarea
                  value={uploadDescription ?? ''}
                  onChange={(e) => {
                    setUploadDescription(e.target.value);
                    if (fileErrors.description) {
                      setFileErrors((prev) => ({ ...prev, description: '' }));
                    }
                  }}
                  className={`w-full px-3 py-2 border rounded ${
                    fileErrors.description ? 'border-red-500' : 'border-gray-300'
                  }`}
                  rows={3}
                />
                {fileErrors.description && (
                  <p className="mt-1 text-sm text-red-600">{fileErrors.description}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  タグ
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
                          : 'bg-gray-200 text-gray-700'
                      }`}
                    >
                      {tag.name}
                    </button>
                  ))}
                </div>
              </div>
              <button
                type="submit"
                disabled={uploading}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {uploading ? 'アップロード中...' : 'アップロード'}
              </button>
            </form>
          ) : (
            <form key="youtube-upload" onSubmit={handleYouTubeUpload} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  YouTube URL *
                </label>
                <input
                  type="url"
                  value={uploadYouTubeUrl ?? ''}
                  onChange={(e) => {
                    setUploadYouTubeUrl(e.target.value);
                    if (youtubeErrors.youtubeUrl) {
                      setYoutubeErrors((prev) => ({ ...prev, youtubeUrl: '' }));
                    }
                  }}
                  placeholder="https://www.youtube.com/watch?v=..."
                  className={`w-full px-3 py-2 border rounded ${
                    youtubeErrors.youtubeUrl ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {youtubeErrors.youtubeUrl && (
                  <p className="mt-1 text-sm text-red-600">{youtubeErrors.youtubeUrl}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  タイトル *
                </label>
                <input
                  type="text"
                  value={uploadTitle ?? ''}
                  onChange={(e) => {
                    setUploadTitle(e.target.value);
                    if (youtubeErrors.title) {
                      setYoutubeErrors((prev) => ({ ...prev, title: '' }));
                    }
                  }}
                  className={`w-full px-3 py-2 border rounded ${
                    youtubeErrors.title ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {youtubeErrors.title && (
                  <p className="mt-1 text-sm text-red-600">{youtubeErrors.title}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  説明
                </label>
                <textarea
                  value={uploadDescription ?? ''}
                  onChange={(e) => {
                    setUploadDescription(e.target.value);
                    if (youtubeErrors.description) {
                      setYoutubeErrors((prev) => ({ ...prev, description: '' }));
                    }
                  }}
                  className={`w-full px-3 py-2 border rounded ${
                    youtubeErrors.description ? 'border-red-500' : 'border-gray-300'
                  }`}
                  rows={3}
                />
                {youtubeErrors.description && (
                  <p className="mt-1 text-sm text-red-600">{youtubeErrors.description}</p>
                )}
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  タグ
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
                          : 'bg-gray-200 text-gray-700'
                      }`}
                    >
                      {tag.name}
                    </button>
                  ))}
                </div>
              </div>
              <button
                type="submit"
                disabled={uploading}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
              >
                {uploading ? '作成中...' : '作成'}
              </button>
            </form>
          )}
        </div>
      </div>
    </div>
  );
}
