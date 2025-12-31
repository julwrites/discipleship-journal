---
id: FEATURES-20251231-205442-LZD
status: pending
title: Configure real Bible AI API integration
priority: high
created: 2025-12-31 20:54:42
category: features
dependencies:
type: task
---

# Configure real Bible AI API integration

## Problem
The application currently only returns mock data from the Bible AI API instead of real Bible passages and AI responses.

## Requirements
1. Obtain real Bible API credentials (e.g., Scripture API Bible key)
2. Configure environment variables for production
3. Test real API integration
4. Update documentation with configuration steps

## Steps
1. Sign up for Bible API service (e.g., scripture.api.bible)
2. Obtain API key
3. Update Cloud Run environment variables:
   - `BIBLE_API_URL`: https://api.scripture.api.bible
   - `BIBLE_API_KEY`: [actual API key]
4. Test integration locally with real API
5. Deploy to production
6. Verify real Bible passages are returned

## Acceptance Criteria
- [ ] Real Bible API credentials obtained
- [ ] Environment variables configured in production
- [ ] Real Bible passages returned in API responses
- [ ] AI chat uses real Bible context
- [ ] Documentation updated with configuration guide
