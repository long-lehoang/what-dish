'use client';

import { cn } from '@shared/lib/utils';
import { Button } from '@shared/ui';
import { useTimer } from '../../hooks/useTimer';

interface CookModeTimerProps {
  initialSeconds: number;
  className?: string;
}

const CIRCLE_SIZE = 160;
const STROKE_WIDTH = 8;
const RADIUS = (CIRCLE_SIZE - STROKE_WIDTH) / 2;
const CIRCUMFERENCE = 2 * Math.PI * RADIUS;

export function CookModeTimer({ initialSeconds, className }: CookModeTimerProps) {
  const { seconds, isRunning, isPaused, start, pause, resume, reset, formattedTime } =
    useTimer(initialSeconds);

  const progress = initialSeconds > 0 ? seconds / initialSeconds : 0;
  const dashOffset = CIRCUMFERENCE * (1 - progress);

  return (
    <div className={cn('flex flex-col items-center gap-4', className)}>
      <div className="relative" style={{ width: CIRCLE_SIZE, height: CIRCLE_SIZE }}>
        <svg width={CIRCLE_SIZE} height={CIRCLE_SIZE} className="-rotate-90" aria-hidden="true">
          {/* Background circle */}
          <circle
            cx={CIRCLE_SIZE / 2}
            cy={CIRCLE_SIZE / 2}
            r={RADIUS}
            fill="none"
            stroke="currentColor"
            strokeWidth={STROKE_WIDTH}
            className="text-gray-700"
          />
          {/* Progress circle */}
          <circle
            cx={CIRCLE_SIZE / 2}
            cy={CIRCLE_SIZE / 2}
            r={RADIUS}
            fill="none"
            stroke="currentColor"
            strokeWidth={STROKE_WIDTH}
            strokeLinecap="round"
            strokeDasharray={CIRCUMFERENCE}
            strokeDashoffset={dashOffset}
            className={cn(
              'transition-[stroke-dashoffset] duration-1000',
              seconds <= 10 ? 'text-red-500' : 'text-primary',
            )}
          />
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span
            className="text-3xl font-bold tabular-nums text-white"
            aria-live="polite"
            aria-label={`Còn lại ${formattedTime}`}
          >
            {formattedTime}
          </span>
        </div>
      </div>

      <div className="flex items-center gap-2">
        {!isRunning && !isPaused && (
          <Button variant="primary" size="sm" onClick={start}>
            Bắt đầu
          </Button>
        )}
        {isRunning && !isPaused && (
          <Button variant="secondary" size="sm" onClick={pause}>
            Tạm dừng
          </Button>
        )}
        {isPaused && (
          <Button variant="primary" size="sm" onClick={resume}>
            Tiếp tục
          </Button>
        )}
        {(isRunning || isPaused) && (
          <Button variant="ghost" size="sm" onClick={reset}>
            Đặt lại
          </Button>
        )}
      </div>
    </div>
  );
}
