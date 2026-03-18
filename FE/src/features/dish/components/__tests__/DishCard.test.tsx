import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { DishCard } from '../DishCard';
import type { Dish } from '../../types';

// Mock next/image
vi.mock('next/image', () => ({
  default: (props: Record<string, unknown>) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={props.src as string} alt={props.alt as string} data-testid="dish-image" />
  ),
}));

// Mock next/link
vi.mock('next/link', () => ({
  default: ({
    href,
    children,
    ...props
  }: {
    href: string;
    children: React.ReactNode;
    onClick?: () => void;
    className?: string;
  }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

// Mock framer-motion
vi.mock('framer-motion', () => ({
  motion: {
    div: ({ children, className }: { children: React.ReactNode; className?: string }) => (
      <div className={className}>{children}</div>
    ),
  },
}));

function createMockDish(overrides?: Partial<Dish>): Dish {
  return {
    id: 'dish-1',
    name: 'Phở Bò',
    slug: 'pho-bo',
    description: 'Vietnamese beef noodle soup',
    imageUrl: 'https://example.com/pho.jpg',
    difficulty: 'MEDIUM',
    prepTime: 30,
    cookTime: 60,
    totalTime: 90,
    servings: 4,
    status: 'PUBLISHED',
    viewCount: 0,
    favoriteCount: 0,
    createdAt: '2024-01-01',
    updatedAt: '2024-01-01',
    ...overrides,
  };
}

describe('DishCard', () => {
  it('renders dish name', () => {
    render(<DishCard dish={createMockDish()} />);
    expect(screen.getByText('Phở Bò')).toBeInTheDocument();
  });

  it('links to the dish detail page', () => {
    render(<DishCard dish={createMockDish()} />);
    const link = screen.getByRole('link');
    expect(link).toHaveAttribute('href', '/dish/pho-bo');
  });

  it('has accessible label', () => {
    render(<DishCard dish={createMockDish()} />);
    expect(screen.getByLabelText('Xem Phở Bò')).toBeInTheDocument();
  });

  it('renders image when imageUrl is provided', () => {
    render(<DishCard dish={createMockDish()} />);
    const img = screen.getByTestId('dish-image');
    expect(img).toHaveAttribute('src', 'https://example.com/pho.jpg');
  });

  it('shows placeholder when no image', () => {
    render(<DishCard dish={createMockDish({ imageUrl: undefined })} />);
    expect(screen.queryByTestId('dish-image')).not.toBeInTheDocument();
  });

  it('calls onClick handler', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();

    render(<DishCard dish={createMockDish()} onClick={onClick} />);
    await user.click(screen.getByRole('link'));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('renders compact variant', () => {
    render(<DishCard dish={createMockDish()} variant="compact" />);
    expect(screen.getByText('Phở Bò')).toBeInTheDocument();
  });

  it('renders overlay variant', () => {
    render(<DishCard dish={createMockDish()} variant="overlay" />);
    expect(screen.getByText('Phở Bò')).toBeInTheDocument();
  });
});
