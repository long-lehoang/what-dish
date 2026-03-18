'use client';

import Link from 'next/link';
import { motion, useInView } from 'framer-motion';
import { useRef } from 'react';
import { DishCard } from '@features/dish';
import { MOCK_DISHES } from '@shared/lib/mock-data';

// ---- Data ----

const FLOATING_FOODS = ['🍜', '🍲', '🥢', '🍛', '🥗', '🍚', '🍖', '🌶️'];
const FEATURED_DISHES = MOCK_DISHES.slice(0, 6);

const STEPS = [
  { icon: '🎴', title: 'Lật bài', desc: 'Bấm nút, nhận ngay gợi ý món ăn ngẫu nhiên' },
  { icon: '📖', title: 'Xem công thức', desc: 'Công thức chi tiết, nấu theo từng bước' },
  { icon: '🗳️', title: 'Vote nhóm', desc: 'Cùng bạn bè bình chọn món ăn yêu thích' },
];

// ---- Animation variants ----

const staggerContainer = {
  hidden: {},
  visible: {
    transition: {
      staggerChildren: 0.12,
      delayChildren: 0.3,
    },
  },
};

const staggerItem = {
  hidden: { opacity: 0, y: 24 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { type: 'spring', stiffness: 260, damping: 20 },
  },
};

// ---- Sub-components ----

function FloatingFoodEmoji({ emoji, index }: { emoji: string; index: number }) {
  const positions = [
    { left: '8%', top: '18%' },
    { left: '85%', top: '12%' },
    { left: '15%', top: '65%' },
    { left: '78%', top: '70%' },
    { left: '50%', top: '8%' },
    { left: '92%', top: '45%' },
    { left: '5%', top: '40%' },
    { left: '65%', top: '80%' },
  ];
  const pos = positions[index % positions.length]!;

  return (
    <motion.span
      className="pointer-events-none absolute select-none text-2xl opacity-30 md:text-4xl md:opacity-40"
      style={{ left: pos.left, top: pos.top }}
      animate={{
        y: [0, -12, 0, 8, 0],
        rotate: [0, 6, -6, 3, 0],
        scale: [1, 1.08, 1, 0.95, 1],
      }}
      transition={{
        duration: 5 + index * 0.7,
        repeat: Infinity,
        ease: 'easeInOut',
        delay: index * 0.4,
      }}
      aria-hidden="true"
    >
      {emoji}
    </motion.span>
  );
}

