'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useState, useEffect } from 'react';
import { getTodosByDate, type Todo } from '@/lib/api';
import Link from 'next/link';

export default function TodosByDatePage() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const dateStr = searchParams.get('date') || new Date().toISOString().split('T')[0];

  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadTodos();
  }, [dateStr]);

  const loadTodos = async () => {
    try {
      setLoading(true);
      setError(null);
      const todosData = await getTodosByDate(dateStr);
      setTodos(todosData);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const date = new Date(dateStr);
  const formattedDate = `${date.getFullYear()}年${date.getMonth() + 1}月${date.getDate()}日`;

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-6">
          <button
            onClick={() => router.back()}
            className="text-blue-600 hover:text-blue-800 underline mb-4"
          >
            ← 戻る
          </button>
          <h1 className="text-4xl font-bold text-gray-900">
            {formattedDate}のTODO
          </h1>
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
            <p className="text-gray-600">この日のTODOはありません</p>
          </div>
        ) : (
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
                      {todo.start_date && todo.end_date && (
                        <span>
                          期間: {new Date(todo.start_date).toLocaleDateString('ja-JP')} ～{' '}
                          {new Date(todo.end_date).toLocaleDateString('ja-JP')}
                        </span>
                      )}
                      {todo.due_date && !todo.start_date && (
                        <span>
                          期限: {new Date(todo.due_date).toLocaleDateString('ja-JP')}
                        </span>
                      )}
                      <span className={todo.completed ? 'text-green-600' : 'text-orange-600'}>
                        {todo.completed ? '完了' : '未完了'}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
