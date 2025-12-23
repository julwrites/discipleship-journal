---
id: TESTING-20251223-033043-KIR
status: pending
title: Fix duplicate key warnings in Dashboard tests
priority: low
created: 2025-12-23 03:30:43
category: testing
dependencies:
type: chore
---

# Fix duplicate key warnings in Dashboard tests

## Description
During the test run, `Dashboard.test.tsx` produces warnings:
`Encountered two children with the same key, note-10.`

This happens because the test mock for `fetchNotes` returns the same set of notes (ids `note-0` to `note-19`) for both page 1 and page 2. When the component appends page 2 to the list, it creates duplicates.

## Acceptance Criteria
- [ ] Modify `Dashboard.test.tsx` to return distinct IDs for the second page of results.
- [ ] Verify that the warning no longer appears in test output.
