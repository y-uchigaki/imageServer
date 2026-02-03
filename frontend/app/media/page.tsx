'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function MediaPage() {
  const router = useRouter();
  
  useEffect(() => {
    router.replace('/media/list');
  }, [router]);

  return null;
}
