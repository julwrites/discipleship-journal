---
id: DOMAIN-20260101-100010-IWD
status: pending
title: Fix Bible AI API integration: incorrect request handling for get passage and ask AI
priority: medium
created: 2026-01-01 10:00:10
category: domain
dependencies: 
type: bug
---

# Fix Bible AI API integration: incorrect request handling for get passage and ask AI

## Problem
Bible AI API integration is not working correctly for both:
1. **Get Bible Passage** (`GET /api/bible/passage?ref=John+3:16`)
2. **Ask AI** (`POST /api/ai/ask`) and **Chat with AI** (`POST /api/chat`)

The API requests are likely using incorrect endpoints, payload formats, or response parsing.

## Current Implementation Analysis

### **File: `api/services/bible_ai.go`**

#### **GetPassage Method** (Line 35-78)
- **Endpoint**: POST to `{API_URL}/query`
- **Payload**: `{"query": {"verses": [reference]}}`
- **Headers**: `X-API-KEY`, `Content-Type: application/json`
- **Response Handling**: Returns raw result map

#### **ChatCompletion Method** (Line 81-197)
- **Endpoint**: POST to `{API_URL}/query` (same as GetPassage)
- **Payload Transformation**:
  - Input: `{"prompt": "...", "verses": [...], "themes": [...], "context": "..."}`
  - Output: `{"query": {"prompt": "..."}, "context": {"user": {"version": "ESV"}, "verses": [...]}}`
- **Response Handling**: Tries multiple formats:
  1. OpenAI format: `result["choices"][0]["message"]["content"]`
  2. `result["text"]`
  3. `result["response"]`
  4. Raw result

### **File: `api/handlers/bible.go`**
- Simple wrapper that calls `Client.GetPassage()` and returns result

### **File: `api/handlers/chat.go`**
- **ChatWithAI**: Calls `Client.ChatCompletion()` with passage context
- **AskAI**: Calls `Client.ChatCompletion()` with general context

## Potential Issues Identified

### 1. **Endpoint Consistency**
- Both `GetPassage` and `ChatCompletion` use `/query` endpoint
- This might be incorrect; different operations might need different endpoints

### 2. **Payload Format**
- The payload structure might not match what the actual Bible AI API expects
- Need to verify against API documentation

### 3. **Response Parsing**
- Complex fallback logic suggests uncertainty about response format
- Might be parsing wrong fields or expecting wrong structure

### 4. **API Configuration**
- Default URL: `https://bible-api-service-779024060388.asia-southeast1.run.app`
- Might be wrong or outdated
- Alternative reference in `.env`: `https://api.scripture.api.bible`

## Required Information for Fix

To properly fix this issue, we need:

1. **API Documentation** for the Bible AI service
2. **Sample requests/responses** for:
   - Getting a Bible passage
   - Chat/completion requests
3. **Correct endpoint URLs**
4. **Required headers and authentication**

## References Mentioned
- **ScriptureBot**: Sample usage reference
- **BibleAIAPI**: Documented APIs
- Location: `~/julwrites/` (need specific paths)

## Next Steps

1. **Examine ScriptureBot and BibleAIAPI references** in `~/julwrites/`
2. **Compare current implementation** with actual API requirements
3. **Fix endpoint URLs and payload formats**
4. **Simplify response parsing** based on actual response structure
5. **Test with real API credentials**

## Blockers
- Need access to API documentation or sample code references
- Need to verify current API URL and authentication method
- Need example requests/responses to understand correct format

## Impact
- **Broken Feature**: Bible passage lookup and AI chat not working
- **User Experience**: API errors or incorrect responses
- **Data Integrity**: Might receive wrong Bible passages or AI responses
