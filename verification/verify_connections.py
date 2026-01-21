from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()

    # Set localStorage before any page load
    context.add_init_script("""
        localStorage.setItem('E2E_TEST_USER', JSON.stringify({uid: 'test-uid', email: 'me@example.com'}));
    """)

    page = context.new_page()

    # Mock Connections API (List)
    page.route("**/api/connections", lambda route: route.fulfill(json=[]))

    # Mock Search API
    def handle_search(route):
        print("Search API called!")
        route.fulfill(json=[
            {"id": "1", "email": "user1@example.com", "username": "CoolUser1"},
            {"id": "2", "email": "user2@example.com"} # No username
        ])

    page.route("**/api/users/search*", handle_search)

    # Go to Connections page directly
    page.goto("http://localhost:5173/connections")

    # Wait for AuthGuard to let us in (check for header)
    expect(page.get_by_role("heading", name="Connections")).to_be_visible()

    # Click "Find People" tab
    page.get_by_role("tab", name="Find People").click()

    # Verify Placeholder
    search_input = page.get_by_placeholder("Search by email or username...")
    expect(search_input).to_be_visible()

    # Type search query
    search_input.fill("test")

    # Wait for results
    print("Waiting for results...")
    # User 1
    user1_card = page.locator(".space-y-2 > div").filter(has_text="CoolUser1").first
    expect(user1_card).to_be_visible()
    expect(user1_card.locator("p.font-medium")).to_have_text("CoolUser1")
    expect(user1_card.locator(".text-muted-foreground").first).to_have_text("user1@example.com")

    # User 2
    user2_card = page.locator(".space-y-2 > div").filter(has_text="user2@example.com").first
    expect(user2_card).to_be_visible()
    expect(user2_card.locator("p.font-medium")).to_have_text("user2@example.com")

    # Take screenshot
    page.screenshot(path="verification/connections_verification.png")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
