'use client';

import { cn } from '@shared/lib/utils';
import { useCookMode } from './CookModeRoot';

export function CookModeControls() {
  const { currentStep, totalSteps, next, prev, exit } = useCookMode();

  const isFirst = currentStep === 0;
  const isLast = currentStep === totalSteps - 1;

  return (
    <nav
      className="flex items-center justify-between border-t border-gray-800 px-4 py-4"
      aria-label="Điều khiển nấu ăn"
    >
      {/* Previous */}
      <button
        onClick={prev}
        disabled={isFirst}
        className={cn(
          'flex h-12 w-12 items-center justify-center rounded-full transition-colors',
          isFirst
            ? 'cursor-not-allowed text-gray-600'
            : 'bg-gray-800 text-white hover:bg-gray-700 active:bg-gray-600',
        )}
        aria-label="Bước trước"
      >
        <svg
          width="20"
          height="20"
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
          aria-hidden="true"
        >
          <path d="M19 12H5M12 19l-7-7 7-7" />
        </svg>
      </button>

      {/* Progress dots */}
      <div className="flex items-center gap-1.5" role="group" aria-label="Tiến độ">
        {Array.from({ length: totalSteps }, (_, i) => (
          <div
            key={i}
            className={cn(
              'h-2 rounded-full transition-all duration-300',
              i === currentStep
                ? 'w-6 bg-primary'
                : i < currentStep
                  ? 'w-2 bg-primary/50'
                  : 'w-2 bg-gray-600',
            )}
            aria-label={`Bước ${i + 1}${i === currentStep ? ' (hiện tại)' : ''}`}
          />
        ))}
      </div>

      {/* Next / Complete / Exit */}
      <div className="flex items-center gap-2">
        {isLast ? (
          <button
            onClick={exit}
            className="flex h-12 items-center justify-center rounded-full bg-green-600 px-5 text-sm font-medium text-white transition-colors hover:bg-green-500 active:bg-green-700"
          >
            Hoàn thành
          </button>
        ) : (
          <button
            onClick={next}
            className="flex h-12 w-12 items-center justify-center rounded-full bg-primary text-white transition-colors hover:bg-primary/90 active:bg-primary/80"
            aria-label="Bước tiếp"
          >
            <svg
              width="20"
              height="20"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              aria-hidden="true"
            >
              <path d="M5 12h14M12 5l7 7-7 7" />
            </svg>
          </button>
        )}

        {!isLast && (
          <button
            onClick={exit}
            className="flex h-12 w-12 items-center justify-center rounded-full text-gray-400 transition-colors hover:bg-gray-800 hover:text-white"
            aria-label="Thoát chế độ nấu"
          >
            <svg
              width="20"
              height="20"
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
    </nav>
  );
}
