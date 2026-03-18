'use client';

import { motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';

export type ToastVariant = 'success' | 'error' | 'info';

export interface ToastData {
  id: string;
  message: string;
  variant: ToastVariant;
}

interface ToastProps {
  toast: ToastData;
  onDismiss: (id: string) => void;
}

const variantStyles: Record<ToastVariant, string> = {
  success: 'bg-emerald-600 text-white',
  error: 'bg-red-600 text-white',
  info: 'bg-sky-600 text-white',
};

const variantIcons: Record<ToastVariant, string> = {
  success: '\u2713',
  error: '\u2715',
  info: '\u2139',
};

export function Toast({ toast: toastData, onDismiss }: ToastProps) {
  return (
    <motion.div
      layout
      initial={{ opacity: 0, y: -20, scale: 0.95 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: -20, scale: 0.95 }}
      transition={{ duration: 0.2, ease: 'easeOut' }}
      role="alert"
      className={cn(
        'flex items-center gap-3 rounded-xl px-4 py-3 shadow-lg',
        variantStyles[toastData.variant],
      )}
    >
      <span className="flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-white/20 text-xs font-bold">
        {variantIcons[toastData.variant]}
      </span>
      <p className="text-sm font-medium">{toastData.message}</p>
      <button
        onClick={() => onDismiss(toastData.id)}
        className="ml-auto shrink-0 rounded-lg p-1 transition-colors hover:bg-white/20"
        aria-label="Dismiss"
      >
        <svg
          className="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </motion.div>
  );
}
