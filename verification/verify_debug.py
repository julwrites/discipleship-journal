
from playwright.sync_api import sync_playwright

def verify_note_editor():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()

        # Inject mock user
        context.add_init_script("""
            localStorage.setItem('E2E_TEST_USER', JSON.stringify({
                uid: 'test-user-id',
                email: 'test@example.com'
            }));
            localStorage.setItem('firebase:authUser:test-api-key:[DEFAULT]', JSON.stringify({
                uid: 'test-user-id',
                email: 'test@example.com',
                stsTokenManager: { accessToken: 'mock-token' }
            }));
        """)

        page = context.new_page()

        # Mock API calls
        # Mock Create Note
        page.route('**/api/notes', lambda route: route.fulfill(
            status=201,
            body='{"id": "new-note-id", "title": "Test Note", "content": ""}',
            headers={'Content-Type': 'application/json'}
        ))

        # Mock Get Note
        page.route('**/api/notes/new-note-id', lambda route: route.fulfill(
            status=200,
            body='{"id": "new-note-id", "title": "Test Note", "content": ""}',
            headers={'Content-Type': 'application/json'}
        ))

        print('Navigating to http://localhost:5173/notes/new')
        page.goto('http://localhost:5173/notes/new')

        # Take a screenshot to debug
        page.screenshot(path='verification/debug.png')
        print('Debug screenshot saved')

        browser.close()

if __name__ == '__main__':
    verify_note_editor()
