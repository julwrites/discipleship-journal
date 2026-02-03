# Bible AI API Setup Guide

This guide explains how to configure the real Bible AI integration for the Discipleship Journal application. The application supports a "Bible AI" service that provides Bible passages and AI-powered chat capabilities.

## Architecture

The backend (`api/`) communicates with an external Bible AI service. This service is configured via two main environment variables or secrets:

1.  `BIBLE_API_URL`: The base URL of the Bible AI service.
2.  `BIBLE_API_KEY`: The authentication key for the service.

The application uses a **Secret Loader** (`api/services/secrets.go`) that prioritizes Google Secret Manager in production but falls back to environment variables (`.env`) for local development.

## Local Development

For local development, you can configure these values in your `.env` file (or `.env.local`):

```bash
# .env
BIBLE_API_URL=https://mock-bible-api.example.com
BIBLE_API_KEY=your-dev-api-key
```

If you do not provide these, the application will fallback to a **Mock Client** (`api/services/bible_ai_mock.go`), which returns hardcoded responses useful for UI testing without external dependencies.

## Production Configuration (Google Cloud)

In production (Cloud Run), we strictly use **Google Secret Manager** for sensitive keys.

### 1. Create Secrets in Google Secret Manager

You need to create two secrets in the Google Cloud Project where the application is deployed.

**BIBLE_API_URL**
```bash
echo -n "https://api.your-bible-ai-service.com" | gcloud secrets create BIBLE_API_URL --data-file=-
```

**BIBLE_API_KEY**
```bash
echo -n "your-production-secret-key" | gcloud secrets create BIBLE_API_KEY --data-file=-
```

### 2. Grant Permissions

Ensure the **Runtime Service Account** (the service account your Cloud Run service uses, e.g., `service-account@project-id.iam.gserviceaccount.com`) has the **Secret Manager Secret Accessor** role (`roles/secretmanager.secretAccessor`).

```bash
gcloud projects add-iam-policy-binding YOUR_PROJECT_ID \
  --member="serviceAccount:YOUR_SERVICE_ACCOUNT_EMAIL" \
  --role="roles/secretmanager.secretAccessor"
```

### 3. Deploy

When the backend starts, it will automatically attempt to fetch these secrets from Secret Manager. No additional environment variable configuration is needed in the `docker-compose.yml` or Cloud Run service definition for these values, provided the `APP_ENV` is set to `production` (or default behavior where secrets are checked).

## Verification

To verify the integration is working in production:

1.  **Check Logs**: The application logs "Using RealBibleAIClient" on startup if the secrets are successfully loaded.
    ```json
    {"level":"INFO", "msg":"Using RealBibleAIClient", "url":"..."}
    ```
2.  **Test Endpoint**: Make a request to the passage endpoint.
    ```bash
    curl "https://your-api-url.com/api/bible/passage?ref=John+3:16"
    ```
    You should receive a JSON response containing the verse text.

## Alternative Providers

The current implementation is optimized for the custom Bible AI API. If you wish to use a different provider (e.g., `api.scripture.api.bible`), you may need to:

1.  Update `BIBLE_API_URL` to the new provider's endpoint.
2.  Update `BIBLE_API_KEY` with the new key.
3.  **Code Change**: Modify `api/services/bible_ai.go` (specifically `GetPassage` and `ChatCompletion`) to handle the different request/response JSON structure of the new provider.

## Troubleshooting

-   **"Bible API not configured"**: This error means neither environment variables nor Secret Manager secrets were found. The app might be using the Mock Client, or `BIBLE_API_URL` is empty.
-   **"invalid response format"**: The API URL might be correct, but the service is returning a response structure that the code doesn't expect. Check `api/services/bible_ai.go` against the API documentation.
