---
id: FEATURES-20251227-034016-BTH
status: completed
title: Frontend Image Support
priority: medium
created: 2025-12-27 03:40:16
category: features
dependencies: []
type: task
---

# Frontend Image Support

## Description
The MVP specifications mention rich media support (Images) in the editor, but the current `RichTextEditor` only supports basic markdown and links. This task aims to add support for images using `@tiptap/extension-image` and improve the UX by replacing `window.prompt` with proper UI components (Popovers).

## Plan
1.  **Dependencies**:
    -   Install `@tiptap/extension-image`.
    -   Install `@radix-ui/react-popover` and `@radix-ui/react-label`.
2.  **UI Components**:
    -   Create `web/src/components/ui/popover.tsx`.
    -   Create `web/src/components/ui/label.tsx`.
3.  **RichTextEditor Update**:
    -   Configure `Image` extension in Tiptap.
    -   Add Image button to toolbar.
    -   Implement Popover for Image URL input.
    -   Refactor Link input to use Popover instead of `window.prompt`.
4.  **Verification**:
    -   Add automated verification script using Playwright.

## Todo
- [x] Install dependencies (`@tiptap/extension-image`, `@radix-ui/react-popover`, `@radix-ui/react-label`) <!-- id: 1 -->
- [x] Add `Popover` component <!-- id: 2 -->
- [x] Add `Label` component <!-- id: 3 -->
- [x] Update `RichTextEditor.tsx` to support Images <!-- id: 4 -->
- [x] Update `RichTextEditor.tsx` to use Popover for Links <!-- id: 5 -->
- [x] Verify changes <!-- id: 6 -->
