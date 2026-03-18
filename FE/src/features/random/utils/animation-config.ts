import type { Variants } from 'framer-motion';

export const PHASE_DURATIONS = {
  shuffle: 800,
  converge: 500,
  select: 600,
  reveal: 500,
  settle: 400,
} as const;

export const SHUFFLE_VARIANTS: Variants = {
  initial: {
    x: 0,
    y: 0,
    rotate: 0,
    opacity: 0,
  },
  animate: (i: number) => ({
    x: Math.cos(i * 1.8) * 120 + (Math.random() - 0.5) * 60,
    y: Math.sin(i * 1.4) * 100 + (Math.random() - 0.5) * 40,
    rotate: (Math.random() - 0.5) * 40,
    opacity: 1,
    transition: {
      duration: 0.8,
      ease: 'easeInOut',
      delay: i * 0.03,
    },
  }),
};

export const CONVERGE_VARIANTS: Variants = {
  initial: (custom: { x: number; y: number; rotate: number }) => ({
    x: custom.x,
    y: custom.y,
    rotate: custom.rotate,
    opacity: 1,
  }),
  animate: (i: number) => ({
    x: 0,
    y: i * -2,
    rotate: (i - 6) * 1.5,
    opacity: 1,
    transition: {
      duration: 0.5,
      ease: [0.25, 0.46, 0.45, 0.94],
      delay: i * 0.02,
    },
  }),
};

export const SELECT_VARIANTS: Variants = {
  initial: {
    scale: 1,
    rotateY: 0,
    opacity: 1,
  },
  animate: {
    scale: 1.3,
    rotateY: 180,
    opacity: 1,
    transition: {
      duration: 0.6,
      ease: [0.34, 1.56, 0.64, 1],
    },
  },
};

export const REVEAL_VARIANTS: Variants = {
  initial: {
    scale: 1.3,
    opacity: 1,
    boxShadow: '0 0 0px rgba(255, 107, 53, 0)',
  },
  animate: {
    scale: 1.1,
    opacity: 1,
    boxShadow: '0 0 40px rgba(255, 107, 53, 0.6)',
    transition: {
      duration: 0.5,
      ease: 'easeOut',
    },
  },
};

export const SETTLE_VARIANTS: Variants = {
  initial: {
    opacity: 1,
    scale: 1,
  },
  animate: (isSelected: boolean) => ({
    opacity: isSelected ? 1 : 0.2,
    scale: isSelected ? 1 : 0.85,
    transition: {
      duration: 0.4,
      ease: 'easeOut',
    },
  }),
};
