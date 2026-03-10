# Database Configuration Guide

This document explains how to configure database connections for the Discipleship Journal application, both for local development and production (Google Cloud SQL).

## Configuration Method

We use individual environment variables for each database component (no single DATABASE_URL).

## Environment Variables

### Required:
- `DB_USERNAME` - Database username
- `DB_PASSWORD` - Database password
- `DB_NAME` - Database name

### Optional:
- `DB_HOST` - Database hostname (default: `localhost`)
- `DB_PORT` - Database port (default: `5432`)
- `CLOUD_SQL_INSTANCE` - Cloud SQL instance name (format: `project:region:instance`)

## Connection String Formats

### Local Development (TCP):
```
postgres://username:password@host:port/dbname?sslmode=disable
```
Example: `postgres://user:password@db:5432/discipleship_journal?sslmode=disable`

### Production - Cloud SQL (Unix Socket):
```
postgres://username:password@/dbname?host=/cloudsql/project:region:instance&sslmode=disable
```
Example: `postgres://myuser:mypassword@/discipleship_journal?host=/cloudsql/myproject:asia-southeast1:myinstance&sslmode=disable`

## Where to Find Cloud SQL Values

### 1. Google Cloud Console → SQL
1. Go to [Google Cloud Console](https://console.cloud.google.com/sql)
2. Select your instance
3. Click on the instance name to view details

### 2. Instance Connection Name
- **Format**: `project-id:region:instance-name`
- **Location**: Instance details page → "Connection name"
- **Example**: `myproject-123456:asia-southeast1:discipleship-journal-db`

### 3. Database Credentials
- **Username**: Set during instance creation or in "Users" tab
- **Password**: Set during user creation (can be reset in "Users" tab)
- **Database Name**: Created in "Databases" tab

### 4. Service Account Permissions
Ensure the Cloud Run service account has:
- `Cloud SQL Client` role
- Or use the service account email in Cloud SQL → "Connections" → "Add network"

## Configuration Examples

### Local Development (.env.local):
```bash
# Individual components (recommended)
DB_USERNAME=user
DB_PASSWORD=password
DB_NAME=discipleship_journal
DB_HOST=db
DB_PORT=5432

# Or legacy format
# DATABASE_URL=postgres://user:password@db:5432/discipleship_journal?sslmode=disable
```

### Production (GitHub Secrets):
Add these secrets in GitHub → Settings → Secrets and variables → Actions:

| Secret Name | Value Example | Description |
|-------------|---------------|-------------|
| `DB_USERNAME` | `myproductionuser` | Database username |
| `DB_PASSWORD` | `SecurePassword123!` | Database password |
| `DB_NAME` | `discipleship_journal` | Database name |
| `CLOUD_SQL_INSTANCE` | `myproject:asia-southeast1:myinstance` | Cloud SQL instance |

### Docker Compose Configuration:
The `docker-compose.yml` file uses individual components with defaults for local development.

## Verification Steps

### 1. Test Connection Locally:
```bash
# Set environment variables
export DB_USERNAME=user
export DB_PASSWORD=password
export DB_NAME=discipleship_journal
export DB_HOST=localhost
export DB_PORT=5432

# Run the application
cd api
go run main.go
```

### 2. Check Health Endpoint:
```bash
curl http://localhost:8080/health
# Should return "OK" if database is connected
```

### 3. Verify Cloud Run Deployment:
After deployment, check:
```bash
# Check health endpoint
curl https://api.yourdomain.com/health

# Check Cloud Run logs
gcloud run logs read [SERVICE_NAME] --region=[REGION]
```

## Troubleshooting

### Common Issues:

1. **"DATABASE_CONNECTION_FAILED" in health check**
   - Check if all required environment variables are set
   - Verify Cloud SQL instance is running
   - Ensure Cloud Run service account has Cloud SQL Client role

2. **Authentication failed**
   - Verify username and password are correct
   - Check if user exists in Cloud SQL
   - Reset password if needed

3. **Connection timeout**
   - Verify Cloud SQL instance allows connections from Cloud Run
   - Check firewall rules and authorized networks
   - Ensure instance has public IP or proper VPC configuration

4. **Database does not exist**
   - Create the database in Cloud SQL
   - Run migrations: `docker-compose run migrator`


## Security Notes

- Never commit passwords or connection strings to version control
- Use GitHub Secrets for production credentials
- Regularly rotate database passwords
- Use minimal privileges for database users
- Enable SSL/TLS for production connections (sslmode=require)
