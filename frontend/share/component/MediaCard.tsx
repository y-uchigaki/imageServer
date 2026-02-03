'use client';

import { useRouter } from 'next/navigation';
import { type Media } from '@/lib/api';

// YouTube URLからIDを抽出する関数
function extractYouTubeId(url: string): string | null {
  const patterns = [
    /(?:youtube\.com\/watch\?v=|youtu\.be\/|youtube\.com\/embed\/|youtube\.com\/shorts\/)([^&\n?#\/]+)/,
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

// YouTube Shortsかどうかを判定する関数
function isYouTubeShorts(url: string): boolean {
  return /youtube\.com\/shorts\//.test(url);
}

interface MediaCardProps {
  media: Media;
  onDelete: (id: string) => void;
}

export default function MediaCard({ media, onDelete }: MediaCardProps) {
  const router = useRouter();

  const handleCardClick = (e: React.MouseEvent) => {
    // 削除ボタンがクリックされた場合は遷移しない
    if ((e.target as HTMLElement).closest('button')) {
      return;
    }
    router.push(`/media/${media.id}/edit`);
  };

  return (
    <div
      className="border border-gray-200 rounded-lg p-4 hover:shadow-lg transition cursor-pointer"
      onClick={handleCardClick}
    >
      <div className="mb-2">
        <h3 className="font-semibold text-lg text-gray-800">
          {media.title}
        </h3>
      </div>

      {media.cloudfront_url && media.type === 'image' && (
        <div className="mb-2">
          <img
            src={media.cloudfront_url}
            alt={media.title}
            className="w-full h-48 object-cover rounded"
          />
        </div>
      )}

      {media.cloudfront_url && media.type === 'audio' && (
        <div className="mb-2">
          <audio
            controls
            className="w-full"
            src={media.cloudfront_url || ''}
          >
            お使いのブラウザは音声再生に対応していません。
          </audio>
        </div>
      )}

      {media.youtube_url && (() => {
        const videoId = extractYouTubeId(media.youtube_url);
        const isShorts = isYouTubeShorts(media.youtube_url);
        // Shortsの場合は縦長（9:16 = 177.78%）、通常の場合は横長（16:9 = 56.25%）
        const aspectRatio = isShorts ? '90%' : '56.25%';
        const maxWidth = isShorts ? 'max-w-sm mx-auto' : '';
        
        return videoId ? (
          <div className="mb-2">
            <div className={`relative w-full ${maxWidth}`} style={{ paddingBottom: aspectRatio }}>
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
          <div className="mb-2">
            <a
              href={media.youtube_url}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:text-blue-800 underline text-sm"
            >
              YouTube動画を開く
            </a>
          </div>
        );
      })()}

      {media.tags.length > 0 && (
        <div className="mb-2">
          <div className="flex flex-wrap gap-1">
            {media.tags.map((tag) => (
              <span
                key={tag.id}
                className="px-2 py-1 bg-gray-100 text-gray-700 text-xs rounded"
              >
                {tag.name}
              </span>
            ))}
          </div>
        </div>
      )}
      
      {media.description && !media.tags.find((tag) => tag.name === '動画') && (
        <p className="text-sm text-gray-600 mb-2"> 
          {media.description}
        </p>
      )}

      <div className="flex items-center justify-between">
        <div className="text-xs text-gray-500">
          <p>作成日: {new Date(media.created_at).toLocaleString('ja-JP')}</p>
        </div>
        <button
          onClick={() => onDelete(media.id)}
          className="bg-red-600 text-white px-3 py-1 rounded text-sm hover:bg-red-700"
        >
          削除
        </button>
      </div>
    </div>
  );
}
