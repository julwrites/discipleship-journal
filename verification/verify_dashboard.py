
import os
import json
from playwright.sync_api import sync_playwright, expect

def verify_dashboard_dialog(page):
    # Mock API responses
    # 1. Login/Me
    page.route("**/api/users/me", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({"id": "user123", "email": "test@example.com", "username": "TestUser"})
    ))

    # 2. Notes list (empty or not, doesn't matter for this test as we test the AI dialog)
    page.route("**/api/notes?*", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({"data": [], "meta": {"total_pages": 1}})
    ))

    # 3. GetNote (called when AI dialog opens)
    page.route("**/api/notes/note123", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({
            "id": "note123",
            "title": "Test Note",
            "content": "<p>Note content</p>",
            "updated_at": "2024-01-01T12:00:00Z"
        })
    ))

    # 4. Ask AI (return HTML content)
    page.route("**/api/ai/ask", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({
            "response": "<h3>AI Answer</h3><p>This is a <b>formatted</b> response.</p><ul><li>Point 1</li><li>Point 2</li></ul>"
        })
    ))

    # 5. Connection requests
    page.route("**/api/connections", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps([])
    ))

    # 6. Reading Plans
    page.route("**/api/my-reading-plans", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps([])
    ))

    # Inject mock user into localStorage
    page.add_init_script("""
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({
            uid: 'test-uid',
            email: 'test@example.com'
        }));
    """)

    # Go to Dashboard
    page.goto("http://localhost:5173/")

    # Wait for page to load
    expect(page.get_by_text("My Journal")).to_be_visible()

    # Trigger AI Dialog manually (since we don't have notes, we can't click the button easily unless we mock notes)
    # Let's mock a note in the list first
    page.route("**/api/notes?*", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body=json.dumps({
            "data": [{
                "id": "note123",
                "title": "Test Note",
                "updated_at": "2024-01-01T12:00:00Z"
            }],
            "meta": {"total_pages": 1}
        })
    ))

    # Reload to get the note
    page.reload()
    expect(page.get_by_text("Test Note")).to_be_visible()

    # Click Ask AI button (sparkles icon)
    # It's in the note card actions. Use first() because there's also a FAB
    # Or better, target the one in the card
    page.locator(".group").first.get_by_title("Ask AI").click()

    # Wait for dialog
    expect(page.get_by_role("dialog")).to_be_visible()
    expect(page.get_by_text("Ask a question about")).to_be_visible()

    # Fill prompt
    page.get_by_placeholder("Ask a question...").fill("Explain this")

    # Click Ask
    page.get_by_role("button", name="Ask").click()

    # Wait for response (which we mocked as HTML)
    # The response contains "AI Answer" in h3
    expect(page.get_by_role("heading", name="AI Answer")).to_be_visible()
    expect(page.get_by_text("This is a formatted response.")).to_be_visible()

    # Take screenshot
    page.screenshot(path="verification/dashboard_ai_dialog.png")
    print("Screenshot saved to verification/dashboard_ai_dialog.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch()
        context = browser.new_context()
        page = context.new_page()
        try:
            verify_dashboard_dialog(page)
        except Exception as e:
            print(f"Error: {e}")
        finally:
            browser.close()
