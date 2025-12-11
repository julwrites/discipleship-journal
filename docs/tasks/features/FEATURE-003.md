---
id: FEATURE-003
status: completed
title: AI Chat Integration
priority: high
created: 2025-12-11 06:09:10
category: unknown
type: task
---

# AI Chat Integration

### Description
Integrate BibleAIAPI for chat and context-aware questions.

### Acceptance Criteria
- [x] API: Integration with BibleAIAPI.
- [x] API: Endpoint for Chat (`POST /api/chat`).
- [x] API: Endpoint for Asking AI (`POST /api/ai/ask`).
- [x] UI: Chat Interface.
- [x] UI: Context selection (Bible passages).
- [x] Logic: Convert Chat response to new Journal Note.

### Implementation Status
- ✅ Backend handlers implemented and verified (`api/handlers/chat.go`).
- ✅ Frontend page scaffolded (`ChatPage.tsx`).
