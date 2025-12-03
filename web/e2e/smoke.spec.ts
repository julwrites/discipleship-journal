import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Discipleship Journal/);
});

test('loads the app', async ({ page }) => {
  await page.goto('/');

  // Basic check for title
  await expect(page).toHaveTitle(/Discipleship Journal/);
});

// Using a more generic selector to debug visibility
test('login page elements are present', async ({ page }) => {
  // Listen for console logs and errors
  page.on('console', msg => console.log('PAGE LOG:', msg.text()));
  page.on('pageerror', exception => console.log(`PAGE ERROR: "${exception}"`));

  await page.goto('/login');

  // Check if we are stuck on Loading...
  const loading = page.getByText('Loading...');
  if (await loading.isVisible()) {
    console.log("App is currently showing Loading...");
  }

  // Let's check if the error boundary is showing up
  const errorHeading = page.getByRole('heading', { name: /Something went wrong/i });

  if (await errorHeading.isVisible()) {
    console.log("Error Boundary Triggered!");
    const errorText = await page.locator('pre').textContent();
    console.log("Error details:", errorText);
  }

  await expect(errorHeading).not.toBeVisible();

  // Check that the heading is visible eventually.
  const heading = page.getByRole('heading', { name: /Discipleship Journal/i });
  await expect(heading).toBeVisible({ timeout: 10000 });
});
