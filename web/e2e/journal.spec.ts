import { test, expect } from '@playwright/test';

test.describe('Journaling (Mocked)', () => {
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
                username: 'testuser',
                settings: {}
            }
        });
    });
  });

  test('should create a new note', async ({ page }) => {
    // Mock Create Note Promise
    let resolveRequest;
    const createNoteRequestPromise = new Promise(resolve => { resolveRequest = resolve; });

    // Stateful mock
    let notes = [];

    // Single route handler for notes to avoid conflicts
    await page.route(/.*\/api\/notes/, async route => {
        const method = route.request().method();

        if (method === 'GET') {
             await route.fulfill({ json: notes });
        } else if (method === 'POST') {
            resolveRequest(route.request());
            const body = route.request().postDataJSON();
            const newNote = {
                id: 'note-new-1',
                title: body.title,
                content: body.content,
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
            };
            notes = [newNote]; // Update state
            await route.fulfill({
                json: newNote
            });
        } else {
             await route.fallback();
        }
    });

    await page.goto('/');

    // Click New Note
    await page.getByTitle('New Note').click();

    // Verify we are on the editor
    await expect(page.locator('.ProseMirror')).toBeVisible();

    // Type Title
    await page.getByPlaceholder('Title').fill('My Daily Journal');

    // Type content
    await page.locator('.ProseMirror').click();
    await page.keyboard.type('# My Daily Journal\n\nToday I learned about grace.');

    // Save
    await page.getByRole('button', { name: 'Save' }).click();

    // Explicitly wait for the Create Request to be made
    const request = await createNoteRequestPromise;
    expect(request).toBeTruthy();
    expect(request.postDataJSON().content.markdown).toContain('# My Daily Journal');

    // After save, we expect to be redirected or see a success message.
    // Or we go back to dashboard and see the note.

    // Let's assume we go back to dashboard manually if needed, or if Save redirects.
    // For now, let's just verify the API call happened (Playwright waits for it if we await route).

    // But since we mocked the POST, we want to see the effect on Dashboard.
    // Let's reload dashboard and mock the list again with one item.

    // Removed duplicate route handler as the stateful one handles it.

    // Go back to dashboard
    await page.goto('/');

    // Verify the note card is visible
    // The dashboard usually shows a snippet.
    // If the content is complex JSON, we might not see "My Daily Journal" easily unless the component renders it.
    // But let's check for "My Daily Journal" text.
    await expect(page.getByText('My Daily Journal')).toBeVisible();
  });
});
