---
id: DOMAIN-20260101-100010-IWD
status: completed
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
- Default URL: `https://mock-bible-api.example.com`
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

## Bible AI API Documentation (Found)

Based on exploration of `~/julwrites/bibleaiapi` and `~/julwrites/ScriptureBot`:

### **API Overview**
- **Single Endpoint**: `POST /query` (both GetPassage and ChatCompletion use this)
- **Authentication**: `X-API-KEY` header
- **Base URL**: `https://mock-bible-api.example.com` (production)
- **Local URL**: `http://localhost:8080`

### **Request Format**
```json
{
  "query": {
    "verses": ["John 3:16"],  // For passage retrieval
    "words": ["grace"],        // For word search
    "prompt": "Explain..."     // For Q&A/chat
  },
  "context": {
    "history": [],            // Optional conversation history
    "schema": "",             // Optional schema for structured output
    "verses": [],             // Optional context verses
    "words": [],              // Optional context words
    "user": {
      "version": "ESV"        // Bible translation version
    }
  }
}
```

**Important**: Exactly one of `verses`, `words`, or `prompt` must be specified in the `query` object.

### **Response Formats**

#### **Verse Response** (for `verses` query)
```json
{
  "verse": "John 3:16 (ESV) For God so loved the world..."
}
```

#### **Word Search Response** (for `words` query)
```json
[
  {
    "verse": "Romans 3:24",
    "url": "https://classic.biblegateway.com/passage/?search=Romans+3%3A24&version=ESV"
  }
]
```

#### **LLM/Open Query Response** (for `prompt` query)
```json
{
  "text": "Jesus fed 5,000 men, not including women and children.",
  "references": [
    {
      "verse": "Matthew 14:21",
      "url": "https://classic.biblegateway.com/passage/?search=Matthew+14%3A21&version=ESV"
    }
  ]
}
```

#### **Error Response**
```json
{
  "error": {
    "code": 401,
    "message": "Invalid API Key"
  }
}
```

### **Example cURL Commands**

**Verse Retrieval:**
```bash
curl -X POST http://localhost:8080/query \
  -H "X-API-KEY: secretapikey" \
  -d '{"query": {"verses": ["John 3:16"]}}'
```

**Word Search:**
```bash
curl -X POST http://localhost:8080/query \
  -H "X-API-KEY: secretapikey" \
  -d '{"query": {"words": ["Grace"]}}'
```

**LLM Prompt with Context:**
```bash
curl -X POST http://localhost:8080/query \
  -H "X-API-KEY: secretapikey" \
  -d '{"query": {"prompt": "Explain this verse"}, "context": {"verses": ["John 3:16"], "user": {"version": "ESV"}}}'
```

### **Key Files for Reference**

#### **BibleAIAPI Repository** (`~/julwrites/bibleaiapi/`)
- `docs/api/openapi.yaml` - Complete OpenAPI specification
- `pkg/client/` - Official Go client library
- `scripts/test_api.sh` - Test scripts with example commands
- `cmd/server/main.go` - Server implementation showing expected requests

#### **ScriptureBot Repository** (`~/julwrites/ScriptureBot/`)
- `pkg/app/api_client.go` - Client implementation consuming the API
- `pkg/app/api_models.go` - Request/response struct definitions
- `pkg/app/passage.go` - Passage retrieval implementation (lines 204-213)
- `pkg/app/ask.go` - Q&A implementation (lines 32-45)

### **Current Configuration in This Project**
- **`.env.local`**: Contains correct API URL and key
- **`docker-compose.yml`**: Sets default Bible API values
- **`api/services/bible_ai.go`**: Current implementation (needs fixing)

### **Issues with Current Implementation**

1. **GetPassage Method** (lines 35-78):
   - ✅ Correct: Uses `POST /query` with `{"query": {"verses": [reference]}}`
   - ❌ Issue: Response parsing expects raw map, but API returns `{"verse": "..."}`

2. **ChatCompletion Method** (lines 81-197):
   - ❌ Issue: Transforms input incorrectly
   - Current: `{"prompt": "...", "verses": [...], "themes": [...], "context": "..."}`
   - Should be: `{"query": {"prompt": "..."}, "context": {"verses": [...], "user": {"version": "ESV"}}}`
   - ❌ Issue: Complex fallback parsing suggests wrong response structure

## Next Steps (Updated with API Documentation)

1. **✅ Examine ScriptureBot and BibleAIAPI references** in `~/julwrites/` - COMPLETED
2. **Fix GetPassage method** in `api/services/bible_ai.go`:
   - Update response parsing to extract `result["verse"]` instead of raw map
   - Add proper error handling for API error responses
3. **Fix ChatCompletion method** in `api/services/bible_ai.go`:
   - Transform input correctly: `{"query": {"prompt": "..."}, "context": {"verses": [...], "user": {"version": "ESV"}}}`
   - Simplify response parsing to extract `result["text"]` and `result["references"]`
   - Remove complex fallback logic (OpenAI format, etc.)
