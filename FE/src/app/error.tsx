'use client';

interface ErrorPageProps {
  error: Error;
  reset: () => void;
}

export default function ErrorPage({ error, reset }: ErrorPageProps) {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center px-4 text-center">
      <div className="mb-4 text-6xl">😵</div>
      <h1 className="mb-2 font-heading text-3xl font-bold">Có lỗi xảy ra</h1>
      <p className="mb-8 text-gray-600 dark:text-gray-400">
        {error.message || 'Đã xảy ra lỗi không mong muốn. Vui lòng thử lại.'}
      </p>
      <button
        onClick={reset}
        className="rounded-full bg-primary px-6 py-3 font-heading font-bold text-white transition-transform hover:scale-105"
      >
        Thử lại
      </button>
    </main>
  );
}
