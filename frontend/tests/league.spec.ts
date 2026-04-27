import { test, expect } from '@playwright/test';

test.describe('League', () => {
  test('list leagues page loads', async ({ page }) => {
    await page.goto('/leagues');
    await expect(page.locator('h1')).toContainText('Leagues');
  });

  test('create league form renders', async ({ page }) => {
    await page.goto('/leagues/create');
    await expect(page.locator('form')).toBeVisible();
    await expect(page.locator('[name="name"]')).toBeVisible();
  });

  test('create league and view details', async ({ page }) => {
    await page.goto('/leagues/create');
    await page.fill('[name="name"]', 'Test League');
    await page.fill('[name="description"]', 'A test league');
    await page.click('[type="submit"]');
    await expect(page).toHaveURL(/\/leagues\/\d+/);
    await expect(page.locator('text=Test League')).toBeVisible();
  });
});