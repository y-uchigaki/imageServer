'use client';

import { useRouter, useParams } from 'next/navigation';
import { useState, useEffect } from 'react';
import { getTodo, updateTodo, deleteTodo, type Todo } from '@/lib/api';
import { z } from 'zod';

const todoSchema = z.object({
  title: z.string().min(1, { message: 'タイトルは必須です' }).max(255, { message: 'タイトルは255文字以内で入力してください' }),
  description: z.string().max(1000, { message: '説明は1000文字以内で入力してください' }).optional().or(z.literal('')),
  startDate: z.string().optional().or(z.literal('')),
  endDate: z.string().optional().or(z.literal('')),
  dueDate: z.string().optional().or(z.literal('')),
  completed: z.boolean(),
});

export default function TodoEditPage() {
  const router = useRouter();
  const params = useParams();
  const todoId = params.id as string;

  const [todo, setTodo] = useState<Todo | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [validationErrors, setValidationErrors] = useState<Record<string, string>>({});

  // フォームの状態
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [dateType, setDateType] = useState<'none' | 'period' | 'due'>('none');
  const [startDate, setStartDate] = useState('');
  const [endDate, setEndDate] = useState('');
  const [dueDate, setDueDate] = useState('');
  const [completed, setCompleted] = useState(false);

  useEffect(() => {
    loadTodo();
  }, [todoId]);

  const loadTodo = async () => {
    try {
      setLoading(true);
      setError(null);
      const todoData = await getTodo(todoId);
      setTodo(todoData);
      setTitle(todoData.title);
      setDescription(todoData.description || '');
      setCompleted(todoData.completed || false);

      if (todoData.start_date && todoData.end_date) {
        setDateType('period');
        setStartDate(todoData.start_date.split('T')[0]);
        setEndDate(todoData.end_date.split('T')[0]);
        setDueDate('');
      } else if (todoData.due_date) {
        setDateType('due');
        setDueDate(todoData.due_date.split('T')[0]);
        setStartDate('');
        setEndDate('');
      } else {
        setDateType('none');
        setStartDate('');
        setEndDate('');
        setDueDate('');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'データの取得に失敗しました');
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setValidationErrors({});
    setError(null);
    setSuccess(null);

    // バリデーション
    const result = todoSchema.safeParse({
      title,
      description: description || undefined,
      startDate: dateType === 'period' ? startDate : undefined,
      endDate: dateType === 'period' ? endDate : undefined,
      dueDate: dateType === 'due' ? dueDate : undefined,
      completed,
    });

    if (!result.success) {
      const errors: Record<string, string> = {};
      result.error.issues.forEach((issue) => {
        const path = issue.path.join('.');
        errors[path] = issue.message;
      });
      setValidationErrors(errors);
      return;
    }

    try {
      setSaving(true);
      await updateTodo(
        todoId,
        title,
        description || undefined,
        dateType === 'period' ? startDate : undefined,
        dateType === 'period' ? endDate : undefined,
        dateType === 'due' ? dueDate : undefined,
        completed
      );
      setSuccess('TODOの更新に成功しました');
      setTimeout(() => {
        router.push('/');
      }, 1500);
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新に失敗しました');
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm('本当にこのTODOを削除しますか？')) {
      return;
    }
    try {
      setError(null);
      await deleteTodo(todoId);
      router.push('/');
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

  if (!todo) {
    return (
      <div className="container mx-auto px-4 py-8">
        <p className="text-red-600">TODOが見つかりません</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
        <div className="mb-4">
          <button
            onClick={() => router.back()}
            className="text-blue-600 hover:text-blue-800 underline"
          >
            ← 戻る
          </button>
        </div>

        <h1 className="text-3xl font-bold mb-6 text-gray-800">TODO編集</h1>

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

        <div className="bg-white rounded-lg shadow p-6 space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              タイトル *
            </label>
            <input
              type="text"
              value={title}
              onChange={(e) => {
                setTitle(e.target.value);
                if (validationErrors.title) {
                  setValidationErrors((prev) => ({ ...prev, title: '' }));
                }
              }}
              className={`w-full px-3 py-2 border rounded ${
                validationErrors.title ? 'border-red-500' : 'border-gray-300'
              }`}
            />
            {validationErrors.title && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.title}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              説明
            </label>
            <textarea
              value={description}
              onChange={(e) => {
                setDescription(e.target.value);
                if (validationErrors.description) {
                  setValidationErrors((prev) => ({ ...prev, description: '' }));
                }
              }}
              className={`w-full px-3 py-2 border rounded ${
                validationErrors.description ? 'border-red-500' : 'border-gray-300'
              }`}
              rows={3}
            />
            {validationErrors.description && (
              <p className="mt-1 text-sm text-red-600">{validationErrors.description}</p>
            )}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              日付設定
            </label>
            <div className="space-y-2">
              <label className="flex items-center">
                <input
                  type="radio"
                  name="dateType"
                  value="none"
                  checked={dateType === 'none'}
                  onChange={() => {
                    setDateType('none');
                    setStartDate('');
                    setEndDate('');
                    setDueDate('');
                  }}
                  className="mr-2"
                />
                日付なし
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  name="dateType"
                  value="period"
                  checked={dateType === 'period'}
                  onChange={() => {
                    setDateType('period');
                    setDueDate('');
                  }}
                  className="mr-2"
                />
                期間（開始日～終了日）
              </label>
              <label className="flex items-center">
                <input
                  type="radio"
                  name="dateType"
                  value="due"
                  checked={dateType === 'due'}
                  onChange={() => {
                    setDateType('due');
                    setStartDate('');
                    setEndDate('');
                  }}
                  className="mr-2"
                />
                期限日のみ
              </label>
            </div>

            {dateType === 'period' && (
              <div className="mt-4 grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    開始日
                  </label>
                  <input
                    type="date"
                    value={startDate}
                    onChange={(e) => setStartDate(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    終了日
                  </label>
                  <input
                    type="date"
                    value={endDate}
                    onChange={(e) => setEndDate(e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded"
                  />
                </div>
              </div>
            )}

            {dateType === 'due' && (
              <div className="mt-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  期限日
                </label>
                <input
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded"
                />
              </div>
            )}
          </div>

          <div>
            <label className="flex items-center">
              <input
                type="checkbox"
                checked={completed}
                onChange={(e) => setCompleted(e.target.checked)}
                className="mr-2"
              />
              完了済み
            </label>
          </div>

          <div className="flex gap-4 pt-4">
            <button
              onClick={handleSave}
              disabled={saving}
              className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:opacity-50"
            >
              {saving ? '保存中...' : '保存'}
            </button>
            <button
              onClick={handleDelete}
              className="bg-red-600 text-white px-6 py-2 rounded hover:bg-red-700"
            >
              削除
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
