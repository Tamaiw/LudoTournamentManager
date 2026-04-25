import { test, expect } from '@playwright/test';

test.describe('Smoke Tests', () => {
  test('full user journey: login -> create league -> create tournament', async ({ page }) => {
    await page.goto('/login');
    await page.fill('[name="email"]', 'admin@example.com');
    await page.fill('[name="password"]', 'adminpassword');
    await page.click('[type="submit"]');
    await expect(page).toHaveURL('/dashboard');

    await page.click('[href="/leagues/create"]');
    await page.fill('[name="name"]', 'Smoke Test League');
    await page.click('[type="submit"]');
    await expect(page).toHaveURL(/\/leagues\/\d+/);

    await page.click('[href="/tournaments/create"]');
    await page.fill('[name="name"]', 'Smoke Test Tournament');
    await page.click('[type="submit"]');
    await expect(page).toHaveURL(/\/tournaments\/\d+/);
  });
});