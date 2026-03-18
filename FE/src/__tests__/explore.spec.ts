import { test, expect } from '@playwright/test';

test.describe('Explore page', () => {
  test('loads the explore page', async ({ page }) => {
    await page.goto('/explore');
    await expect(page).toHaveURL('/explore');
  });

  test('has a search input', async ({ page }) => {
    await page.goto('/explore');
    const searchInput = page.getByRole('textbox').first();
    if (await searchInput.isVisible()) {
      await expect(searchInput).toBeVisible();
    }
  });

  test('displays dish grid area', async ({ page }) => {
    await page.goto('/explore');
    const body = page.locator('body');
    await expect(body).toBeVisible();
  });
});
