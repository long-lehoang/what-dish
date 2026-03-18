'use client';

import { motion, type MotionStyle, type Variants, type TargetAndTransition } from 'framer-motion';
import { cn } from '@shared/lib/utils';

interface MiniCardProps {
  style?: MotionStyle;
  animate?: TargetAndTransition | string;
  variants?: Variants;
  custom?: number;
  className?: string;
}

const FOOD_EMOJIS = [
  '\uD83C\uDF5C',
  '\uD83C\uDF5B',
  '\uD83C\uDF72',
  '\uD83C\uDF63',
  '\uD83C\uDF5D',
  '\uD83C\uDF5E',
  '\uD83E\uDD5F',
  '\uD83E\uDD62',
  '\uD83C\uDF54',
  '\uD83C\uDF5F',
  '\uD83C\uDF69',
  '\uD83E\uDD58',
];

export function MiniCard({ style, animate, variants, custom, className }: MiniCardProps) {
  const emoji =
    FOOD_EMOJIS[
      typeof custom === 'number'
        ? custom % FOOD_EMOJIS.length
        : Math.floor(Math.random() * FOOD_EMOJIS.length)
    ];

  return (
    <motion.div
      style={style}
      animate={animate}
      variants={variants}
      custom={custom}
      className={cn(
        'flex h-20 w-[60px] items-center justify-center rounded-lg shadow-md',
        'bg-gradient-to-br from-primary to-secondary',
        'select-none',
        className,
      )}
      aria-hidden="true"
    >
      <div className="flex flex-col items-center gap-0.5">
        <span className="text-xl">{emoji}</span>
        <div className="h-px w-6 bg-white/30" />
        <div className="h-px w-4 bg-white/20" />
      </div>
    </motion.div>
  );
}
