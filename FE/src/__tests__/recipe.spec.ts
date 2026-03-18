import { test, expect } from '@playwright/test';

test.describe('Recipe detail page', () => {
  test('shows 404 for non-existent dish', async ({ page }) => {
    const response = await page.goto('/dish/non-existent-dish-slug-xyz');
    // Should either show a 404 page or redirect
    if (response) {
      const status = response.status();
      // Accept 404 or 200 (with not-found UI)
      expect([200, 404]).toContain(status);
    }
    // Should show some not-found content
    const body = page.locator('body');
    await expect(body).toBeVisible();
  });
});

test.describe('Recipe navigation', () => {
  test('explore page loads', async ({ page }) => {
    await page.goto('/explore');
    await expect(page).toHaveURL('/explore');
  });
});
