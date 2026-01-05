
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
        page.route('**/api/notes', lambda route: route.fulfill(
            status=201,
            body='{"id": "new-note-id", "title": "Test Note", "content": ""}',
            headers={'Content-Type': 'application/json'}
        ))

        page.route('**/api/notes/new-note-id', lambda route: route.fulfill(
            status=200,
            body='{"id": "new-note-id", "title": "Test Note", "content": ""}',
            headers={'Content-Type': 'application/json'}
        ))

        page.route('**/api/bible/passage*', lambda route: route.fulfill(
            status=200,
            body='{"verse": "<p>In the beginning was the Word.</p>", "text": "<p>In the beginning was the Word.</p>"}',
            headers={'Content-Type': 'application/json'}
        ))

        page.route('**/api/ai/ask', lambda route: route.fulfill(
            status=200,
            body='{"response": "<p>This is an <strong>HTML</strong> response.</p>"}',
            headers={'Content-Type': 'application/json'}
        ))

        print('Navigating to http://localhost:5173/notes/new')
        page.goto('http://localhost:5173/notes/new')

        print('Waiting for editor selector')
        page.wait_for_selector('.ProseMirror', timeout=60000)

        # 1. Test Scripture Insertion
        print('Testing Scripture Insertion')
        page.click('button:has-text("Add Scripture")')
        page.fill('textarea[placeholder*="John 3:16"]', 'John 1:1')
        page.click('button:has-text("Search")')
        page.wait_for_timeout(1000)

        page.screenshot(path='verification/scripture_dialog.png')

        page.click('button:has-text("Insert into Note")')
        page.wait_for_timeout(500)
        page.keyboard.press('Escape')
        page.wait_for_timeout(500)

        # 2. Test AI Insertion
        print('Testing AI Insertion')
        page.click('button:has-text("Ask AI")', force=True)
        page.wait_for_selector('textarea[placeholder*="Ask a question"]')
        page.fill('textarea[placeholder*="Ask a question"]', 'What is this?')

        # Use simpler selector
        page.click('div[role="dialog"] >> button:has-text("Ask")')
        page.wait_for_timeout(1000)

        page.screenshot(path='verification/ai_dialog.png')

        page.click('button:has-text("Insert into Note")')
        page.wait_for_timeout(500)

        page.screenshot(path='verification/editor_result.png')

        browser.close()

if __name__ == '__main__':
    verify_note_editor()
