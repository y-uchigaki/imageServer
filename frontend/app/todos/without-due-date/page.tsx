'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { getTodosWithoutDueDate, type Todo } from '@/lib/api';

export default function TodosWithoutDueDatePage() {
  const router = useRouter();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [hasMore, setHasMore] = useState(true);
  const [offset, setOffset] = useState(0);
  const limit = 20;
  const observerTarget = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadInitialData();
  }, []);

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
  }, [hasMore, loadingMore, loading]);

  const loadInitialData = async () => {
    try {
      setLoading(true);
      setError(null);
      const response = await getTodosWithoutDueDate(0, limit);
      setTodos(response.todos || []);
      setHasMore(response.has_more ?? false);
      setOffset(response.todos?.length || 0);
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
      const response = await getTodosWithoutDueDate(offset, limit);
      setTodos((prev) => [...prev, ...(response.todos || [])]);
      setHasMore(response.has_more ?? false);
      setOffset((prev) => prev + (response.todos?.length || 0));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoadingMore(false);
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
          <h1 className="text-4xl font-bold text-gray-900">期限未設定のTODO</h1>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
            {error}
          </div>
        )}

        {loading ? (
          <p className="text-gray-600">読み込み中...</p>
        ) : todos.length === 0 ? (
          <div className="bg-white rounded-lg shadow p-6">
            <p className="text-gray-600">期限未設定のTODOはありません</p>
          </div>
        ) : (
          <>
            <div className="space-y-4">
              {todos.map((todo) => (
                <div
                  key={todo.id}
                  className="bg-white rounded-lg shadow p-6 hover:shadow-lg transition cursor-pointer"
                  onClick={() => router.push(`/todos/${todo.id}/edit`)}
                >
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="text-xl font-semibold text-gray-800 mb-2">
                        {todo.title}
                      </h3>
                      {todo.description && (
                        <p className="text-gray-600 mb-2">{todo.description}</p>
                      )}
                      <div className="flex gap-4 text-sm text-gray-500">
                        <span className={todo.completed ? 'text-green-600' : 'text-orange-600'}>
                          {todo.completed ? '完了' : '未完了'}
                        </span>
                        <span>
                          作成日: {new Date(todo.created_at).toLocaleDateString('ja-JP')}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            {/* 無限スクロール用の監視要素 */}
            <div ref={observerTarget} className="h-10 flex items-center justify-center mt-4">
              {loadingMore && (
                <p className="text-gray-600">読み込み中...</p>
              )}
              {!hasMore && todos.length > 0 && (
                <p className="text-gray-500 text-sm">すべてのTODOを表示しました</p>
              )}
            </div>
          </>
        )}
      </div>
    </div>
  );
}
