import { test, expect } from '@playwright/test';

test.describe('Chat & AI (Mocked)', () => {
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
  });

  test('should send a message and receive an AI response', async ({ page }) => {
    // Mock AI Ask Endpoint
    await page.route('**/api/ai/ask', async route => {
        await route.fulfill({
            json: {
                response: "This is a mocked AI response about grace."
            }
        });
    });

    await page.goto('/chat');

    // Type message
    await page.getByPlaceholder('Ask a question...').fill('What is grace?');

    // Click Send
    await page.getByRole('button', { name: 'Send' }).click();

    // Verify User Message is visible
    await expect(page.getByText('What is grace?')).toBeVisible();

    // Verify AI Response is visible
    await expect(page.getByText('This is a mocked AI response about grace.')).toBeVisible();
  });
});
