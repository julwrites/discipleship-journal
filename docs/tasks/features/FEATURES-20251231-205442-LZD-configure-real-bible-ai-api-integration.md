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
1. **Verify current configuration**: The project already uses a custom BibleAIAPI service
   - **API URL**: `https://mock-bible-api.example.com`
   - **API Key**: `REPLACE_WITH_YOUR_BIBLE_API_KEY`
   - **Location**: Configured in `.env.local` and `docker-compose.yml`

2. **Fix API integration issues** (prerequisite: see task `DOMAIN-20260101-100010-IWD`):
   - Correct request/response handling in `api/services/bible_ai.go`
   - Test with real API using cURL examples

3. **Configure production environment**:
   - Ensure `BIBLE_API_URL` and `BIBLE_API_KEY` are set in Cloud Run
   - Use Google Secret Manager for production secrets

4. **Test integration**:
   - Use cURL commands to verify API connectivity
   - Test GetPassage endpoint with real Bible references
   - Test AI chat endpoints with sample prompts

5. **Alternative service** (if needed):
   - Scripture API Bible: `https://api.scripture.api.bible`
   - Requires separate API key registration
   - May require different request/response format

## Bible AI API Documentation
See task `DOMAIN-20260101-100010-IWD` for complete API documentation including:
- Request/response formats
- Example cURL commands
- Authentication details
- Error handling patterns

## Acceptance Criteria
- [ ] **API integration fixed** (dependency: task `DOMAIN-20260101-100010-IWD`)
- [ ] **Current BibleAIAPI service verified** and working
- [ ] **Production environment** configured with correct secrets
- [ ] **Real Bible passages** returned in API responses
- [ ] **AI chat** uses real Bible context
- [ ] **Documentation updated** with actual configuration (not scripture.api.bible)

## Current Status
- ✅ **Credentials already exist**: API key configured in `.env.local`
- ✅ **API URL configured**: BibleAIAPI service URL set
- ❌ **Integration broken**: Request/response handling needs fixing
- ❌ **Production verification**: Need to test with real API

## Prerequisite Task
Complete `DOMAIN-20260101-100010-IWD` first to fix the API integration issues before configuring production deployment.
