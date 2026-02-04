from playwright.sync_api import sync_playwright
import time
import json
import os

def test_auto_save():
    # Ensure verification directory exists
    os.makedirs("/home/jules/verification", exist_ok=True)

    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()

        # Capture console logs
        page.on("console", lambda msg: print(f"CONSOLE: {msg.text}"))

        print("Navigating to home...")
        # Set E2E_TEST_USER in localStorage before loading the app
        page.goto("http://localhost:3000/")

        user_json = json.dumps({
            "uid": "test-user",
            "email": "test@example.com",
            "displayName": "Test User"
        })

        print("Setting mock user...")
        page.evaluate(f"localStorage.setItem('E2E_TEST_USER', '{user_json}')")

        # Navigate to new note
        print("Navigating to new note...")
        page.goto("http://localhost:3000/notes/new")

        # Wait for editor to load (title input)
        print("Waiting for editor...")
        page.wait_for_selector("input[placeholder='Title']")

        # Type in title
        print("Typing title...")
        page.fill("input[placeholder='Title']", "Auto Save Test")

        # Type in content
        print("Typing content...")
        page.click(".ProseMirror")
        page.keyboard.type("Testing auto save...")

        # Wait for auto-save to trigger (2s debounce + extra)
        print("Waiting for auto-save...")
        time.sleep(3.5)

        # Check for indicators.
        content = page.content()

        if "Error saving" in content:
            print("Detected 'Error saving'.")
        elif "Saved at" in content:
            print("Detected 'Saved at'.")
        elif "Saving..." in content:
            print("Detected 'Saving...'.")
        else:
            print("No auto-save indicator found.")

        page.screenshot(path="/home/jules/verification/verification.png")
        print("Screenshot saved.")

        browser.close()

if __name__ == "__main__":
    test_auto_save()
