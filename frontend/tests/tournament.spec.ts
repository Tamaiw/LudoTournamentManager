import { test, expect } from '@playwright/test';

test.describe('Tournament', () => {
  test('create tournament form renders', async ({ page }) => {
    await page.goto('/tournaments/create');
    await expect(page.locator('form')).toBeVisible();
  });

  test('view tournament details', async ({ page }) => {
    await page.goto('/tournaments/1');
    await expect(page.locator('text=Tournament')).toBeVisible();
  });
});