## 2024-04-13 - Add Missing ARIA Labels to Icon-Only Buttons
**Learning:** React Router `navigate("/")` is commonly bound to "Back" arrow icon buttons across pages without screen-reader descriptions. The `Dashboard.tsx` uses multiple icon-only floating/header buttons.
**Action:** When adding an icon button component, immediately add `aria-label` attribute if text is omitted, ensuring that tools like VoiceOver correctly interpret the button's action.
