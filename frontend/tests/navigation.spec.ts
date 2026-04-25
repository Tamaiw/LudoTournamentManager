import { test, expect } from '@playwright/test';

test.describe('Navigation', () => {
  test('sidebar navigation links work', async ({ page }) => {
    await page.goto('/dashboard');
    await page.click('[href="/leagues"]');
    await expect(page).toHaveURL(/\/leagues/);
  });

  test('breadcrumb navigation', async ({ page }) => {
    await page.goto('/leagues/1');
    await expect(page.locator('[aria-label="breadcrumb"]')).toBeVisible();
  });

  test('404 page renders correctly', async ({ page }) => {
    await page.goto('/nonexistent-page');
    await expect(page.locator('text=404')).toBeVisible();
  });
});