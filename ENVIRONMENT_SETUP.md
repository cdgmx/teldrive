# Environment Variable Configuration

TelDrive now supports environment variable overrides for sensitive configuration values. This allows you to keep secrets out of your configuration files while maintaining security.

## Quick Setup

### Option 1: Using .env file (Recommended for production)

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your actual values:**
   ```bash
   # Required: JWT secret for authentication (minimum 32 characters)
   TELDRIVE_JWT_SECRET=your-super-secret-jwt-key-here-at-least-32-chars
   
   # Required: Your Telegram usernames (comma-separated)
   TELDRIVE_JWT_ALLOWED_USERS=YourTelegramUsername1,YourTelegramUsername2
   
   # Required: Encryption key for uploads (exactly 32 characters)
   TELDRIVE_TG_UPLOADS_ENCRYPTION_KEY=your-32-character-encryption-key-here
   
   # Optional: Redis password (if using Redis)
   TELDRIVE_CACHE_REDIS_PASS=your-redis-password-if-needed
   
   # Optional: Database password (for extra security)
   POSTGRES_PASSWORD=your-secure-database-password
   ```

3. **Start with Docker Compose:**
   ```bash
   docker-compose up -d
   ```

### Option 2: Quick start without .env (Development only)

For quick testing, you can start immediately:
```bash
docker-compose up -d
```

⚠️ **Warning**: This uses default placeholder values that are **NOT SECURE** for production. You'll see values like:
- JWT Secret: `CHANGE-THIS-JWT-SECRET-IN-PRODUCTION-MIN-32-CHARS`
- Allowed Users: `YourTelegramUsername`
- Encryption Key: `CHANGE-THIS-32-CHAR-ENCRYPTION-KEY`

**Always create a proper `.env` file for production use!**

## Environment Variables

| Variable | Description | Required | Example |
|----------|-------------|----------|---------|
| `TELDRIVE_JWT_SECRET` | JWT signing secret | ✅ | `my-super-secret-jwt-key-32-chars-min` |
| `TELDRIVE_JWT_ALLOWED_USERS` | Comma-separated usernames | ✅ | `user1,user2,user3` |
| `TELDRIVE_TG_UPLOADS_ENCRYPTION_KEY` | Upload encryption key (32 chars) | ✅ | `abcdef1234567890abcdef1234567890` |
| `TELDRIVE_CACHE_REDIS_PASS` | Redis password | ❌ | `redis-password` |
| `POSTGRES_PASSWORD` | Database password | ❌ | `secure-db-password` |

## How It Works

1. **Config Priority:** Environment variables override values in `config.toml`
2. **Naming Convention:** Config keys like `jwt.secret` become `TELDRIVE_JWT_SECRET`
3. **Fallbacks:** If env var is not set, uses value from `config.toml`
4. **Security:** Sensitive values are never stored in git-tracked files

## Deployment Platforms

### Docker Compose
Already configured! Just create your `.env` file.

### Other Platforms
Set these environment variables in your deployment platform:
- Heroku: Use `heroku config:set`
- Railway: Set in environment variables section
- DigitalOcean App Platform: Set in environment variables
- Kubernetes: Use ConfigMaps/Secrets

## Security Notes

- ✅ `.env` is git-ignored (never committed)
- ✅ `config.toml` contains no secrets (safe to commit)
- ✅ Environment variables take precedence
- ✅ Works with any deployment platform
- ⚠️ Keep your `.env` file secure and backed up separately