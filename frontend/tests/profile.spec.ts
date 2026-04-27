import { test, expect } from '@playwright/test';

test.describe('Profile', () => {
  test('view own profile', async ({ page }) => {
    await page.goto('/profile');
    await expect(page.locator('[role="banner"]')).toBeVisible();
  });

  test('edit profile name', async ({ page }) => {
    await page.goto('/profile');
    await page.click('[role="button"][name="Edit"]');
    await page.fill('[name="displayName"]', 'New Name');
    await page.click('[type="submit"]');
    await expect(page.locator('text=New Name')).toBeVisible();
  });
});