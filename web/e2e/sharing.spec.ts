import { test, expect } from '@playwright/test';

test.describe('Note Sharing (Mocked)', () => {
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

    // 3. Mock Existing Note
    await page.route('**/api/notes/note-1', async route => {
        await route.fulfill({
            json: {
                id: 'note-1',
                title: 'My Shared Note',
                content: { markdown: '# Sharing is Caring' },
                created_at: new Date().toISOString(),
                updated_at: new Date().toISOString()
            }
        });
    });

    // 4. Mock Groups List
    await page.route('**/api/groups', async route => {
        await route.fulfill({
            json: [
                { id: 'g1', name: 'Study Group', description: 'Weekly study', role: 'member' },
                { id: 'g2', name: 'Prayer Group', description: 'Daily prayer', role: 'admin' }
            ]
        });
    });
  });

  test('should share a note with a group', async ({ page }) => {
    // Mock Share API
    let shareRequest;
    await page.route('**/api/groups/g1/shares', async route => {
        if (route.request().method() === 'POST') {
            shareRequest = route.request();
            await route.fulfill({ status: 200, json: { success: true } });
        } else {
            await route.fallback();
        }
    });

    // Handle Alert
    const dialogDismissedPromise = new Promise<void>(resolve => {
        page.once('dialog', async dialog => {
            expect(dialog.message()).toBe('Note shared!');
            await dialog.accept();
            resolve();
        });
    });

    await page.goto('/notes/note-1');

    // Wait for note to load
    await expect(page.locator('.ProseMirror')).toBeVisible();
    await expect(page.locator('.ProseMirror')).toContainText('Sharing is Caring');

    // Click Share
    await page.getByRole('button', { name: 'Share' }).click();

    // Verify Dialog Open
    await expect(page.getByRole('heading', { name: 'Share to Group' })).toBeVisible();

    // Select Group
    await page.locator('select').selectOption('g1');

    // Add Comment
    await page.getByPlaceholder('Add a comment (optional)...').fill('Check this out!');

    // Click Share Note
    await page.getByRole('button', { name: 'Share Note' }).click();

    // Verify Request
    expect(shareRequest).toBeTruthy();
    const postData = shareRequest.postDataJSON();
    expect(postData.note_id).toBe('note-1');
    expect(postData.comment).toBe('Check this out!');

    await dialogDismissedPromise;
  });
});
