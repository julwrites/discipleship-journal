import { test, expect } from '@playwright/test';

test.describe('Authentication Flow (Mocked)', () => {

  test('should show login page when not authenticated', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveURL('/login');
    await expect(page.getByText('Discipleship Journal')).toBeVisible();
    await expect(page.getByRole('button', { name: /Sign in with Google/i })).toBeVisible();
  });

  test('should redirect to dashboard if authenticated via localStorage mock', async ({ page }) => {
    // Set the mock user in localStorage before navigation
    await page.addInitScript(() => {
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-user-123',
            email: 'test@example.com',
            displayName: 'Test User',
            photoURL: 'https://example.com/photo.jpg'
        }));
    });

    // Mock API requests that might be triggered on Dashboard load
    await page.route('**/api/users/me', async route => {
        await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify({
                id: 'uuid-123',
                firebase_uid: 'test-user-123',
                email: 'test@example.com',
                username: 'testuser',
                settings: { bible_version: 'ESV' }
            })
        });
    });

    await page.route('**/api/notes**', async route => {
        await route.fulfill({
            status: 200,
            contentType: 'application/json',
            body: JSON.stringify([])
        });
    });

    await page.goto('/');
    await expect(page).toHaveURL('/'); // Should stay on root (Dashboard)

    // Verify Dashboard elements
    await expect(page.getByText('Journal Entries')).toBeVisible();
    await expect(page.getByRole('button', { name: /New Note/i })).toBeVisible();
  });
});
