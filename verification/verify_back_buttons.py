import time
from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    # 1. Login (Mock)
    # Inject E2E user into localStorage to bypass login
    page.add_init_script("""
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-user-id',
            email: 'test@example.com',
            displayName: 'Test User'
        }));
    """)

    # Go to Dashboard
    print("Navigating to Dashboard...")
    page.goto("http://localhost:3000/")

    # Wait for dashboard to load
    try:
        expect(page.get_by_text("My Journal")).to_be_visible(timeout=10000)
    except Exception as e:
        print(f"Failed to find 'My Journal'. Current URL: {page.url}")
        page.screenshot(path="verification/debug_dashboard_fail.png")
        raise e

    # 2. Verify "Templates" menu item
    print("Checking Templates menu...")
    page.get_by_role("button", name="Resources").click()
    templates_menu_item = page.get_by_role("menuitem", name="Templates")
    expect(templates_menu_item).to_be_visible()

    # Navigate to Templates
    templates_menu_item.click()

    # 3. Verify Templates Page
    print("Verifying Templates Page...")
    expect(page.get_by_role("heading", name="Templates")).to_be_visible()
    expect(page.get_by_text("Create reusable AI prompts for devotionals, studies, or mentoring.")).to_be_visible()

    # Verify Back Button
    back_button = page.get_by_role("button", name="Back to Dashboard")
    expect(back_button).to_be_visible()

    # Take screenshot of Templates Page
    page.screenshot(path="verification/templates_page.png")
    print("Templates Page verified.")

    # Click Back Button
    back_button.click()
    expect(page.get_by_text("My Journal")).to_be_visible()

    # 4. Verify Memory Verses Page
    print("Verifying Memory Verses Page...")
    page.get_by_role("button", name="Resources").click()
    page.get_by_role("menuitem", name="Memory Verses").click()

    expect(page.get_by_role("heading", name="Scripture Memory")).to_be_visible()

    # Verify Back Button
    back_button = page.get_by_role("button", name="Back to Dashboard")
    expect(back_button).to_be_visible()

    # Take screenshot of Memory Verses Page
    page.screenshot(path="verification/memory_verses_page.png")
    print("Memory Verses Page verified.")

    # Click Back Button
    back_button.click()
    expect(page.get_by_text("My Journal")).to_be_visible()

    print("All verifications passed!")
    browser.close()

with sync_playwright() as playwright:
    run(playwright)
