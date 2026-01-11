import time
from playwright.sync_api import sync_playwright

def run():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        # 1. Load context with mock user/token
        context = browser.new_context()
        page = context.new_page()

        # 2. Setup mock routes
        # Mock GET /api/users/me
        page.route("**/api/users/me", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"id": "user1", "username": "Test User", "settings": {"bible_version": "ESV"}}'
        ))

        # Mock GET /api/bible/versions
        page.route("**/api/bible/versions*", lambda route: route.fulfill(
            status=200,
            content_type="application/json",
            body='{"data": [{"id": "ESV", "name": "English Standard Version", "abbreviation": "ESV", "language": "English"}, {"id": "NIV", "name": "New International Version", "abbreviation": "NIV", "language": "English"}, {"id": "KJV", "name": "King James Version", "abbreviation": "KJV", "language": "English"}], "meta": {"total": 3, "page": 1, "limit": 20}}'
        ))

        # Mock GET /api/bible/passage?ref=John+3%3A16&version=ESV
        page.route("**/api/bible/passage*", lambda route: route.fulfill(
             status=200,
             content_type="application/json",
             body='{"verse": "For God so loved the world...", "reference": "John 3:16", "version": "ESV"}'
        ))

        # 3. Inject User into LocalStorage to bypass Auth
        page.goto("http://localhost:5173")
        page.evaluate("localStorage.setItem('E2E_TEST_USER', 'true')")
        page.reload()

        # 4. Navigate to Settings to see the Selector
        page.goto("http://localhost:5173/settings")
        time.sleep(2) # Wait for load

        # 5. Screenshot Settings
        page.screenshot(path="verification/settings_page.png")
        print("Screenshot saved to verification/settings_page.png")

        # 6. Click on selector to open popover
        page.get_by_role("combobox").click()
        time.sleep(1)
        page.screenshot(path="verification/settings_popover.png")
        print("Screenshot saved to verification/settings_popover.png")

        browser.close()

if __name__ == "__main__":
    run()
