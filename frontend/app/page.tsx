'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { getMediaList, getTagList, type Media, type Tag } from '@/lib/api';
import Calendar from '@/share/component/Calendar';

export default function Home() {
  const [mediaList, setMediaList] = useState<Media[]>([]);
  const [tagList, setTagList] = useState<Tag[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  
  // カレンダーの年月
  const now = new Date();
  const [currentYear, setCurrentYear] = useState(now.getFullYear());
  const [currentMonth, setCurrentMonth] = useState(now.getMonth() + 1);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);
      const [media, tags] = await Promise.all([
        getMediaList(),
        getTagList(),
      ]);
      setMediaList(media);
      setTagList(tags);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <header className="mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-2">
            Image Server API テスト
          </h1>
          <p className="text-gray-600">
            バックエンドAPIの動作を確認できるフロントエンドです
          </p>
        </header>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-semibold mb-4 text-gray-800">
              メディア管理
            </h2>
            <p className="text-gray-600 mb-4">
              現在のメディア数: {loading ? '読み込み中...' : mediaList.length}
            </p>
            <div className="space-y-2">
              <Link
                href="/media/upload"
                className="inline-block bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition mr-2"
              >
                アップロード
              </Link>
              <Link
                href="/media/list"
                className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
              >
                メディア一覧
              </Link>
            </div>
          </div>

          <div className="bg-white rounded-lg shadow p-6">
            <h2 className="text-2xl font-semibold mb-4 text-gray-800">
              タグ管理
            </h2>
            <p className="text-gray-600 mb-4">
              現在のタグ数: {loading ? '読み込み中...' : tagList.length}
            </p>
            <Link
              href="/tags"
              className="inline-block bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition"
            >
              タグ管理ページへ
            </Link>
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6 mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-semibold text-gray-800">カレンダー</h2>
            <div className="flex items-center gap-4">
              <button
                onClick={() => {
                  if (currentMonth === 1) {
                    setCurrentYear(currentYear - 1);
                    setCurrentMonth(12);
                  } else {
                    setCurrentMonth(currentMonth - 1);
                  }
                }}
                className="px-3 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
              >
                ← 前月
              </button>
              <span className="text-lg font-semibold text-gray-800">
                {currentYear}年{currentMonth}月
              </span>
              <button
                onClick={() => {
                  if (currentMonth === 12) {
                    setCurrentYear(currentYear + 1);
                    setCurrentMonth(1);
                  } else {
                    setCurrentMonth(currentMonth + 1);
                  }
                }}
                className="px-3 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300"
              >
                次月 →
              </button>
            </div>
          </div>
          <Calendar
            year={currentYear}
            month={currentMonth}
            onDateClick={(date) => {
              const year = date.getFullYear();
              const month = String(date.getMonth() + 1).padStart(2, '0');
              const day = String(date.getDate()).padStart(2, '0');
              const dateStr = `${year}-${month}-${day}`;
              window.location.href = `/todos/date?date=${dateStr}`;
            }}
          />
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-4 text-gray-800">TODO管理</h2>
          <div className="space-y-2">
            <Link
              href="/todos/create"
              className="inline-block bg-green-600 text-white px-4 py-2 rounded hover:bg-green-700 transition mr-2"
            >
              TODOを作成
            </Link>
            <Link
              href="/todos/without-due-date"
              className="inline-block bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700 transition"
            >
              期限未設定
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
