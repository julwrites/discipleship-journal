---
id: DOMAIN-20251213-092126-EZZ
status: pending
title: Refactor Connections API
priority: medium
created: 2025-12-13 09:21:26
category: domain
dependencies:
type: task
---

# Refactor Connections API

Standardize the Connections API to use RESTful verbs. Current implementation uses query params for actions on DELETE/PUT.
Proposed changes:
- DELETE /api/connections/{id} (No action param)
- POST /api/connections/{id}/accept (New endpoint for accepting)
