# Secret Management

The application uses a dual-strategy for managing configuration secrets:

1.  **Google Secret Manager**: Primary source for production/staging environments.
2.  **Environment Variables**: Fallback for local development or if Secret Manager is unavailable.

## Secret Loading Strategy

The application initializes a `SecretLoader` at startup. When a configuration value is needed, it follows this precedence:

1.  **Google Secret Manager**: Attempts to fetch the secret with the given name from the GCP Project defined in `GOOGLE_CLOUD_PROJECT`.
2.  **Environment Variable**: If the secret doesn't exist in Secret Manager or the client fails, it looks for an environment variable with the same name.

## Required Secrets

The following secrets should be defined in Google Secret Manager for the production environment:

### Database Configuration
| Secret Name | Description | Example |
|-------------|-------------|---------|
| `DB_USERNAME` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `secret_password` |
| `DB_NAME` | Database name | `discipleship_journal` |
| `DB_HOST` | Database host (or socket path for Cloud SQL) | `127.0.0.1` |
| `DB_PORT` | Database port | `5432` |
| `CLOUD_SQL_INSTANCE` | Cloud SQL Connection Name | `project:region:instance` |

### External APIs
| Secret Name | Description |
|-------------|-------------|
| `BIBLE_API_URL` | Base URL for the Bible API (e.g. RealBibleAI) |
| `BIBLE_API_KEY` | API Key for the Bible API |
| `LLM_SYSTEM_PROMPTS` | JSON string containing system prompts |

### Security
| Secret Name | Description | Example |
|-------------|-------------|---------|
| `CORS_ALLOWED_ORIGINS` | Comma-separated list of allowed origins | `https://journal.example.com,https://staging.journal.example.com` |

## Local Development

For local development, you do not need to set up Google Secret Manager. Simply create a `.env` file in the `api/` directory with the required values. The `godotenv` library will load these as environment variables, which the `SecretLoader` will use as fallback.
