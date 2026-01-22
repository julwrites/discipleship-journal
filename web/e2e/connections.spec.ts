import { test, expect } from '@playwright/test';

test.describe('Connections (Mocked)', () => {
  test.beforeEach(async ({ page }) => {
    // 1. Mock Auth
    await page.addInitScript(() => {
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-user-123',
            email: 'test@example.com',
            displayName: 'Test User'
        }));
    });

    // 2. Mock User Profile
    await page.route('**/api/users/me', async route => {
        await route.fulfill({
            json: {
                id: 'uuid-123',
                firebase_uid: 'test-user-123',
                username: 'testuser'
            }
        });
    });

    // 3. Mock Connections List (Initially empty)
    await page.route('**/api/connections', async route => {
         await route.fulfill({ json: [] });
    });

    // 4. Mock Pending Requests (Initially empty)
     await page.route('**/api/connections/requests', async route => {
         await route.fulfill({ json: [] });
    });
  });

  test('should search for a user and send a request', async ({ page }) => {
    await page.goto('/connections');

    // Click "Find Users" tab if it exists
    await page.getByRole('tab', { name: 'Find People' }).click();

    // Mock Search API
    await page.route('**/api/users/search?q=friend*', async route => {
        await route.fulfill({
            json: [{
                id: 'uuid-friend',
                username: 'friend',
                email: 'friend@example.com',
                avatar_url: ''
            }]
        });
    });

    // Mock Send Request API
    await page.route('**/api/connections/request', async route => {
        const body = route.request().postDataJSON();
        expect(body.receiver_id).toBe('uuid-friend');
        await route.fulfill({ status: 201, json: { success: true } });
    });

    // Type in search box
    await page.getByPlaceholder('Search by email or username...').fill('friend');

    // Wait for results (auto-search)
    await expect(page.getByText('friend@example.com')).toBeVisible();

    // Click "Connect" button
    await page.getByRole('button', { name: 'Connect' }).click();

    // Verify Toast
    await expect(page.getByText('Request sent!')).toBeVisible();
  });
});
