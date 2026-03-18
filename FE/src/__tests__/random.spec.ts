import { test, expect } from '@playwright/test';

test.describe('Random page', () => {
  test('loads the random page', async ({ page }) => {
    await page.goto('/random');
    await expect(page).toHaveURL('/random');
  });

  test('has a shuffle/start button', async ({ page }) => {
    await page.goto('/random');
    const button = page.getByRole('button').first();
    await expect(button).toBeVisible();
  });

  test('can navigate to random from homepage', async ({ page }) => {
    await page.goto('/');
    const ctaLink = page.getByRole('link', { name: /random|ngẫu nhiên|ăn gì/i }).first();
    if (await ctaLink.isVisible()) {
      await ctaLink.click();
      await expect(page).toHaveURL('/random');
    }
  });
});

test.describe('Homepage', () => {
  test('renders the homepage', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/tối nay ăn gì/i);
  });

  test('has Vietnamese content', async ({ page }) => {
    await page.goto('/');
    // The page should contain Vietnamese text
    const body = page.locator('body');
    await expect(body).toContainText(/ăn/i);
  });

  test('has html lang="vi"', async ({ page }) => {
    await page.goto('/');
    const lang = await page.locator('html').getAttribute('lang');
    expect(lang).toBe('vi');
  });
});
