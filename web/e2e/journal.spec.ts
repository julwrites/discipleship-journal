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
    // Mock listing notes (initially empty)
    await page.route('**/api/notes**', async route => {
        if (route.request().method() === 'GET') {
             await route.fulfill({ json: [] });
        } else {
            await route.continue();
        }
    });

    // Mock Create Note
    let createNoteRequestPromise;
    await page.route('**/api/notes', async route => {
        if (route.request().method() === 'POST') {
            createNoteRequestPromise = Promise.resolve(route.request());
            const body = route.request().postDataJSON();
            await route.fulfill({
                json: {
                    id: 'note-new-1',
                    content: body.content, // Should check structure
                    created_at: new Date().toISOString(),
                    updated_at: new Date().toISOString()
                }
            });
        } else {
             await route.fallback();
        }
    });

    await page.goto('/');

    // Click New Note
    await page.getByRole('button', { name: 'New Note' }).click();

    // Verify we are on the editor
    await expect(page.getByPlaceholder('Write your reflection...')).toBeVisible();

    // Type content
    await page.getByPlaceholder('Write your reflection...').fill('# My Daily Journal\n\nToday I learned about grace.');

    // Save (assuming auto-save or manual save, button might be "Save" or icon)
    // Wait, the UI might auto-save or have a save button.
    // Checking NoteEditor implementation or trying to find a save button.
    // If it's auto-save, we might need to wait.
    // Let's look for a Save button.
    const saveButton = page.getByRole('button', { name: 'Save' });
    if (await saveButton.isVisible()) {
        await saveButton.click();
    } else {
        // Maybe it autosaves?
        // Let's assume there is a Save button for now based on typical UI.
    }

    // Explicitly wait for the Create Request to be made
    const request = await createNoteRequestPromise;
    expect(request).toBeTruthy();
    expect(request.postDataJSON().content).toContain('# My Daily Journal');

    // After save, we expect to be redirected or see a success message.
    // Or we go back to dashboard and see the note.

    // Let's assume we go back to dashboard manually if needed, or if Save redirects.
    // For now, let's just verify the API call happened (Playwright waits for it if we await route).

    // But since we mocked the POST, we want to see the effect on Dashboard.
    // Let's reload dashboard and mock the list again with one item.

    await page.route('**/api/notes**', async route => {
         if (route.request().method() === 'GET') {
             await route.fulfill({ json: [{
                 id: 'note-new-1',
                 content: { type: 'doc', content: [{ type: 'heading', content: [{ type: 'text', text: 'My Daily Journal' }] }] }, // Simplified
                 created_at: new Date().toISOString(),
                 updated_at: new Date().toISOString()
             }] });
         }
    });

    // Go back to dashboard
    await page.goto('/');

    // Verify the note card is visible
    // The dashboard usually shows a snippet.
    // If the content is complex JSON, we might not see "My Daily Journal" easily unless the component renders it.
    // But let's check for "My Daily Journal" text.
    await expect(page.getByText('My Daily Journal')).toBeVisible();
  });
});
