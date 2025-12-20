---
id: FEATURES-20251220-092428-DIJ
status: pending
title: Passwordless Sign-in
priority: medium
created: 2025-12-20 09:24:28
category: features
dependencies: FOUNDATION-002
type: task
---

# Passwordless Sign-in

Implement passwordless sign-in using email links (Magic Links) via Firebase Authentication.

## Requirements
- Allow users to sign in by entering their email address.
- Send a verification link to the email.
- Complete the sign-in process when the user clicks the link.
- Handle deep links if necessary.
- Update `LoginPage.tsx` to include the Email Link option.
