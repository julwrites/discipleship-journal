## 2024-05-15 - Missing aria-labels on icon-only buttons
**Learning:** Found several icon-only buttons using only the `title` attribute without `aria-label` for screen readers across the dashboard and editor components.
**Action:** Always verify that icon-only buttons have an explicit `aria-label` alongside or instead of just `title` to ensure consistent accessibility for screen reader users.
