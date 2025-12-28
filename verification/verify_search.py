from playwright.sync_api import sync_playwright

def verify_advanced_search():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        # Use mock auth mode
        context = browser.new_context()
        page = context.new_page()

        page.on("console", lambda msg: print(f"BROWSER CONSOLE: {msg.text}"))
        page.on("pageerror", lambda err: print(f"BROWSER ERROR: {err}"))

        # Inject E2E_TEST_USER for mock auth
        page.add_init_script("""
            localStorage.setItem('E2E_TEST_USER', JSON.stringify({
                uid: 'test-user-123',
                email: 'test@example.com',
                displayName: 'Test User'
            }));
        """)

        # Mock API responses
        # Mock syncUser
        page.route("**/api/users/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"id": "user-uuid-123", "firebase_uid": "test-user-123", "email": "test@example.com"}'
        ))

        # Mock fetchNotes
        page.route("**/api/notes*", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"data": [{"id": "note-1", "title": "Test Note 1", "updated_at": "2023-01-01T00:00:00Z"}, {"id": "note-2", "title": "Test Note 2", "updated_at": "2023-01-02T00:00:00Z"}], "meta": {"total": 2, "page": 1, "limit": 20, "total_pages": 1}}'
        ))

        # Go to Dashboard with mock API key to trigger mock mode
        print("Navigating to dashboard...")
        try:
             page.goto("http://localhost:5173/?FIREBASE_API_KEY=mock-key")
        except Exception as e:
             print(f"Failed to load page: {e}")
             return

        # Wait for notes to load
        print("Waiting for notes...")
        try:
            page.wait_for_selector("text=Test Note 1", timeout=5000)
        except Exception as e:
            print(f"Timeout waiting for notes: {e}")
            page.screenshot(path="verification/failed_load.png")
            print("Saved failure screenshot to verification/failed_load.png")
            return

        print("Notes loaded. Opening filter...")
        # Find Filter button
        filter_btn = page.get_by_title("Filter & Sort")
        if not filter_btn.is_visible():
            print("Filter button not found by title")
            page.screenshot(path="verification/failed_btn.png")
            return

        filter_btn.click()

        # Wait for popover content
        print("Waiting for popover...")
        page.wait_for_selector("text=Sort By")

        # Take screenshot of the filter popover
        page.screenshot(path="verification/advanced_search.png")
        print("Screenshot saved to verification/advanced_search.png")

        browser.close()

if __name__ == "__main__":
    verify_advanced_search()
