import { test, expect } from '@playwright/test';

test.describe('Groups Feature', () => {
    test.beforeEach(async ({ page }) => {
        // Set E2E_TEST_USER in localStorage to simulate logged-in user
        // We need to do this before page load, but we can't access localStorage before page.goto.
        // So we goto '/', inject script, then reload or just let it update.
        // Actually, playwright has `addInitScript` which runs before page loads.
        await page.addInitScript(() => {
            localStorage.setItem('E2E_TEST_USER', JSON.stringify({
                uid: 'test-user-id',
                displayName: 'Test User',
                email: 'test@example.com'
            }));
        });
    });

    test('Create and Join Group Flow', async ({ page }) => {
        // 1. Mock API Responses
        // Debug requests
        page.on('request', request => console.log('>>', request.method(), request.url()));

        // List Groups (initially empty)
        await page.route('**/api/groups', async (route) => {
            if (route.request().method() === 'GET') {
                console.log('Intercepted GET /api/groups');
                await route.fulfill({ json: [] });
            } else if (route.request().method() === 'POST') {
                console.log('Intercepted POST /api/groups');
                // Create Group
                await route.fulfill({ status: 201, json: { id: 'new-group-id' } });
            } else {
                await route.continue();
            }
        });

        // Search Groups
        await page.route('**/api/groups/search*', async (route) => {
             await route.fulfill({ json: [
                 { id: 'found-group-id', name: 'Found Group', description: 'Description', role: '' }
             ] });
        });

        // Join Group
        await page.route('**/api/groups/*/join', async (route) => {
            await route.fulfill({ status: 200 });
        });

        // Get Members
        await page.route('**/api/groups/*/members', async (route) => {
            await route.fulfill({ json: [
                { user_id: 'u1', display_name: 'Test User', email: 'test@example.com', role: 'admin' }
            ] });
        });

        // 2. Go to Groups Page
        await page.goto('/groups');

        // Check if we are on Groups page
        await expect(page.locator('h1')).toHaveText('Groups');

        // 3. Create Group
        await page.click('text=Create Group');
        await page.fill('input[placeholder="Group Name"]', 'My New Group');
        await page.fill('textarea[placeholder="Description"]', 'This is a test group');

        // Intercept the POST and the subsequent GET
        await page.route('**/api/groups', async (route) => {
            if (route.request().method() === 'POST') {
                 console.log('Intercepted POST /api/groups');
                 await route.fulfill({ status: 201, json: { id: 'new-group-id' } });
            } else if (route.request().method() === 'GET') {
                console.log('Intercepted GET /api/groups - returning new list');
                await route.fulfill({ json: [
                    { id: 'new-group-id', name: 'My New Group', description: 'This is a test group', role: 'admin' }
                ] });
            } else {
                await route.continue();
            }
        });

        // Submit
        // Use a more robust selector for the submit button in the dialog
        await page.click('button:has-text("Create") >> visible=true >> nth=-1');

        // Verify group appears
        await expect(page.locator('text=My New Group')).toBeVisible({ timeout: 10000 });

        // 4. Search and Join
        await page.click('text=Find Groups');
        await page.fill('input[placeholder="Search groups..."]', 'Found');
        await page.click('button:has-text("Search")');

        await expect(page.locator('text=Found Group')).toBeVisible();

        // Mock join success alert
        // We need to set up the handler before triggering the dialog.
        // And since `dialog.accept()` is async but `page.on` callback is not awaited by the event emitter in the same way,
        // we should ensure it handles it.
        // Playwright's default behavior for dialogs is to dismiss them unless a handler is set.
        // If we set a handler, we must handle it.

        // Wait for dialog event
        const dialogPromise = page.waitForEvent('dialog');

        await page.click('text=Join Group');

        const dialog = await dialogPromise;
        await dialog.accept();

        // Verify we switch back or refresh? The logic calls `handleSearch` and `fetchMyGroups`.
        // We just check if alert was triggered (implicitly handled) and maybe button state changes if we mocked the search response to update.
    });
});
