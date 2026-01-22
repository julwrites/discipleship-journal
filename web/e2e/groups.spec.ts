import { test, expect } from '@playwright/test';

// Common mock data
const mockUser = {
  uid: 'test-user-123',
  email: 'test@example.com',
  displayName: 'Test User',
  photoURL: 'https://example.com/photo.jpg'
};

const mockMember = {
    user_id: 'member-1',
    display_name: 'Member One',
    email: 'member@example.com',
    role: 'member'
};

test.describe('Groups (Mocked)', () => {
  test.beforeEach(async ({ page }) => {
    // Mock API requests
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

    await page.route('**/api/groups', async (route) => {
      if (route.request().method() === 'GET') {
        await route.fulfill({
            contentType: 'application/json',
            body: JSON.stringify([{ id: 'g1', name: 'Bible Study', description: 'Weekly study', role: 'admin' }])
        });
      } else if (route.request().method() === 'POST') {
        await route.fulfill({ status: 201, body: JSON.stringify({ id: 'g2' }) });
      }
    });

    // Mock Groups Search
    await page.route('**/api/groups/search*', async (route) => {
      await route.fulfill({
          contentType: 'application/json',
          body: JSON.stringify([{ id: 'g3', name: 'Prayer Warriors', description: 'Prayer group', role: '' }])
      });
    });

    // Mock Join
    await page.route('**/api/groups/g3/join', async (route) => {
      await route.fulfill({ status: 200 });
    });

    // Mock Shares List
    await page.route('**/api/groups/*/shares', async (route) => {
        await route.fulfill({
            contentType: 'application/json',
            body: JSON.stringify([])
        });
    });

    // Mock Members List
    // We mock ALL GET requests to members endpoint for ANY group or specific group g1
    await page.route('**/api/groups/*/members*', async (route) => {
        if (route.request().method() === 'GET') {
            await route.fulfill({
                contentType: 'application/json',
                body: JSON.stringify([
                    { user_id: 'uuid-123', display_name: 'Test User', email: 'test@example.com', role: 'admin' },
                    mockMember
                ])
            });
        } else if (route.request().method() === 'POST') {
             await route.fulfill({ status: 201, json: { success: true } });
        } else if (route.request().method() === 'DELETE') {
             await route.fulfill({ status: 200, json: { success: true } });
        }
    });

    // Mock User Search for Adding
    await page.route('**/api/users/search*', async (route) => {
        await route.fulfill({
            contentType: 'application/json',
            body: JSON.stringify([{ id: 'new-user', display_name: 'New User', email: 'new@example.com', username: 'newuser' }])
        });
    });

    // Mock Connections for Add Member dialog
    await page.route('**/api/connections', async (route) => {
        await route.fulfill({
            contentType: 'application/json',
            body: JSON.stringify([{
                id: 'conn-1',
                requester_id: 'uuid-123',
                receiver_id: 'new-user',
                status: 'accepted',
                requester_email: 'test@example.com',
                receiver_email: 'new@example.com'
            }])
        });
    });

    // Inject user
    await page.addInitScript((user) => {
        localStorage.setItem('E2E_TEST_USER', JSON.stringify(user));
    }, mockUser);

    await page.goto('/groups');
  });

  test('should display my groups and expand members', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Groups' })).toBeVisible();
    await expect(page.locator('text=Bible Study')).toBeVisible();

    // Click to expand
    await page.locator('text=Bible Study').click();

    // Click Members Tab
    await page.getByRole('tab', { name: 'Members' }).click();

    // Wait for the Members section to appear.
    await expect(page.getByRole('heading', { name: 'Members' })).toBeVisible();

    // Verify members are listed
    await expect(page.locator('text=Test User')).toBeVisible();
    await expect(page.locator('text=Member One')).toBeVisible();
  });

  test('should add a member as admin', async ({ page }) => {
    await page.locator('text=Bible Study').click();

    // Click Members Tab
    await page.getByRole('tab', { name: 'Members' }).click();
    await expect(page.getByRole('heading', { name: 'Members' })).toBeVisible();

    await page.getByRole('button', { name: 'Add Member' }).click();

    // Wait for dialog and connection to appear
    await expect(page.locator('text=Add Member to Bible Study')).toBeVisible();
    await expect(page.locator('text=new@example.com')).toBeVisible();
    await page.getByRole('button', { name: 'Add' }).first().click();

    // Verify dialog closed
    await expect(page.locator('text=Add Member to Bible Study')).not.toBeVisible();
  });

  test('should remove a member as admin', async ({ page }) => {
    // Setup dialog listener before action
    page.once('dialog', async dialog => {
        await dialog.accept();
    });

    await page.locator('text=Bible Study').click();

    // Click Members Tab
    await page.getByRole('tab', { name: 'Members' }).click();
    await expect(page.getByRole('heading', { name: 'Members' })).toBeVisible();
    await expect(page.locator('text=Member One')).toBeVisible();

    const memberRow = page.locator('div.flex.justify-between.items-center', { hasText: 'Member One' });
    await memberRow.getByRole('button').click();
  });
});
