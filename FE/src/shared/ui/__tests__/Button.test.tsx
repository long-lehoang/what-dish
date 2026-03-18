import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from '../Button';

describe('Button', () => {
  it('renders children text', () => {
    render(<Button>Click me</Button>);
    expect(screen.getByRole('button', { name: 'Click me' })).toBeInTheDocument();
  });

  it('handles click events', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();

    render(<Button onClick={onClick}>Click</Button>);
    await user.click(screen.getByRole('button'));

    expect(onClick).toHaveBeenCalledTimes(1);
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button disabled>Disabled</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('is disabled when isLoading is true', () => {
    render(<Button isLoading>Loading</Button>);
    expect(screen.getByRole('button')).toBeDisabled();
  });

  it('shows spinner when loading', () => {
    render(<Button isLoading>Loading</Button>);
    const svg = screen.getByRole('button').querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('does not show spinner when not loading', () => {
    render(<Button>Normal</Button>);
    const svg = screen.getByRole('button').querySelector('svg');
    expect(svg).not.toBeInTheDocument();
  });

  it('does not fire click when disabled', async () => {
    const user = userEvent.setup();
    const onClick = vi.fn();

    render(
      <Button disabled onClick={onClick}>
        Nope
      </Button>,
    );
    await user.click(screen.getByRole('button'));

    expect(onClick).not.toHaveBeenCalled();
  });

  it('applies custom className', () => {
    render(<Button className="custom-class">Styled</Button>);
    expect(screen.getByRole('button')).toHaveClass('custom-class');
  });

  it('forwards ref', () => {
    const ref = vi.fn();
    render(<Button ref={ref}>Ref</Button>);
    expect(ref).toHaveBeenCalledWith(expect.any(HTMLButtonElement));
  });
});
