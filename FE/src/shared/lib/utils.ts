import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]): string {
  return twMerge(clsx(inputs));
}

export function formatDuration(minutes: number): string {
  if (minutes < 60) {
    return `${minutes} phút`;
  }

  const hours = Math.floor(minutes / 60);
  const remainingMinutes = minutes % 60;

  if (remainingMinutes === 0) {
    return `${hours} giờ`;
  }

  return `${hours} giờ ${remainingMinutes} phút`;
}

export function formatCurrency(vnd: number): string {
  const thousands = Math.round(vnd / 1000);
  return `~${thousands}k`;
}

export function slugify(text: string): string {
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '') // Remove diacritics
    .replace(/đ/g, 'd')
    .replace(/Đ/g, 'd')
    .replace(/[^a-z0-9\s-]/g, '')
    .replace(/\s+/g, '-')
    .replace(/-+/g, '-')
    .replace(/^-|-$/g, '');
}

export function getRandomExcluding<T>(items: T[], excludeIds: string[], idKey: keyof T): T | null {
  const available = items.filter((item) => !excludeIds.includes(String(item[idKey])));

  if (available.length === 0) {
    return null;
  }

  const index = Math.floor(Math.random() * available.length);
  return available[index] ?? null;
}

export function isExpired(timestamp: number, days: number): boolean {
  const expiryMs = days * 24 * 60 * 60 * 1000;
  return Date.now() - timestamp > expiryMs;
}

export function generateRoomCode(): string {
  const chars = 'ABCDEFGHJKLMNPQRSTUVWXYZ23456789'; // Exclude ambiguous: 0/O, 1/I
  let code = '';
  for (let i = 0; i < 6; i++) {
    code += chars[Math.floor(Math.random() * chars.length)];
  }
  return code;
}
