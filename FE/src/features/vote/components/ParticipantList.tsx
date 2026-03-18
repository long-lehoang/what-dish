'use client';

import { motion, AnimatePresence, LayoutGroup } from 'framer-motion';
import { cn } from '@shared/lib/utils';
import type { Participant } from '../types';

interface ParticipantListProps {
  participants: Participant[];
  className?: string;
}

export function ParticipantList({ participants, className }: ParticipantListProps) {
  return (
    <LayoutGroup>
      <div
        className={cn('flex flex-wrap gap-3', className)}
        role="list"
        aria-label="Danh sách người tham gia"
      >
        <AnimatePresence>
          {participants.map((participant) => (
            <motion.div
              key={participant.name}
              layout
              initial={{ scale: 0, opacity: 0 }}
              animate={{ scale: 1, opacity: 1 }}
              exit={{ scale: 0, opacity: 0 }}
              transition={{ type: 'spring', stiffness: 300, damping: 25 }}
              className="flex flex-col items-center gap-1"
              role="listitem"
            >
              <div className="relative">
                <div
                  className="flex h-12 w-12 items-center justify-center rounded-full text-lg font-bold text-white"
                  style={{ backgroundColor: participant.avatarColor }}
                  aria-hidden="true"
                >
                  {participant.name.charAt(0).toUpperCase()}
                </div>
                {participant.hasVoted && (
                  <div className="absolute -bottom-0.5 -right-0.5 flex h-5 w-5 items-center justify-center rounded-full bg-green-500 text-white">
                    <svg
                      width="12"
                      height="12"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      strokeWidth="3"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      aria-hidden="true"
                    >
                      <path d="M20 6L9 17l-5-5" />
                    </svg>
                  </div>
                )}
              </div>
              <span className="max-w-[4rem] truncate text-xs text-gray-600 dark:text-gray-400">
                {participant.name}
              </span>
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </LayoutGroup>
  );
}
