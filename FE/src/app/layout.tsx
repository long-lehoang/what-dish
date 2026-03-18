import type { Metadata, Viewport } from 'next';
import { Be_Vietnam_Pro, Inter } from 'next/font/google';
import '@/styles/globals.css';
import { ThemeProvider } from '@/shared/providers/ThemeProvider';
import { ToastProvider } from '@/shared/providers/ToastProvider';
import { Navbar } from '@/shared/ui/Navbar';

const beVietnamPro = Be_Vietnam_Pro({
  subsets: ['vietnamese', 'latin'],
  weight: ['500', '700'],
  variable: '--font-be-vietnam-pro',
  display: 'swap',
});

const inter = Inter({
  subsets: ['vietnamese', 'latin'],
  weight: ['400', '500'],
  variable: '--font-inter',
  display: 'swap',
});

export const metadata: Metadata = {
  title: {
    default: 'Tối Nay Ăn Gì — Hết phân vân, lật là ăn!',
    template: '%s | Tối Nay Ăn Gì',
  },
  description:
    'Ứng dụng giúp bạn quyết định hôm nay ăn gì qua trải nghiệm lật bài ngẫu nhiên, công thức nấu ăn chi tiết, và bình chọn nhóm.',
  manifest: '/manifest.json',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'default',
    title: 'Ăn Gì',
  },
};

export const viewport: Viewport = {
  themeColor: '#FF6B35',
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html
      lang="vi"
      className={`${beVietnamPro.variable} ${inter.variable}`}
      suppressHydrationWarning
    >
      <body className="min-h-screen bg-background font-body text-gray-900 antialiased dark:bg-dark-bg dark:text-gray-100">
        <ThemeProvider>
          <ToastProvider>
            <Navbar />
            {children}
          </ToastProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
