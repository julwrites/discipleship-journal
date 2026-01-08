from playwright.sync_api import sync_playwright, expect
import time

def verify_unsaved_changes(page):
    # Enable mock mode
    page.context.add_init_script("""
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-user',
            email: 'test@example.com',
            displayName: 'Test User'
        }));
    """)

    # Mock API
    page.route("**/api/notes?*", lambda route: route.fulfill(
        status=200, content_type="application/json", body='{"data": [], "meta": {"total_pages": 1}}'
    ))
    page.route("**/api/notes", lambda route: route.fulfill(
        status=200, content_type="application/json", body='{"data": [], "meta": {"total_pages": 1}}'
    ))
    page.route("**/api/notes/new", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"id": "new-note-id", "title": "New Note", "content": "Content"}'
    ))
    page.route("**/api/notes/*", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"id": "note-id", "title": "Dirty Title", "content": "Content"}'
    ))

    # Listen for console logs
    page.on("console", lambda msg: print(f"Console: {msg.text}"))
    page.on("pageerror", lambda err: print(f"PageError: {err}"))

    print("Navigating to Dashboard...")
    page.goto("http://localhost:5173/")

    # Check if dashboard loaded
    expect(page.get_by_text("My Journal")).to_be_visible()

    print("Clicking New Note...")
    # Use robust locator
    page.locator('button[title="New Note"]').click()

    # Wait for editor to load
    page.wait_for_selector('input[placeholder="Title"]')
    print("Editor loaded.")

    print("Modifying note...")
    page.fill('input[placeholder="Title"]', 'Dirty Title')

    # Click Back
    print("Clicking Back...")
    page.click("text=← Back")

    # Expect Alert Dialog
    print("Waiting for dialog...")
    dialog = page.locator("div[role='alertdialog']")
    expect(dialog).to_be_visible()
    expect(dialog).to_contain_text("Unsaved Changes")

    print("Taking screenshot...")
    page.screenshot(path="verification/unsaved_changes_dialog.png")

    # Click Cancel
    print("Clicking Cancel...")
    page.get_by_role("button", name="Cancel").click()
    expect(dialog).not_to_be_visible()

    # Click Back again and Leave
    print("Leaving...")
    page.click("text=← Back")
    page.get_by_role("button", name="Leave").click()

    # Should be on dashboard
    expect(page.get_by_text("My Journal")).to_be_visible()
    print("Verified: returned to dashboard.")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch()
        page = browser.new_page()
        try:
            verify_unsaved_changes(page)
        except Exception as e:
            print(f"Error: {e}")
            page.screenshot(path="verification/error_retry.png")
        finally:
            browser.close()
