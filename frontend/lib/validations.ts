import { z } from 'zod';

// メディアアップロード（ファイル）のバリデーションスキーマ
export const fileUploadSchema = z.object({
  file: z
    .instanceof(File, { message: 'ファイルを選択してください' })
    .refine((file) => file.size > 0, { message: 'ファイルが空です' })
    .refine(
      (file) => {
        const validTypes = [
          'image/jpeg',
          'image/png',
          'image/gif',
          'image/webp',
          'audio/mpeg',
          'audio/mp3',
          'audio/wav',
          'audio/wave',
          'audio/x-wav',
        ];
        return validTypes.some((type) => file.type.startsWith(type.split('/')[0]));
      },
      { message: '画像または音声ファイルを選択してください' }
    )
    .refine((file) => file.size <= 100 * 1024 * 1024, {
      message: 'ファイルサイズは100MB以下にしてください',
    }),
  title: z
    .string()
    .min(1, { message: 'タイトルは必須です' })
    .max(255, { message: 'タイトルは255文字以内で入力してください' }),
  description: z
    .string()
    .max(1000, { message: '説明は1000文字以内で入力してください' })
    .optional()
    .or(z.literal('')),
  tagIds: z.array(z.string().uuid({ message: '無効なタグIDです' })).optional(),
});

// メディアアップロード（YouTube）のバリデーションスキーマ
export const youtubeUploadSchema = z.object({
  youtubeUrl: z
    .string()
    .min(1, { message: 'YouTube URLは必須です' })
    .url({ message: '有効なURLを入力してください' })
    .refine(
      (url) => {
        const youtubePatterns = [
          /^https?:\/\/(www\.)?(youtube\.com|youtu\.be)\/.+/,
          /^https?:\/\/youtube\.com\/watch\?v=[\w-]+/,
          /^https?:\/\/youtu\.be\/[\w-]+/,
          /^https?:\/\/youtube\.com\/embed\/[\w-]+/,
          /^https?:\/\/youtube\.com\/shorts\/[\w-]+/,
        ];
        return youtubePatterns.some((pattern) => pattern.test(url));
      },
      { message: '有効なYouTube URLを入力してください' }
    ),
  title: z
    .string()
    .min(1, { message: 'タイトルは必須です' })
    .max(255, { message: 'タイトルは255文字以内で入力してください' }),
  description: z
    .string()
    .max(1000, { message: '説明は1000文字以内で入力してください' })
    .optional()
    .or(z.literal('')),
  tagIds: z.array(z.string().uuid({ message: '無効なタグIDです' })).optional(),
});

// タグ作成・編集のバリデーションスキーマ
export const tagSchema = z.object({
  name: z
    .string()
    .min(1, { message: 'タグ名は必須です' })
    .max(255, { message: 'タグ名は255文字以内で入力してください' })
    .regex(/^[^\s]+$/, { message: 'タグ名に空白を含めることはできません' }),
  type: z.enum(['all', 'image', 'audio', 'video'], {
    errorMap: () => ({ message: '有効なタグタイプを選択してください' }),
  }),
});

// 型エクスポート
export type FileUploadInput = z.infer<typeof fileUploadSchema>;
export type YouTubeUploadInput = z.infer<typeof youtubeUploadSchema>;
export type TagInput = z.infer<typeof tagSchema>;
