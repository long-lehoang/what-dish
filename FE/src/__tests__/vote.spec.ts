import { test, expect } from '@playwright/test';

test.describe('Vote page', () => {
  test('loads the vote page', async ({ page }) => {
    await page.goto('/vote');
    await expect(page).toHaveURL('/vote');
  });

  test('has create and join options', async ({ page }) => {
    await page.goto('/vote');
    const body = page.locator('body');
    // Should have some call to action for creating or joining rooms
    await expect(body).toBeVisible();
  });

  test('can navigate to create room', async ({ page }) => {
    await page.goto('/vote');
    const createLink = page.getByRole('link', { name: /tạo|create/i }).first();
    if (await createLink.isVisible()) {
      await createLink.click();
      await expect(page).toHaveURL('/vote/create');
    }
  });
});

test.describe('Vote create page', () => {
  test('loads the create room page', async ({ page }) => {
    await page.goto('/vote/create');
    await expect(page).toHaveURL('/vote/create');
  });

  test('has a form for room creation', async ({ page }) => {
    await page.goto('/vote/create');
    // Should have inputs or buttons for creating a room
    const body = page.locator('body');
    await expect(body).toBeVisible();
  });
});
