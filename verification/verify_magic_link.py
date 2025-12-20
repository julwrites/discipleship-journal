
from playwright.sync_api import sync_playwright

def verify_login_page():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Navigate to login page
        page.goto("http://localhost:5173/login")

        # Click on Magic Link tab
        page.get_by_role("tab", name="Magic Link").click()

        # Check if the input and button are visible
        page.get_by_placeholder("name@example.com").fill("test@example.com")

        # Take a screenshot
        page.screenshot(path="verification/magic_link_tab.png")

        browser.close()

if __name__ == "__main__":
    verify_login_page()
