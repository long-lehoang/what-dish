'use client';

import { useCallback } from 'react';
import { cn } from '@shared/lib/utils';

interface SearchBarProps {
  value: string;
  onChange: (value: string) => void;
  className?: string;
}

export function SearchBar({ value, onChange, className }: SearchBarProps) {
  const handleClear = useCallback(() => {
    onChange('');
  }, [onChange]);

  return (
    <div className={cn('relative', className)}>
      {/* Magnifying glass icon */}
      <svg
        className="absolute left-3.5 top-1/2 h-5 w-5 -translate-y-1/2 text-gray-400"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        strokeWidth="2"
        strokeLinecap="round"
        strokeLinejoin="round"
        aria-hidden="true"
      >
        <circle cx="11" cy="11" r="8" />
        <path d="m21 21-4.3-4.3" />
      </svg>

      <input
        type="text"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Tìm món ăn..."
        inputMode="search"
        className="w-full rounded-xl border border-gray-200 bg-white py-3 pl-11 pr-10 text-sm text-gray-900 placeholder:text-gray-400 focus:border-primary focus:outline-none focus:ring-2 focus:ring-primary/20 dark:border-gray-700 dark:bg-dark-card dark:text-gray-100 dark:placeholder:text-gray-500"
        aria-label="Tìm kiếm món ăn"
      />

      {value.length > 0 && (
        <button
          onClick={handleClear}
          className="absolute right-3 top-1/2 -translate-y-1/2 rounded-full p-0.5 text-gray-400 transition-colors hover:text-gray-600 dark:hover:text-gray-300"
          aria-label="Xóa tìm kiếm"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth="2"
            strokeLinecap="round"
            strokeLinejoin="round"
            aria-hidden="true"
          >
            <path d="M18 6L6 18M6 6l12 12" />
          </svg>
        </button>
      )}
    </div>
  );
}
