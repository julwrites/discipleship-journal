# Troubleshooting Guide

## Database Migrations

### Dirty Database State ("Dirty database version 1")

**Symptom:**
Database migrations fail with an error like:
```
failed to apply migrations: Dirty database version 1. Fix and force version.
```

**Cause:**
A previous migration attempt failed mid-execution, leaving the database in an inconsistent ("dirty") state. For version 1 (initial migration), this often happens because the database user lacks permissions to create extensions (e.g., `uuid-ossp`) or schemas.

**Resolution:**
1.  **Run the fix script:**
    Use the provided helper script to clean the dirty state and attempt to create the necessary extension.
    ```bash
    ./scripts/fix_staging_db.sh
    ```
    You will be prompted for:
    -   **Cloud SQL Instance Connection Name:** (e.g., `discipleship-journal-52a2c:asia-southeast1:discipleship-journal-postgresql`)
    -   **Database Name:** (e.g., `discipleship_journal_staging`)
    -   **User:** Ensure you use a **superuser** (like `postgres`) to grant permissions or create extensions if the application user cannot.

2.  **Verify Permissions:**
    Ensure the application user (defined in secrets) has the necessary permissions on the `public` schema and can create extensions if required.
    Often, the `uuid-ossp` extension must be created by a superuser once.

### "Failed to parse system prompts JSON"

**Symptom:**
Application logs show:
```
level=WARN msg="Failed to parse system prompts JSON" error="invalid character '\\n' in string literal"
```

**Cause:**
The `LLM_SYSTEM_PROMPTS` secret in Google Secret Manager contains unescaped newlines in the JSON string. JSON strings cannot contain raw newlines; they must be escaped as `\n`.

**Resolution:**
1.  Go to Google Secret Manager.
2.  Edit the `LLM_SYSTEM_PROMPTS` secret.
3.  Ensure the content is a valid JSON string on a single line, or properly escaped.
    **Incorrect:**
    ```json
    "You are a helpful assistant.
    Do your best."
    ```
    **Correct:**
    ```json
    "You are a helpful assistant.\nDo your best."
    ```
