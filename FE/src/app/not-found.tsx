import Link from 'next/link';

export default function NotFoundPage() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center px-4 text-center">
      <div className="mb-4 text-6xl">🍽️</div>
      <h1 className="mb-2 font-heading text-3xl font-bold">404 — Không tìm thấy</h1>
      <p className="mb-8 text-gray-600 dark:text-gray-400">
        Trang bạn tìm kiếm không tồn tại hoặc đã bị xóa.
      </p>
      <Link
        href="/"
        className="rounded-full bg-primary px-6 py-3 font-heading font-bold text-white transition-transform hover:scale-105"
      >
        Về trang chủ
      </Link>
    </main>
  );
}
