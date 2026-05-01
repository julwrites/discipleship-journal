
## 2024-05-18 - Added aria-label to Dashboard icon buttons
**Learning:** React components returning icon-only buttons via `lucide-react` icons (like `Users`, `BookOpen`) inside Radix UI primitives (`Button`, `DropdownMenuTrigger`) must explicitly have `aria-label` attributes set to be accessible by screen readers. The `title` attribute is primarily for visual tooltips and does not sufficiently substitute for `aria-label` for assistive technologies in all contexts.
**Action:** When adding new icon-only buttons to the UI, particularly within navigation bars or dropdown triggers, always include a descriptive `aria-label` attribute alongside the `title` attribute.