function AnimatedSection({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  const ref = useRef(null);
  const isInView = useInView(ref, { once: true, margin: '-80px' });

  return (
    <motion.section
      ref={ref}
      initial="hidden"
      animate={isInView ? 'visible' : 'hidden'}
      variants={staggerContainer}
      className={className}
    >
      {children}
    </motion.section>
  );
}

// ---- Page ----

export default function HomePage() {
  return (
    <main className="overflow-hidden">
      {/* ===== Hero ===== */}
      <section className="relative flex min-h-[90vh] flex-col items-center justify-center px-4 text-center">
        {/* Gradient bg */}
        <div className="absolute inset-0 bg-gradient-to-b from-primary/5 via-background to-background dark:from-primary/10 dark:via-dark-bg dark:to-dark-bg" />
        {/* Radial glow */}
        <div className="absolute left-1/2 top-1/3 h-[500px] w-[500px] -translate-x-1/2 -translate-y-1/2 rounded-full bg-primary/10 blur-3xl dark:bg-primary/5" />

        {/* Floating food emojis */}
        <div className="absolute inset-0 overflow-hidden">
          {FLOATING_FOODS.map((emoji, i) => (
            <FloatingFoodEmoji key={i} emoji={emoji} index={i} />
          ))}
        </div>

        {/* Hero content */}
        <motion.div
          className="relative z-10 mx-auto max-w-2xl"
          initial="hidden"
          animate="visible"
          variants={staggerContainer}
        >
          {/* Animated bowl */}
          <motion.div className="mb-6" variants={staggerItem}>
            <motion.span
              className="inline-block text-7xl md:text-8xl"
              animate={{ rotate: [0, -8, 8, -4, 0] }}
              transition={{ duration: 2.5, repeat: Infinity, repeatDelay: 3, ease: 'easeInOut' }}
            >
              🍜
            </motion.span>
          </motion.div>

          {/* Title with gradient */}
          <motion.h1
            className="mb-4 font-heading text-5xl font-bold tracking-tight text-gray-900 dark:text-white md:text-7xl"
            variants={staggerItem}
          >
            <span className="animate-shimmer bg-gradient-to-r from-primary via-secondary to-primary bg-[length:200%_auto] bg-clip-text text-transparent">
              Tối Nay
            </span>{' '}
            Ăn Gì?
          </motion.h1>

          {/* Tagline */}
          <motion.p
            className="mb-10 font-heading text-lg text-gray-500 dark:text-gray-400 md:text-xl"
            variants={staggerItem}
          >
            Hết phân vân — lật là ăn!
          </motion.p>

          {/* CTA button */}
          <motion.div variants={staggerItem}>
            <Link
              href="/random"
              className="group relative inline-flex items-center gap-2.5 overflow-hidden rounded-full bg-gradient-to-r from-primary to-secondary px-10 py-4 font-heading text-lg font-bold text-white shadow-xl transition-all duration-300 hover:shadow-2xl hover:shadow-primary/25 active:scale-95"
            >
              <span className="relative z-10">Hôm nay ăn gì?</span>
              <motion.span
                className="relative z-10 inline-block text-xl"
                animate={{ rotateY: [0, 180, 360] }}
                transition={{ duration: 2, repeat: Infinity, repeatDelay: 3 }}
              >
                🎴
              </motion.span>
              {/* Reverse gradient on hover */}
              <span className="absolute inset-0 bg-gradient-to-r from-secondary to-primary opacity-0 transition-opacity duration-300 group-hover:opacity-100" />
            </Link>
          </motion.div>

          {/* Scroll hint */}
          <motion.div
            className="mt-20"
            animate={{ y: [0, 10, 0] }}
            transition={{ duration: 2, repeat: Infinity, ease: 'easeInOut' }}
          >
            <span className="text-sm text-gray-400 dark:text-gray-500">Kéo xuống khám phá</span>
            <div className="mt-1 text-gray-300 dark:text-gray-600">↓</div>
          </motion.div>
        </motion.div>
      </section>

      {/* ===== Featured Dishes ===== */}
      <AnimatedSection className="px-4 py-20">
        <div className="mx-auto max-w-6xl">
          <motion.div className="mb-10 text-center" variants={staggerItem}>
            <span className="mb-3 inline-block rounded-full bg-primary/10 px-4 py-1 text-sm font-medium text-primary">
              Gợi ý hôm nay
            </span>
            <h2 className="font-heading text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">
              Món ngon đang chờ bạn
            </h2>
          </motion.div>

          <motion.div
            className="scrollbar-hide flex snap-x snap-mandatory gap-5 overflow-x-auto pb-4 md:grid md:grid-cols-3 md:overflow-visible"
            variants={staggerContainer}
          >
            {FEATURED_DISHES.map((dish) => (
              <motion.div
                key={dish.id}
                variants={staggerItem}
                className="w-[280px] flex-shrink-0 snap-start md:w-auto"
              >
                <DishCard dish={dish} variant="overlay" />
              </motion.div>
            ))}
          </motion.div>

          <motion.div className="mt-8 text-center" variants={staggerItem}>
            <Link
              href="/explore"
              className="inline-flex items-center gap-1 text-sm font-semibold text-primary transition-colors hover:text-secondary"
            >
              Xem tất cả món →
            </Link>
          </motion.div>
        </div>
      </AnimatedSection>

      {/* ===== How It Works ===== */}
      <AnimatedSection className="bg-white/60 px-4 py-20 dark:bg-dark-card/30">
        <div className="mx-auto max-w-4xl">
          <motion.div className="mb-12 text-center" variants={staggerItem}>
            <span className="mb-3 inline-block rounded-full bg-accent/15 px-4 py-1 text-sm font-medium text-accent">
              Đơn giản
            </span>
            <h2 className="font-heading text-3xl font-bold text-gray-900 dark:text-white md:text-4xl">
              3 bước, hết phân vân
            </h2>
          </motion.div>

          <div className="grid gap-8 md:grid-cols-3">
            {STEPS.map((step, i) => (
              <motion.div
                key={i}
                variants={staggerItem}
                whileHover={{ y: -6 }}
                transition={{ type: 'spring', stiffness: 300, damping: 20 }}
                className="relative rounded-2xl bg-white p-8 text-center shadow-sm dark:bg-dark-card"
              >
                {/* Step number pill */}
                <div className="absolute -top-3 left-1/2 -translate-x-1/2 rounded-full bg-gradient-to-r from-primary to-secondary px-3.5 py-1 text-xs font-bold text-white shadow-sm">
                  {i + 1}
                </div>
                {/* Connector line (between cards on desktop) */}
                {i < STEPS.length - 1 && (
                  <div className="absolute -right-4 top-1/2 hidden h-0.5 w-8 bg-gradient-to-r from-primary/30 to-transparent md:block" />
                )}
                <motion.div
                  className="mb-4 text-5xl"
                  whileHover={{ scale: 1.15, rotate: [0, -8, 8, 0] }}
                  transition={{ duration: 0.4 }}
                >
                  {step.icon}
                </motion.div>
                <h3 className="mb-2 font-heading text-lg font-bold text-gray-900 dark:text-white">
                  {step.title}
                </h3>
                <p className="text-sm leading-relaxed text-gray-500 dark:text-gray-400">
                  {step.desc}
                </p>
              </motion.div>
            ))}
          </div>
        </div>
      </AnimatedSection>

      {/* ===== Feature Cards ===== */}
      <AnimatedSection className="px-4 py-20">
        <div className="mx-auto grid max-w-4xl gap-6 md:grid-cols-2">
          {/* Explore card */}
          <motion.div variants={staggerItem}>
            <Link
              href="/explore"
              className="group relative block overflow-hidden rounded-3xl bg-gradient-to-br from-accent/20 via-primary/5 to-transparent p-8 transition-all duration-300 hover:shadow-xl dark:from-accent/10 dark:via-primary/5"
            >
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-accent/10 blur-2xl transition-all duration-500 group-hover:scale-150 group-hover:bg-accent/20" />
              <motion.div
                className="relative mb-5 text-5xl"
                whileHover={{ rotate: [0, -15, 15, 0], scale: 1.1 }}
                transition={{ duration: 0.5 }}
              >
                🔍
              </motion.div>
              <h3 className="relative mb-2 font-heading text-xl font-bold text-gray-900 dark:text-white">
                Khám phá món
              </h3>
              <p className="relative mb-4 text-sm leading-relaxed text-gray-600 dark:text-gray-400">
                Duyệt hàng trăm công thức, tìm kiếm theo nguyên liệu, loại món, hoặc độ khó.
              </p>
              <span className="relative inline-flex items-center gap-1 text-sm font-semibold text-primary transition-transform duration-300 group-hover:translate-x-2">
                Khám phá ngay
                <span className="transition-transform duration-300 group-hover:translate-x-1">
                  →
                </span>
              </span>
            </Link>
          </motion.div>

          {/* Vote card */}
          <motion.div variants={staggerItem}>
            <Link
              href="/vote"
              className="group relative block overflow-hidden rounded-3xl bg-gradient-to-br from-secondary/15 via-primary/5 to-transparent p-8 transition-all duration-300 hover:shadow-xl dark:from-secondary/10 dark:via-primary/5"
            >
              <div className="absolute -right-6 -top-6 h-32 w-32 rounded-full bg-secondary/10 blur-2xl transition-all duration-500 group-hover:scale-150 group-hover:bg-secondary/20" />
              <motion.div
                className="relative mb-5 text-5xl"
                whileHover={{ rotate: [0, -15, 15, 0], scale: 1.1 }}
                transition={{ duration: 0.5 }}
              >
                🗳️
              </motion.div>
              <h3 className="relative mb-2 font-heading text-xl font-bold text-gray-900 dark:text-white">
                Vote nhóm
              </h3>
              <p className="relative mb-4 text-sm leading-relaxed text-gray-600 dark:text-gray-400">
                Tạo phòng bình chọn, mời bạn bè cùng quyết định ăn gì hôm nay.
              </p>
              <span className="relative inline-flex items-center gap-1 text-sm font-semibold text-secondary transition-transform duration-300 group-hover:translate-x-2">
                Tạo phòng vote
                <span className="transition-transform duration-300 group-hover:translate-x-1">
                  →
                </span>
              </span>
            </Link>
          </motion.div>
        </div>
      </AnimatedSection>

      {/* ===== Bottom CTA ===== */}
      <section className="px-4 pb-20 pt-4">
        <motion.div
          initial={{ opacity: 0, scale: 0.92 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true, margin: '-60px' }}
          transition={{ type: 'spring', stiffness: 200, damping: 20 }}
          className="animate-shimmer mx-auto max-w-lg overflow-hidden rounded-3xl bg-gradient-to-br from-primary via-secondary to-primary bg-[length:200%_200%] p-12 text-center text-white shadow-2xl"
        >
          <motion.span
            className="mb-3 inline-block text-5xl"
            animate={{ scale: [1, 1.15, 1] }}
            transition={{ duration: 2, repeat: Infinity, ease: 'easeInOut' }}
          >
            🍽️
          </motion.span>
          <h2 className="mb-2 font-heading text-3xl font-bold">Đói chưa?</h2>
          <p className="mb-8 text-white/75">Để Tối Nay Ăn Gì giúp bạn quyết định!</p>
          <Link
            href="/random"
            className="inline-block rounded-full bg-white px-10 py-3.5 font-heading font-bold text-primary shadow-lg transition-all duration-300 hover:scale-105 hover:shadow-xl active:scale-95"
          >
            Bắt đầu ngay!
          </Link>
        </motion.div>
      </section>
    </main>
  );
}
