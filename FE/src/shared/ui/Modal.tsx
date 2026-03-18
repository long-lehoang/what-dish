'use client';

import { useEffect, type ReactNode } from 'react';
import { AnimatePresence, motion } from 'framer-motion';
import { cn } from '@shared/lib/utils';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: ReactNode;
  title?: string;
  className?: string;
}

export function Modal({ isOpen, onClose, children, title, className }: ModalProps) {
  useEffect(() => {
    function handleKeyDown(event: KeyboardEvent) {
      if (event.key === 'Escape') {
        onClose();
      }
    }

    if (isOpen) {
      document.addEventListener('keydown', handleKeyDown);
      document.body.style.overflow = 'hidden';
    }

    return () => {
      document.removeEventListener('keydown', handleKeyDown);
      document.body.style.overflow = '';
    };
  }, [isOpen, onClose]);

  return (
    <AnimatePresence>
      {isOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
          {/* Overlay */}
          <motion.div
            className="absolute inset-0 bg-black/50 backdrop-blur-sm"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={onClose}
            aria-hidden="true"
          />

          {/* Content */}
          <motion.div
            role="dialog"
            aria-modal="true"
            aria-label={title}
            className={cn(
              'relative z-10 mx-4 w-full max-w-lg rounded-2xl bg-white p-6 shadow-xl dark:bg-dark-card',
              className,
            )}
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.2, ease: 'easeOut' }}
          >
            {title && (
              <h2 className="mb-4 text-xl font-bold text-gray-900 dark:text-gray-100">{title}</h2>
            )}
            {children}
          </motion.div>
        </div>
      )}
    </AnimatePresence>
  );
}
