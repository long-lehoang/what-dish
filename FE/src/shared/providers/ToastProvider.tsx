'use client';

import { createContext, useCallback, useContext, useRef, useState, type ReactNode } from 'react';
import { AnimatePresence } from 'framer-motion';
import { Toast, type ToastData, type ToastVariant } from '@shared/ui/Toast';

interface ToastOptions {
  variant?: ToastVariant;
  duration?: number;
}

interface ToastContextValue {
  toast: (message: string, options?: ToastOptions) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

const DEFAULT_DURATION = 3000;

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<ToastData[]>([]);
  const counterRef = useRef(0);

  const dismissToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const showToast = useCallback(
    (message: string, options?: ToastOptions) => {
      const id = `toast-${++counterRef.current}`;
      const variant = options?.variant ?? 'info';
      const duration = options?.duration ?? DEFAULT_DURATION;

      const newToast: ToastData = { id, message, variant };
      setToasts((prev) => [...prev, newToast]);

      setTimeout(() => {
        dismissToast(id);
      }, duration);
    },
    [dismissToast],
  );

  return (
    <ToastContext.Provider value={{ toast: showToast }}>
      {children}
      <div
        className="pointer-events-none fixed inset-x-0 top-4 z-[100] flex flex-col items-center gap-2 px-4"
        aria-live="polite"
      >
        <AnimatePresence mode="popLayout">
          {toasts.map((t) => (
            <div key={t.id} className="pointer-events-auto w-full max-w-sm">
              <Toast toast={t} onDismiss={dismissToast} />
            </div>
          ))}
        </AnimatePresence>
      </div>
    </ToastContext.Provider>
  );
}

export function useToast(): ToastContextValue {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
}
