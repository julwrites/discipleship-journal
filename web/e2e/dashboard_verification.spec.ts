import { test, expect } from '@playwright/test';

test('dashboard resources and user dropdowns', async ({ page }) => {
  page.on('console', msg => console.log('PAGE LOG:', msg.text()));

  // Mock API Responses
  await page.route('**/api/users/me', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        id: 'test-user-id',
        firebase_uid: 'test-firebase-uid',
        email: 'test@example.com',
        username: 'TestUser',
        settings: {}
      })
    });
  });

  await page.route('**/api/notes*', async route => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        data: [],
        meta: { total: 0, page: 1, limit: 20, total_pages: 1 }
      })
    });
  });

  await page.goto('/');

  // Inject User
  await page.evaluate(() => {
    console.log('Injecting E2E_TEST_USER');
    localStorage.setItem('E2E_TEST_USER', JSON.stringify({
       uid: 'test-firebase-uid',
       email: 'test@example.com',
       displayName: 'Test User'
    }));
    // We can't access import.meta.env inside evaluate directly unless we pass it, but we can assume it's there.
  });

  await page.reload();

  // Wait for load
  await page.waitForTimeout(1000);

  // Take screenshot for debugging
  await page.screenshot({ path: 'verification/debug_dashboard.png' });

  // Check if we are on login
  if (await page.getByText('Sign in to your account').isVisible()) {
      console.log('Still on Login Page');
  }

  await expect(page.getByText('My Journal')).toBeVisible();

  // 1. Verify Resources Dropdown
  const resourcesBtn = page.getByTitle('Resources');
  await expect(resourcesBtn).toBeVisible();
  await resourcesBtn.click();

  await expect(page.getByRole('menuitem', { name: 'Reading Plans' })).toBeVisible();
  await expect(page.getByRole('menuitem', { name: 'Memory Verses' })).toBeVisible();
  await expect(page.getByRole('menuitem', { name: 'Study Templates' })).toBeVisible();

  // Close menu
  await page.keyboard.press('Escape');

  // 2. Verify User Dropdown
  const userBtn = page.getByTitle('User Menu');
  await expect(userBtn).toBeVisible();
  await userBtn.click();

  await expect(page.getByRole('menuitem', { name: 'Settings' })).toBeVisible();
  await expect(page.getByRole('menuitem', { name: 'Sign Out' })).toBeVisible();

  await page.screenshot({ path: 'verification/dashboard_dropdowns.png' });
});
