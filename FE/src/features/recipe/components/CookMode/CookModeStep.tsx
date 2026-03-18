'use client';

import { AnimatePresence, motion } from 'framer-motion';
import { useCookMode } from './CookModeRoot';

const slideVariants = {
  enter: (direction: number) => ({
    x: direction > 0 ? 300 : -300,
    opacity: 0,
  }),
  center: {
    x: 0,
    opacity: 1,
  },
  exit: (direction: number) => ({
    x: direction > 0 ? -300 : 300,
    opacity: 0,
  }),
};

export function CookModeStep() {
  const { currentStep, totalSteps, steps } = useCookMode();
  const step = steps[currentStep];

  if (!step) return null;

  return (
    <div className="flex flex-1 flex-col items-center justify-center px-6">
      <p className="mb-6 text-sm font-medium uppercase tracking-wider text-gray-400">
        Bước {currentStep + 1}/{totalSteps}
      </p>

      <AnimatePresence mode="wait" custom={1}>
        <motion.div
          key={currentStep}
          custom={1}
          variants={slideVariants}
          initial="enter"
          animate="center"
          exit="exit"
          transition={{ duration: 0.3, ease: 'easeInOut' }}
          className="max-w-lg text-center"
        >
          <p className="text-2xl font-medium leading-relaxed text-white md:text-3xl">
            {step.instruction}
          </p>

          {step.tip && (
            <p className="mt-6 rounded-xl bg-amber-500/10 px-4 py-3 text-base text-amber-300">
              <span aria-hidden="true">&#x1F4A1; </span>
              Mẹo: {step.tip}
            </p>
          )}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}
