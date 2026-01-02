# Security Practices

This document outlines the security practices and guidelines for the Discipleship Journal project.

## 🔐 Secret Management

### Never Commit Secrets
- **Never** commit API keys, passwords, or tokens to the repository
- **Always** use environment variables or secret managers
- **Use** `.env` files for local development (add to `.gitignore`)

### Environment Variables
Use environment variables for all sensitive configuration:

```bash
# Local development (.env file)
DATABASE_URL=postgres://user:password@localhost:5432/db
BIBLE_API_KEY=your_actual_key_here
GOOGLE_CLOUD_PROJECT=your-project-id

# Production (Cloud Run/Firebase)
# Set via Cloud Run environment variables or Secret Manager
```

### Secret Scanning
We use pre-commit hooks to prevent accidental secret commits:

```bash
# Install pre-commit hooks
./scripts/setup-pre-commit.sh

# Run manually
pre-commit run --all-files
```

## 🛡️ Pre-commit Hooks

### Installed Hooks
1. **detect-secrets** - Scans for 100+ secret patterns
2. **gitleaks** - Advanced secret detection
3. **Custom pattern checks** - Common secret patterns
4. **Basic code quality** - Trailing whitespace, YAML validation, etc.

### To Skip Hooks (Emergency Only)
```bash
git commit --no-verify  # NOT RECOMMENDED
```

## 🚨 Emergency Procedures

### If You Accidentally Commit a Secret
1. **Immediately** rotate the compromised secret
2. **Rewrite git history** to remove the secret:
   ```bash
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch PATH_TO_FILE" \
     --prune-empty --tag-name-filter cat -- --all
   ```
3. **Force push** to all branches:
   ```bash
   git push origin --force --all
   git push origin --force --tags
   ```
4. **Notify team members** to reset their local repos

### Secret Rotation Schedule
- API Keys: Rotate every 90 days or when compromised
- Database passwords: Rotate every 180 days
- Service accounts: Rotate annually

## 📁 File Structure Safety

### Always Gitignore
```
.env
.env.local
.env.production
*.key
*.pem
*.crt
service-account.json
firebase-adminsdk-*.json
```

### Safe Defaults in Code
```yaml
# docker-compose.yml - SAFE
BIBLE_API_KEY: ${BIBLE_API_KEY:-REPLACE_WITH_YOUR_BIBLE_API_KEY}

# docker-compose.yml - UNSAFE
BIBLE_API_KEY: ${BIBLE_API_KEY:-real_secret_key_here}
```

## 🔍 Code Review Security Checklist

### Before Merging PRs
- [ ] No hardcoded secrets
- [ ] Environment variables used properly
- [ ] No exposed API endpoints without auth
- [ ] Input validation present
- [ ] Error messages don't leak sensitive info
- [ ] Dependencies are up-to-date and secure

### Security-Focused PR Template
```markdown
## Security Considerations
- [ ] No secrets in code
- [ ] Environment variables documented
- [ ] Authentication/authorization reviewed
- [ ] Input validation implemented
- [ ] Error handling secure
```

## 🛠️ Tools & Automation

### CI/CD Security
- GitHub Actions secrets for deployment
- Automated dependency vulnerability scanning
- Regular security audits

### Monitoring
- Log sensitive operations (auth attempts, data access)
- Monitor for unusual patterns
- Set up alerts for security events

## 📚 Resources

### Learning
- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [GitHub Secret Scanning](https://docs.github.com/en/code-security/secret-scanning)
- [Pre-commit Hooks](https://pre-commit.com/)

### Tools
- [detect-secrets](https://github.com/Yelp/detect-secrets)
- [gitleaks](https://github.com/gitleaks/gitleaks)
- [trivy](https://github.com/aquasecurity/trivy) (container scanning)

---

**Remember**: Security is everyone's responsibility. When in doubt, ask for a security review before merging changes.