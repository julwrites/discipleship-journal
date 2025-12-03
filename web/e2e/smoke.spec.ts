import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Discipleship Journal/);
});

test('login page elements are present', async ({ page }) => {
  await page.goto('/login');

  // Check that the heading is visible eventually.
  const heading = page.getByRole('heading', { name: /Discipleship Journal/i });
  await expect(heading).toBeVisible({ timeout: 10000 });
});
