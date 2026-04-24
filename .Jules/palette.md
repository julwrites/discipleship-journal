## 2024-04-24 - Context-Aware ARIA Labels on Card Actions
**Learning:** Icon-only buttons within iterated components (like `NoteCard`) require context-aware `aria-label`s (e.g., "Delete [Item Title]"). Generic labels like "Delete" or `title` attributes alone create ambiguity for screen reader users when multiple cards are present on the same page.
**Action:** Always interpolate identifying item information into `aria-label`s for actions inside lists and cards.