4. **Test with real API credentials** using example cURL commands above
5. **Consider using Go client library** from `~/julwrites/bibleaiapi/pkg/client/` for consistency

## Specific Fixes Required

### **GetPassage Fix** (`api/services/bible_ai.go:35-78`)
```go
// Current response parsing (line ~70):
// return result, nil

// Should be:
verse, ok := result["verse"].(string)
if !ok {
    return nil, fmt.Errorf("invalid response format: missing 'verse' field")
}
return verse, nil
```

### **ChatCompletion Fix** (`api/services/bible_ai.go:81-197`)
```go
// Current payload transformation (lines ~110-130):
// payload := map[string]interface{}{
//     "prompt": prompt,
//     "verses": verses,
//     "themes": themes,
//     "context": context,
// }

// Should be:
queryPayload := map[string]interface{}{
    "query": map[string]interface{}{
        "prompt": prompt,
    },
    "context": map[string]interface{}{
        "verses": verses,
        "user": map[string]interface{}{
            "version": "ESV", // or get from user preferences
        },
    },
}

// Response parsing should extract:
// text := result["text"].(string)
// references := result["references"].([]interface{})
```

## Blockers (Resolved)
- ✅ **API Documentation**: Found in BibleAIAPI repository
- ✅ **Sample Code References**: Found in ScriptureBot repository
- ✅ **API URL and Authentication**: Verified in `.env.local` and BibleAIAPI docs
- ✅ **Example Requests/Responses**: Provided in cURL examples above

## Updated Implementation Plan

1. **Update `bible_ai.go`** with correct request/response handling
2. **Test GetPassage endpoint** with `curl -X GET "http://localhost:8080/api/bible/passage?ref=John+3:16"`
3. **Test Chat endpoints** with sample prompts
4. **Verify response formats** match expected structures
5. **Update tests** to reflect new API behavior

## Impact
- **Broken Feature**: Bible passage lookup and AI chat not working
- **User Experience**: API errors or incorrect responses
- **Data Integrity**: Might receive wrong Bible passages or AI responses

## Implementation Notes (2026-01-04)

### Root Cause Identified
1. **Frontend/Backend Response Field Mismatch**: Frontend expected `text` or `content` field, but backend returned `verse` field (as per Bible AI API response format).
2. **Mock vs Real API Inconsistency**: Mock client returns `text` field, real API returns `verse` field, causing frontend to show "Passage found but no text returned" when using real API.

### Fixes Applied
1. **Frontend (`web/src/pages/NoteEditor.tsx`)**: Updated to check `verse` first, then `text`, then `content`.
2. **Backend Real Client (`api/services/bible_ai.go`)**: Added `text` field to response (duplicate of `verse`) for compatibility.
3. **Backend Mock Client (`api/services/bible_ai_mock.go`)**: Added `verse` field (duplicate of `text`) for consistency.
4. **Tests Updated**: Updated frontend test mocks to include both `verse` and `text` fields.

### Additional Improvements
- Added logging in `GetPassage` method to aid debugging when verse is empty.
- Added fallback parsing to check for `text` field in API response for robustness.

### Verification
- Backend tests pass.
- Frontend tests pass.
- Frontend build succeeds.
- Response format now compatible with both mock and real API.

### HTML to Markdown Conversion Fix (2026-01-04)

#### Root Cause Identified
1. **Bible AI API returns HTML-formatted verses**: The API returns verses with HTML tags (`<p>`, `<b>`, `<i>`, `<sup>`, etc.) for formatting.
2. **Frontend preview showed raw HTML tags**: The preview dialog displayed HTML tags instead of formatted text because the backend was returning HTML without conversion.

#### Fixes Applied
1. **Added `parsePassageFromHTML` function** in `api/services/bible_ai.go`:
   - Converts common Bible HTML tags to Markdown equivalents:
     - `<b>`, `<strong>` → `**bold**`
     - `<i>`, `<em>` → `*italic*`
     - `<sup>` → `^superscript^`
     - `<br>` → newline
     - `<p>` → paragraph separation with blank lines
     - `<h1>-<h4>` → bold headings with proper spacing
   - Handles nested tags and unknown elements gracefully.
2. **Updated `GetPassage` method** to convert HTML to Markdown before returning verses.
3. **Enhanced frontend preview** to use `ReactMarkdown` for proper rendering of Markdown formatting.
4. **Added comprehensive tests** for HTML parsing function with various tag combinations.

#### Verification
- Backend tests pass (including new HTML parsing tests).
- Frontend tests pass with updated mock responses.
- HTML tags are now converted to Markdown formatting.
- Preview dialog displays formatted text instead of raw HTML.
