---
date: 2026-01-31
title: "CORS configuration requires Google Secret Manager update for staging frontend origin"
tags: []
created: 2026-01-31 13:33:41
---

The backend API uses CORS_ALLOWED_ORIGINS secret from Google Secret Manager to validate origins. Staging frontend origin (https://discipleship-journal-staging.firebaseapp.com) must be included in this comma-separated list. The secret is loaded at runtime via SecretLoader in main.go. Update with gcloud secrets commands.
