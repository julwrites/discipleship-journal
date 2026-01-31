---
date: 2026-01-31
title: "Use import.meta.env.DEV for development-only logging in frontend"
tags: []
created: 2026-01-31 14:16:13
---

Wrap console.log statements with import.meta.env.DEV checks to reduce production console noise while preserving debugging capability. Service worker InvalidStateError should be suppressed as it occurs when unregistering active workers.
