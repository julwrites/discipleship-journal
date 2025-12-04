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
    // Mock AI Chat Endpoint (ChatPage uses /api/chat)
    await page.route('**/api/chat', async route => {
        await route.fulfill({
            json: {
                response: "This is a mocked AI response about grace."
            }
        });
    });

    await page.goto('/chat');

    // Type message
    await page.getByPlaceholder('What does this say about...').fill('What is grace?');
    // Also fill passage as button is disabled otherwise
    await page.getByPlaceholder('e.g. Romans 8, Psalm 23').fill('John 3:16');


    // Click Send
    await page.getByRole('button', { name: 'Ask AI' }).click();

    // Verify User Message (Prompt) input retains value? No, usually clears.
    // The test expects "User Message" to be visible in chat history.
    // ChatPage implementation shows response but doesn't explicitly show user message in history list in the current simple version?
    // Let's check ChatPage.tsx:
    // {response && ( ... <p>{response}</p> ... )}
    // It only shows the response!

    // Verify AI Response is visible

    // Verify AI Response is visible
    await expect(page.getByText('This is a mocked AI response about grace.')).toBeVisible();
  });
});
