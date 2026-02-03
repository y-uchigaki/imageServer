'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { getTodosByDateRange, type Todo } from '@/lib/api';

interface CalendarProps {
  year: number;
  month: number;
  onDateClick?: (date: Date) => void;
}

export default function Calendar({ year, month, onDateClick }: CalendarProps) {
  const router = useRouter();
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loading, setLoading] = useState(true);

  // 月の最初の日と最後の日を計算
  const firstDay = new Date(year, month - 1, 1);
  const lastDay = new Date(year, month, 0);
  const startOfMonth = new Date(year, month - 1, 1);
  const endOfMonth = new Date(year, month, 0, 23, 59, 59);

  // 日付をYYYY-MM-DD形式の文字列に変換（ローカルタイムゾーン）
  const formatDateString = (date: Date): string => {
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    return `${year}-${month}-${day}`;
  };

  useEffect(() => {
    loadTodos();
  }, [year, month]);

  const loadTodos = async () => {
    try {
      setLoading(true);
      const startDate = formatDateString(startOfMonth);
      const endDate = formatDateString(endOfMonth);
      const todosData = await getTodosByDateRange(startDate, endDate);
      setTodos(todosData);
    } catch (err) {
      console.error('Failed to load todos:', err);
    } finally {
      setLoading(false);
    }
  };

  // カレンダーの日付配列を生成
  const getCalendarDays = () => {
    const days: (Date | null)[] = [];
    const startDate = new Date(firstDay);
    startDate.setDate(startDate.getDate() - startDate.getDay()); // 週の最初の日（日曜日）に合わせる

    for (let i = 0; i < 42; i++) {
      const date = new Date(startDate);
      date.setDate(startDate.getDate() + i);
      days.push(date);
    }

    return days;
  };

  // 特定の日付のTODOを取得
  const getTodosForDate = (date: Date): Todo[] => {
    const dateStr = formatDateString(date);
    return todos.filter((todo) => {
      // 期間指定のTODO
      if (todo.start_date && todo.end_date) {
        const startStr = todo.start_date.split('T')[0];
        const endStr = todo.end_date.split('T')[0];
        return dateStr >= startStr && dateStr <= endStr;
      }
      // 期限日のみのTODO
      if (todo.due_date) {
        const dueStr = todo.due_date.split('T')[0];
        return dateStr === dueStr;
      }
      return false;
    });
  };

  // 日付が現在の月に属するか
  const isCurrentMonth = (date: Date | null): boolean => {
    if (!date) return false;
    return date.getMonth() === month - 1 && date.getFullYear() === year;
  };

  // 日付が今日か
  const isToday = (date: Date | null): boolean => {
    if (!date) return false;
    const today = new Date();
    return (
      date.getDate() === today.getDate() &&
      date.getMonth() === today.getMonth() &&
      date.getFullYear() === today.getFullYear()
    );
  };

  const handleDateClick = (date: Date) => {
    const dateStr = formatDateString(date);
    if (onDateClick) {
      onDateClick(date);
    } else {
      router.push(`/todos/date?date=${dateStr}`);
    }
  };

  const days = getCalendarDays();
  const weekDays = ['日', '月', '火', '水', '木', '金', '土'];

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="grid grid-cols-7 gap-1">
        {/* 曜日ヘッダー */}
        {weekDays.map((day, index) => (
          <div
            key={index}
            className="text-center font-semibold text-gray-700 py-2"
          >
            {day}
          </div>
        ))}

        {/* カレンダーの日付 */}
        {days.map((date, index) => {
          if (!date) return <div key={index} />;

          const dateTodos = getTodosForDate(date);
          const isCurrentMonthDay = isCurrentMonth(date);
          const isTodayDay = isToday(date);

          return (
            <div
              key={index}
              className={`
                min-h-[80px] border border-gray-200 p-1 cursor-pointer
                ${isCurrentMonthDay ? 'bg-white' : 'bg-gray-50'}
                ${isTodayDay ? 'ring-2 ring-blue-500' : ''}
                hover:bg-blue-50 transition
              `}
              onClick={() => handleDateClick(date)}
            >
              <div
                className={`
                  text-sm font-medium mb-1
                  ${isCurrentMonthDay ? 'text-gray-900' : 'text-gray-400'}
                  ${isTodayDay ? 'text-blue-600 font-bold' : ''}
                `}
              >
                {date.getDate()}
              </div>

              {/* TODO表示 */}
              <div className="space-y-0.5">
                {dateTodos.map((todo) => {
                  const isPeriod = todo.start_date && todo.end_date;
                  const isDueDate = todo.due_date && !todo.start_date;

                  if (isPeriod) {
                    // 期間指定のTODO - 矢印で表示
                    const dateStr = formatDateString(date);
                    const startStr = todo.start_date!.split('T')[0];
                    const endStr = todo.end_date!.split('T')[0];
                    const isStart = dateStr === startStr;
                    const isEnd = dateStr === endStr;
                    const isMiddle = dateStr > startStr && dateStr < endStr;

                    return (
                      <div
                        key={todo.id}
                        className={`
                          text-xs px-1 py-0.5 rounded truncate
                          ${todo.completed ? 'bg-gray-300 text-gray-600' : 'bg-blue-200 text-blue-800'}
                        `}
                        title={todo.title}
                      >
                        {isStart && '→ '}
                        {isMiddle && '→ '}
                        {isEnd && '← '}
                        {todo.title}
                      </div>
                    );
                  } else if (isDueDate) {
                    // 期限日のみのTODO
                    return (
                      <div
                        key={todo.id}
                        className={`
                          text-xs px-1 py-0.5 rounded truncate
                          ${todo.completed ? 'bg-gray-300 text-gray-600' : 'bg-green-200 text-green-800'}
                        `}
                        title={todo.title}
                      >
                        {todo.title}
                      </div>
                    );
                  }
                  return null;
                })}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
