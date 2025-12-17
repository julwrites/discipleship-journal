import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Discipleship Journal/);
});

test('login page elements are present', async ({ page }) => {
  await page.goto('/login');

  // Check that the heading is visible eventually.
  // Using a more relaxed selector to find the heading by its text content directly
  // This helps avoid issues with specific role accessibility in some contexts
  // Use exact: true to avoid matching the footer copyright text
  const heading = page.getByText('Discipleship Journal', { exact: true });
  await expect(heading).toBeVisible({ timeout: 10000 });
});
