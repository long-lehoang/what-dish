import { test, expect } from '@playwright/test';

test.describe('Homepage', () => {
  test('renders title and Vietnamese content', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/tối nay ăn gì/i);
    await expect(page.locator('html')).toHaveAttribute('lang', 'vi');
  });

  test('shows featured dishes section with dish cards', async ({ page }) => {
    await page.goto('/');

    // Featured section should render dish cards
    const dishCards = page.getByRole('link', { name: /xem /i });
    await expect(dishCards.first()).toBeVisible({ timeout: 10000 });

    // Should have multiple dish cards
    const count = await dishCards.count();
    expect(count).toBeGreaterThanOrEqual(3);
  });

  test('featured dish cards link to detail pages', async ({ page }) => {
    await page.goto('/');

    const firstCard = page.getByRole('link', { name: /xem /i }).first();
    await expect(firstCard).toBeVisible({ timeout: 10000 });

    const href = await firstCard.getAttribute('href');
    expect(href).toMatch(/^\/dish\//);
  });

  test('can navigate to random page from CTA', async ({ page }) => {
    await page.goto('/');

    const cta = page.getByRole('link', { name: /hôm nay ăn gì|bắt đầu ngay/i }).first();
    await expect(cta).toBeVisible();
    await cta.click();
    await expect(page).toHaveURL('/random');
  });
});

test.describe('Random page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/random');
  });

  test('loads with shuffle button and title', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /tối nay/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /bắt đầu chọn món|lật bài/i })).toBeVisible();
  });

  test('filter bar shows category chips', async ({ page }) => {
    // Wait for categories to load from API/mock
    const categoryChip = page.getByRole('button', { name: /Bún\/Phở|Cơm|Xào|Lẩu/i }).first();
    await expect(categoryChip).toBeVisible({ timeout: 10000 });

    // Verify multiple category types are shown
    await expect(page.getByRole('button', { name: 'Cơm' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Lẩu' })).toBeVisible();
  });

  test('filter bar expands to show difficulty and time options', async ({ page }) => {
    // Click the filter toggle button
    const filterToggle = page.getByRole('button', { name: /bộ lọc/i }).first();
    await expect(filterToggle).toBeVisible();
    await filterToggle.click();

    // Difficulty options should appear
    await expect(page.getByText('ĐỘ KHÓ', { exact: false })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Dễ' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Trung bình' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Khó' })).toBeVisible();

    // Cook time options should appear
    await expect(page.getByText('THỜI GIAN NẤU', { exact: false })).toBeVisible();
    await expect(page.getByRole('button', { name: '< 30 phút' })).toBeVisible();
    await expect(page.getByRole('button', { name: '< 1 giờ' })).toBeVisible();
  });

  test('dish pool shows dishes with names', async ({ page }) => {
    // Wait for the dish pool to load
    const dishCount = page.getByText(/\d+ món trong giỏ xoay/);
    await expect(dishCount).toBeVisible({ timeout: 10000 });

    // Should show actual dish cards
    const dishLinks = page.locator('section').last().getByRole('link');
    const count = await dishLinks.count();
    expect(count).toBeGreaterThan(0);
  });

  test('clicking a category filter updates dish pool', async ({ page }) => {
    // Wait for initial load
    await expect(page.getByText(/\d+ món trong giỏ xoay/)).toBeVisible({ timeout: 10000 });

    // Get initial dish count text
    const countText = await page.getByText(/\d+ món trong giỏ xoay/).textContent();
    const initialCount = parseInt(countText?.match(/(\d+)/)?.[1] ?? '0');

    // Click a specific category to filter
    const comChip = page.getByRole('button', { name: 'Cơm' });
    if (await comChip.isVisible()) {
      await comChip.click();

      // Wait for the dish pool to re-render with filtered results
      await expect(page.getByText(/\d+ món trong giỏ xoay/)).toBeVisible({ timeout: 10000 });

      // Count should change (fewer dishes)
      const newText = await page.getByText(/\d+ món trong giỏ xoay/).textContent();
      const newCount = parseInt(newText?.match(/(\d+)/)?.[1] ?? '0');
      expect(newCount).toBeLessThanOrEqual(initialCount);
    }
  });

  test('dish type filter returns matching dishes (not zero)', async ({ page }) => {
    // This test catches the bug where dish_type_id was NULL in the DB,
    // causing all category filters to return 0 results.
    await expect(page.getByText(/\d+ món trong giỏ xoay/)).toBeVisible({ timeout: 10000 });

    // Click a category chip that should have at least 1 matching dish
    const categoryChips = page.getByRole('button', { name: /Cơm|Phở|Bún|Gỏi|Bánh/i });
    const firstChip = categoryChips.first();
    if (await firstChip.isVisible()) {
      await firstChip.click();

      // Wait for the dish pool to update after filter
      await expect(page.getByText(/\d+ món trong giỏ xoay/)).toBeVisible({ timeout: 10000 });

      // The dish pool must show > 0 dishes — if it shows 0,
      // it means dish_type_id is not assigned in the DB
      const countText = await page.getByText(/\d+ món trong giỏ xoay/).textContent();
      const count = parseInt(countText?.match(/(\d+)/)?.[1] ?? '0');
      expect(count).toBeGreaterThan(0);
    }
  });

  test('clicking shuffle button starts animation', async ({ page }) => {
    const shuffleBtn = page.getByRole('button', { name: /bắt đầu chọn món|lật bài/i });
    await expect(shuffleBtn).toBeVisible();

    await shuffleBtn.click();

    // Button should become disabled during animation
    await expect(shuffleBtn).toBeDisabled({ timeout: 2000 });
  });
});
