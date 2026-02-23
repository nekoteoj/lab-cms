# Configuration Guide

This document describes all configuration options for Lab CMS. All configuration is done via environment variables.

## Quick Start

1. Copy the example configuration:
   ```bash
   cp configs/.env.example .env
   ```

2. Generate a secure session secret:
   ```bash
   openssl rand -base64 32
   ```

3. Edit `.env` and set the required values (at minimum `SESSION_SECRET` and `ROOT_ADMIN_PASSWORD`)

4. Run the application

## Configuration Reference

### Server Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | HTTP server port |
| `ENV` | `development` | Environment mode: `development` or `production` |

**Environment Modes:**
- **development**: Relaxed security rules, verbose logging allowed
- **production**: Strict security enforced, debug logging disabled

### Database Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `./data/lab-cms.db` | Path to SQLite database file |

### Session & Security

| Variable | Default | Description |
|----------|---------|-------------|
| `SESSION_SECRET` | *(required)* | Secret key for session signing (32+ chars recommended) |
| `SESSION_MAX_AGE` | `24` | Session lifetime in hours |
| `COOKIE_SECURE` | `false` (dev), `true` (prod) | HTTPS-only cookies |
| `COOKIE_HTTPONLY` | `true` | Prevent JavaScript cookie access |
| `COOKIE_SAMESITE` | `strict` | CSRF protection level |
| `CSRF_ENABLED` | `true` | Enable CSRF token validation |
| `TRUSTED_PROXIES` | *(empty)* | Comma-separated proxy IPs |

**Cookie SameSite Values:**
- `strict`: Most secure, cookies never sent cross-site
- `lax`: Cookies sent on top-level navigation (login flows)
- `none`: No protection (not recommended)

### Initial Admin Setup

| Variable | Default | Description |
|----------|---------|-------------|
| `ROOT_ADMIN_USERNAME` | `admin` | Initial admin username |
| `ROOT_ADMIN_PASSWORD` | *(required)* | Initial admin password (8+ chars) |

**Note:** These credentials create the first admin account on application startup. Change the password immediately after first login.

### File Uploads

| Variable | Default | Description |
|----------|---------|-------------|
| `UPLOAD_PATH` | `./uploads` | Directory for uploaded files |
| `MAX_UPLOAD_SIZE` | `10485760` (10MB) | Maximum upload size in bytes |

### Logging

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Log verbosity: `debug`, `info`, `warn`, `error` |

**Log Levels:**
- `debug`: All messages (development only)
- `info`: General operational messages
- `warn`: Warning messages
- `error`: Errors only

## Security Best Practices

### Production Checklist

Before deploying to production:

- [ ] Set `ENV=production`
- [ ] Generate strong `SESSION_SECRET` (32+ random characters)
- [ ] Set `COOKIE_SECURE=true` (requires HTTPS)
- [ ] Verify `CSRF_ENABLED=true`
- [ ] Set `COOKIE_SAMESITE=strict`
- [ ] Change from `LOG_LEVEL=debug` to `info` or higher
- [ ] Set strong `ROOT_ADMIN_PASSWORD` (8+ characters)
- [ ] Configure `TRUSTED_PROXIES` if behind a reverse proxy

### Generating Secrets

Generate cryptographically secure secrets:

```bash
# Session secret (44 characters base64 = 32 bytes)
openssl rand -base64 32

# Alternative: 64 character hex string
openssl rand -hex 32
```

### Environment File Security

⚠️ **Never commit `.env` files to version control!**

Add to `.gitignore`:
```
.env
*.env
.env.*
```

For production deployments, use:
- Docker secrets
- Kubernetes secrets
- Cloud provider secret managers (AWS Secrets Manager, Azure Key Vault, etc.)
- Systemd service environment files with restricted permissions

## Example Configurations

### Development

```env
PORT=8080
ENV=development
DATABASE_URL=./data/lab-cms.db
SESSION_SECRET=dev-secret-change-in-production
SESSION_MAX_AGE=24
COOKIE_SECURE=false
COOKIE_HTTPONLY=true
COOKIE_SAMESITE=strict
CSRF_ENABLED=true
ROOT_ADMIN_USERNAME=admin
ROOT_ADMIN_PASSWORD=devpass123
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
LOG_LEVEL=debug
```

### Production

```env
PORT=8080
ENV=production
DATABASE_URL=./data/lab-cms.db
SESSION_SECRET=YOUR_GENERATED_SECRET_HERE_MIN_32_CHARS
SESSION_MAX_AGE=24
COOKIE_SECURE=true
COOKIE_HTTPONLY=true
COOKIE_SAMESITE=strict
CSRF_ENABLED=true
TRUSTED_PROXIES=127.0.0.1,10.0.0.0/8
ROOT_ADMIN_USERNAME=admin
ROOT_ADMIN_PASSWORD=SecurePass123!
UPLOAD_PATH=./uploads
MAX_UPLOAD_SIZE=10485760
LOG_LEVEL=info
```

## Validation

The configuration is validated on startup. Common errors:

| Error | Solution |
|-------|----------|
| `SESSION_SECRET is required` | Generate and set a session secret |
| `SESSION_SECRET must be at least 32 characters in production` | Use a longer secret in production |
| `ROOT_ADMIN_PASSWORD is required` | Set an initial admin password |
| `LOG_LEVEL cannot be 'debug' in production` | Change to `info`, `warn`, or `error` |
| `CSRF_ENABLED cannot be false in production` | Set to `true` |

## Troubleshooting

### Session Issues
If users are logged out unexpectedly:
- Check `SESSION_MAX_AGE` value
- Verify `SESSION_SECRET` is consistent across app restarts
- Ensure `COOKIE_SECURE` matches your HTTPS setup

### File Upload Errors
- Verify `UPLOAD_PATH` directory exists and is writable
- Check `MAX_UPLOAD_SIZE` is sufficient for your files
- Ensure disk has available space

### CSRF Token Errors
- Ensure `CSRF_ENABLED` matches between config and form handling
- Verify cookies are being sent (check browser dev tools)
- Check `COOKIE_SAMESITE` isn't blocking requests

## Advanced Configuration

### Behind a Reverse Proxy

When running behind Nginx, Apache, or a load balancer:

```env
TRUSTED_PROXIES=127.0.0.1,10.0.0.1
COOKIE_SECURE=true
```

The `TRUSTED_PROXIES` setting ensures client IP addresses are correctly identified.

### Custom Upload Directory

For production, store uploads outside the application directory:

```env
UPLOAD_PATH=/var/www/lab-cms/uploads
```

Ensure the directory exists and has correct permissions:
```bash
mkdir -p /var/www/lab-cms/uploads
chmod 750 /var/www/lab-cms/uploads
```
