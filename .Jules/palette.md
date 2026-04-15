## 2024-10-25 - Icon-only Back Buttons
**Learning:** Icon-only navigational buttons (like `<ArrowLeft />`) are frequently used in top-level secondary pages (Settings, Tags, Memory Verses, etc.) but often lack `aria-label`s, making them invisible to screen readers.
**Action:** Always add `aria-label="Go back"` to back buttons that only contain an icon to ensure users navigating via screen reader know their purpose. Similarly, ensure "Edit" and "Delete" icon buttons have aria labels.
