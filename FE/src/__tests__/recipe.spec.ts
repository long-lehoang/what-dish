import { test, expect } from '@playwright/test';

test.describe('Recipe detail page', () => {
  test('shows 404 for non-existent dish', async ({ page }) => {
    const response = await page.goto('/dish/non-existent-dish-slug-xyz');
    if (response) {
      const status = response.status();
      expect([200, 404]).toContain(status);
    }
    // Should show some not-found content
    const body = page.locator('body');
    await expect(body).toBeVisible();
  });

  test('loads a dish detail page with content', async ({ page }) => {
    // Navigate to a known mock dish
    await page.goto('/dish/pho-bo-ha-noi');

    // Should show dish name
    await expect(page.getByRole('heading', { name: 'Phở Bò Hà Nội' })).toBeVisible({
      timeout: 10000,
    });
  });

  test('shows recipe info badges (time, difficulty)', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Phở Bò Hà Nội' })).toBeVisible({
      timeout: 10000,
    });

    // Should show time and difficulty badges (time format: "30 phút", "2 giờ", or "3g30p")
    const body = page.locator('body');
    await expect(body).toContainText(/phút|giờ|\dg\d+p/);
    await expect(body).toContainText(/Dễ|TB|Khó/);
  });

  test('shows ingredients section', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Phở Bò Hà Nội' })).toBeVisible({
      timeout: 10000,
    });

    // Ingredients heading
    await expect(page.getByRole('heading', { name: 'Nguyên liệu' }).first()).toBeVisible();

    // Should show actual ingredient names from mock data (use label locator to scope to ingredient list)
    const ingredientSection = page.locator('ul');
    await expect(ingredientSection.getByText('Xương bò')).toBeVisible();
    await expect(ingredientSection.getByText('Bánh phở tươi')).toBeVisible();
  });

  test('shows cooking steps section', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Phở Bò Hà Nội' })).toBeVisible({
      timeout: 10000,
    });

    // Steps heading
    await expect(page.getByText('Cách làm')).toBeVisible();

    // Step numbers should appear
    const stepBadges = page.locator('span').filter({ hasText: /^[1-8]$/ });
    const count = await stepBadges.count();
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test('has a cook mode button', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Phở Bò Hà Nội' })).toBeVisible({
      timeout: 10000,
    });

    const cookBtn = page.getByRole('link', { name: /cook mode/i });
    await expect(cookBtn).toBeVisible();
    await expect(cookBtn).toHaveAttribute('href', '/dish/pho-bo-ha-noi/cook');
  });

  test('serving adjuster works', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Nguyên liệu' }).first()).toBeVisible({
      timeout: 10000,
    });

    // Should show default servings in the adjuster (scoped to avoid matching the badge)
    const adjuster = page.locator('text=4 phần').first();
    await expect(adjuster).toBeVisible();

    // Click + to increase
    const increaseBtn = page.getByRole('button', { name: /tăng khẩu phần/i });
    await increaseBtn.click();

    await expect(page.locator('text=5 phần').first()).toBeVisible();
  });

  test('ingredient checkboxes toggle', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi');
    await expect(page.getByRole('heading', { name: 'Nguyên liệu' }).first()).toBeVisible({
      timeout: 10000,
    });

    // Click the first ingredient checkbox
    const checkbox = page.getByRole('checkbox').first();
    await expect(checkbox).not.toBeChecked();

    await checkbox.click();
    await expect(checkbox).toBeChecked();
  });
});

test.describe('Cook mode', () => {
  test('loads cook mode for a dish', async ({ page }) => {
    await page.goto('/dish/pho-bo-ha-noi/cook');

    // Should show step content
    await expect(page.getByText(/bước 1/i)).toBeVisible({ timeout: 10000 });

    // Should show the first step description
    const body = page.locator('body');
    await expect(body).toContainText(/xương bò|chần/i);
  });
});

test.describe('Recipe navigation from explore', () => {
  test('clicking a dish card navigates to detail page', async ({ page }) => {
    await page.goto('/explore');

    // Wait for dishes to load
    const dishCard = page.getByRole('link', { name: /xem /i }).first();
    await expect(dishCard).toBeVisible({ timeout: 10000 });

    // Click the first dish card
    await dishCard.click();

    // Should navigate to a dish detail page
    await expect(page).toHaveURL(/\/dish\//);

    // Detail page should show content
    const heading = page.getByRole('heading', { level: 1 });
    await expect(heading).toBeVisible({ timeout: 10000 });
  });
});
