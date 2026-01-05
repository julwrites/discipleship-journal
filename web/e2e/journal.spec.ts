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
             await route.fulfill({
                 json: {
                     data: notes,
                     meta: { total: notes.length, page: 1, limit: 20, total_pages: 1 }
                 }
             });
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
                json: { id: newNote.id }
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
    // Expect content to be a string or contain the text directly (depending on how we type it)
    // Now we save as string (HTML) or string (Markdown if we typed markdown but removed md support?)
    // Wait, Tiptap HTML of "# My Daily Journal" depends on extension.
    // If we have Markdown extension, we can type markdown but output might be HTML.
    // However, the test types: `page.keyboard.type('# My Daily Journal\n\nToday I learned about grace.');`
    // Tiptap with Markdown extension might convert this to `<h1>My Daily Journal</h1>...` on the fly.
    // So the stored content might be HTML: `<h1>My Daily Journal</h1>...`.
    // So we should check for text containment rather than strict structure.
    const content = request.postDataJSON().content;
    expect(JSON.stringify(content)).toContain('My Daily Journal');

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
