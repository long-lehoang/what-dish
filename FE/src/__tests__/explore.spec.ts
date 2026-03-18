import { test, expect } from '@playwright/test';

test.describe('Explore page', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/explore');
  });

  test('loads with title and search input', async ({ page }) => {
    await expect(page.getByRole('heading', { name: /khám phá/i })).toBeVisible();
    const searchInput = page.getByRole('textbox').first();
    await expect(searchInput).toBeVisible();
  });

  test('displays dish cards from API', async ({ page }) => {
    // Wait for dishes to load (from API or mock fallback)
    const dishCards = page.getByRole('link', { name: /xem /i });
    await expect(dishCards.first()).toBeVisible({ timeout: 10000 });

    const count = await dishCards.count();
    expect(count).toBeGreaterThan(0);
  });

  test('search filters dishes by name', async ({ page }) => {
    // Wait for initial dishes to load
    const dishCards = page.getByRole('link', { name: /xem /i });
    await expect(dishCards.first()).toBeVisible({ timeout: 10000 });

    // Type a search query
    const searchInput = page.getByRole('textbox').first();
    await searchInput.fill('phở');

    // Wait for debounced search + API response by checking results update
    await expect(page.locator('body')).toContainText(/phở/i, { timeout: 10000 });
  });

  test('filter button opens filter sheet', async ({ page }) => {
    const filterBtn = page.getByRole('button', { name: /bộ lọc/i });
    await expect(filterBtn).toBeVisible();

    await filterBtn.click();

    // Filter sheet should show difficulty options (radio inputs inside labels)
    await expect(page.getByLabel('Dễ')).toBeVisible();
    await expect(page.getByLabel('Trung bình')).toBeVisible();
    await expect(page.getByLabel('Khó')).toBeVisible();

    // Filter sheet should show cook time options
    await expect(page.getByRole('button', { name: '< 30 phút' })).toBeVisible();

    // Should have apply and clear buttons
    await expect(page.getByRole('button', { name: /áp dụng/i })).toBeVisible();
    await expect(page.getByRole('button', { name: /xóa bộ lọc/i })).toBeVisible();
  });

  test('applying a filter shows active filter chip', async ({ page }) => {
    // Open filter sheet
    const filterBtn = page.getByRole('button', { name: /bộ lọc/i });
    await filterBtn.click();

    // Wait for sheet to be visible, then select difficulty
    const easyRadio = page.getByLabel('Dễ');
    await expect(easyRadio).toBeVisible();
    await easyRadio.click();

    // Apply
    const applyBtn = page.getByRole('button', { name: /áp dụng/i });
    await applyBtn.click();

    // Active filter chip should appear (chip has role="listitem" via explicit ARIA)
    await expect(page.locator('[aria-label*="Xóa bộ lọc: Dễ"]')).toBeVisible({
      timeout: 5000,
    });
  });

  test('applying difficulty filter returns matching dishes', async ({ page }) => {
    // Wait for initial dishes to load
    const dishCards = page.getByRole('link', { name: /xem /i });
    await expect(dishCards.first()).toBeVisible({ timeout: 10000 });

    // Open filter, select "Dễ" difficulty, apply
    await page.getByRole('button', { name: /bộ lọc/i }).click();
    await expect(page.getByLabel('Dễ')).toBeVisible();
    await page.getByLabel('Dễ').click();
    await page.getByRole('button', { name: /áp dụng/i }).click();

    // Should still show dish cards (filter returns results, not empty)
    const filteredCards = page.getByRole('link', { name: /xem /i });
    await expect(filteredCards.first()).toBeVisible({ timeout: 10000 });
    const count = await filteredCards.count();
    expect(count).toBeGreaterThan(0);
  });

  test('dish cards link to detail pages', async ({ page }) => {
    const dishCards = page.getByRole('link', { name: /xem /i });
    await expect(dishCards.first()).toBeVisible({ timeout: 10000 });

    const href = await dishCards.first().getAttribute('href');
    expect(href).toMatch(/^\/dish\//);
  });
});
