from playwright.sync_api import sync_playwright

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    page = browser.new_page()

    # Mock API responses
    page.route("**/api/users/me", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"id": "user-1", "email": "test@example.com", "settings": {"bible_version": "ESV"}}'
    ))
    page.route("**/api/tags", lambda route: route.fulfill(
        status=200, content_type="application/json", body='[]'
    ))
    page.route("**/api/notes/new", lambda route: route.fulfill(status=200)) # Just in case

    # Set mock user
    page.add_init_script("""
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-user-id',
            email: 'test@example.com',
            displayName: 'Test User'
        }));
    """)

    page.goto("http://localhost:5173/notes/new")

    # Wait for Title
    try:
        page.wait_for_selector('input[placeholder="Title"]', timeout=5000)
        # Wait for Tag Input - using text selector as reliable fallback
        page.wait_for_selector('button:has-text("Add tag...")', timeout=5000)
    except:
        pass

    page.screenshot(path="/home/jules/verification/editor_with_tags_fixed.png")
    browser.close()

if __name__ == "__main__":
    with sync_playwright() as p:
        run(p)
