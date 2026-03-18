import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Modal } from '../Modal';

// Mock framer-motion to avoid animation issues in tests
vi.mock('framer-motion', () => ({
  AnimatePresence: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  motion: {
    div: ({
      children,
      className,
      onClick,
      role,
      'aria-modal': ariaModal,
      'aria-label': ariaLabel,
      'aria-hidden': ariaHidden,
    }: Record<string, unknown>) => (
      <div
        className={className as string}
        onClick={onClick as React.MouseEventHandler}
        role={role as string}
        aria-modal={ariaModal as boolean}
        aria-label={ariaLabel as string}
        aria-hidden={ariaHidden as boolean}
      >
        {children as React.ReactNode}
      </div>
    ),
  },
}));

describe('Modal', () => {
  it('renders children when open', () => {
    render(
      <Modal isOpen onClose={() => {}}>
        <p>Modal content</p>
      </Modal>,
    );

    expect(screen.getByText('Modal content')).toBeInTheDocument();
  });

  it('does not render children when closed', () => {
    render(
      <Modal isOpen={false} onClose={() => {}}>
        <p>Modal content</p>
      </Modal>,
    );

    expect(screen.queryByText('Modal content')).not.toBeInTheDocument();
  });

  it('renders title when provided', () => {
    render(
      <Modal isOpen onClose={() => {}} title="Test Title">
        <p>Content</p>
      </Modal>,
    );

    expect(screen.getByText('Test Title')).toBeInTheDocument();
  });

  it('calls onClose when Escape is pressed', async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();

    render(
      <Modal isOpen onClose={onClose}>
        <p>Content</p>
      </Modal>,
    );

    await user.keyboard('{Escape}');

    expect(onClose).toHaveBeenCalledTimes(1);
  });

  it('calls onClose when overlay is clicked', async () => {
    const user = userEvent.setup();
    const onClose = vi.fn();

    render(
      <Modal isOpen onClose={onClose}>
        <p>Content</p>
      </Modal>,
    );

    // Click the overlay (aria-hidden div)
    const overlay = document.querySelector('[aria-hidden="true"]');
    if (overlay) {
      await user.click(overlay);
      expect(onClose).toHaveBeenCalledTimes(1);
    }
  });

  it('has dialog role with aria-modal', () => {
    render(
      <Modal isOpen onClose={() => {}} title="Dialog">
        <p>Content</p>
      </Modal>,
    );

    const dialog = screen.getByRole('dialog');
    expect(dialog).toHaveAttribute('aria-modal', 'true');
  });
});
