'use client';

import { forwardRef, type ButtonHTMLAttributes } from 'react';
import { cn } from '@shared/lib/utils';

type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost';
type ButtonSize = 'sm' | 'md' | 'lg';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: ButtonVariant;
  size?: ButtonSize;
  isLoading?: boolean;
}

const variantStyles: Record<ButtonVariant, string> = {
  primary: 'bg-primary text-white hover:bg-primary/90 active:bg-primary/80',
  secondary: 'bg-secondary text-white hover:bg-secondary/90 active:bg-secondary/80',
  outline:
    'border-2 border-primary text-primary hover:bg-primary/10 active:bg-primary/20 dark:border-primary dark:text-primary',
  ghost:
    'text-gray-700 hover:bg-gray-100 active:bg-gray-200 dark:text-gray-300 dark:hover:bg-gray-800 dark:active:bg-gray-700',
};

const sizeStyles: Record<ButtonSize, string> = {
  sm: 'px-3 py-1.5 text-sm rounded-lg',
  md: 'px-5 py-2.5 text-base rounded-xl',
  lg: 'px-7 py-3.5 text-lg rounded-xl',
};

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(function Button(
  { variant = 'primary', size = 'md', isLoading = false, disabled, className, children, ...props },
  ref,
) {
  const isDisabled = disabled || isLoading;

  return (
    <button
      ref={ref}
      disabled={isDisabled}
      className={cn(
        'inline-flex items-center justify-center font-medium transition-colors duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/50 focus-visible:ring-offset-2',
        variantStyles[variant],
        sizeStyles[size],
        isDisabled && 'cursor-not-allowed opacity-50',
        className,
      )}
      {...props}
    >
      {isLoading && (
        <svg
          className="mr-2 h-4 w-4 animate-spin"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          aria-hidden="true"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
          />
        </svg>
      )}
      {children}
    </button>
  );
});
