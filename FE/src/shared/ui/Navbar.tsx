'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { motion } from 'framer-motion';
import { useTheme } from '@shared/providers/ThemeProvider';
import { cn } from '@shared/lib/utils';

const NAV_LINKS = [
  { href: '/random', label: 'Random', icon: '🎲' },
  { href: '/explore', label: 'Khám phá', icon: '🔍' },
  { href: '/vote', label: 'Vote', icon: '🗳️' },
];

function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="relative flex h-9 w-9 items-center justify-center rounded-full bg-gray-100 transition-colors hover:bg-gray-200 dark:bg-dark-card dark:hover:bg-gray-700"
      aria-label={theme === 'dark' ? 'Chuyển sang chế độ sáng' : 'Chuyển sang chế độ tối'}
    >
      <motion.span
        key={theme}
        initial={{ scale: 0, rotate: -90 }}
        animate={{ scale: 1, rotate: 0 }}
        transition={{ type: 'spring', stiffness: 300, damping: 20 }}
        className="text-base"
      >
        {theme === 'dark' ? '☀️' : '🌙'}
      </motion.span>
    </button>
  );
}

export function Navbar() {
  const pathname = usePathname();

  return (
    <header className="sticky top-0 z-30 border-b border-gray-200/50 bg-background/80 backdrop-blur-lg dark:border-gray-800/50 dark:bg-dark-bg/80">
      <nav className="mx-auto flex h-14 max-w-6xl items-center justify-between px-4">
        {/* Logo / Home link */}
        <Link
          href="/"
          className="flex items-center gap-2 font-heading text-lg font-bold text-gray-900 transition-colors hover:text-primary dark:text-white dark:hover:text-primary"
        >
          <span className="text-xl">🍜</span>
          <span className="hidden sm:inline">Ăn Gì</span>
        </Link>

        {/* Nav links */}
        <div className="flex items-center gap-1">
          {NAV_LINKS.map((link) => {
            const isActive = pathname.startsWith(link.href);
            return (
              <Link
                key={link.href}
                href={link.href}
                className={cn(
                  'relative rounded-full px-3 py-1.5 text-sm font-medium transition-colors',
                  isActive
                    ? 'text-primary'
                    : 'text-gray-500 hover:text-gray-900 dark:text-gray-400 dark:hover:text-gray-100',
                )}
              >
                {isActive && (
                  <motion.span
                    layoutId="nav-active"
                    className="absolute inset-0 rounded-full bg-primary/10 dark:bg-primary/20"
                    transition={{ type: 'spring', stiffness: 350, damping: 30 }}
                  />
                )}
                <span className="relative flex items-center gap-1.5">
                  <span className="text-sm md:text-base">{link.icon}</span>
                  <span className="hidden md:inline">{link.label}</span>
                </span>
              </Link>
            );
          })}

          {/* Divider */}
          <div className="mx-1.5 h-5 w-px bg-gray-200 dark:bg-gray-700" />

          {/* Theme toggle */}
          <ThemeToggle />
        </div>
      </nav>
    </header>
  );
}
