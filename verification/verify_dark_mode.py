from playwright.sync_api import sync_playwright

def verify_dark_mode():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        context = browser.new_context()
        page = context.new_page()

        # Inject E2E_TEST_USER for mock auth and force Dark Mode
        page.add_init_script("""
            localStorage.setItem('E2E_TEST_USER', JSON.stringify({
                uid: 'test-user-123',
                email: 'test@example.com',
                displayName: 'Test User'
            }));
            localStorage.setItem('ui-theme-mode', 'dark');
        """)

        # Mock API responses (minimal)
        page.route("**/api/users/me", lambda route: route.fulfill(
            status=200, content_type="application/json", body='{"id": "user-uuid-123"}'
        ))
        page.route("**/api/notes*", lambda route: route.fulfill(
            status=200, content_type="application/json", body='{"data": [], "meta": {"total": 0}}'
        ))

        print("Navigating to dashboard...")
        try:
             page.goto("http://localhost:5173/?FIREBASE_API_KEY=mock-key")
        except Exception as e:
             print(f"Failed to load page: {e}")
             return

        print("Waiting for dashboard...")
        try:
             page.wait_for_selector("text=My Journal", timeout=5000)
        except Exception as e:
             print("Dashboard didn't load (or title mismatch).")
             page.screenshot(path="verification/failed_dark.png")
             return

        # Open Filter
        print("Opening filter...")
        filter_btn = page.get_by_title("Filter & Sort")
        if filter_btn.is_visible():
            filter_btn.click()
            page.wait_for_selector("text=Sort By")
        else:
            print("Filter button not visible")
            return

        # Take screenshot
        page.screenshot(path="verification/dark_mode_filter.png")
        print("Screenshot saved to verification/dark_mode_filter.png")

        browser.close()

if __name__ == "__main__":
    verify_dark_mode()
