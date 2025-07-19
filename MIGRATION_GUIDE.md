# 📦 TelDrive Migration Guide - Simplified Deployment

This guide helps you migrate from the previous multi-variant Docker Compose approach to the new simplified two-mode deployment.

## 🎯 What Changed

### Before (5 variants)
- `docker-compose.local.yml` - Local builds
- `docker-compose.registry.yml` - Registry images  
- `docker-compose.custom-ui.yml` - Custom UI builds
- `docker-compose.dev.yml` - Development
- Offline mode support

### After (2 modes)
- `docker-compose.yml` - **Production** (builds from source)
- `docker-compose.dev.yml` - **Development** (optimized workflow)

## 🔄 Migration Commands

### From docker-compose.local.yml
```bash
# Old command
docker-compose -f docker-compose.local.yml up -d

# New command (same functionality)
docker-compose up -d
```

### From docker-compose.custom-ui.yml  
```bash
# Old command
docker-compose -f docker-compose.custom-ui.yml up -d

# New command (identical functionality)
docker-compose up -d
```

### From docker-compose.registry.yml
```bash
# Old registry-based approach
docker-compose -f docker-compose.registry.yml up -d

# New source-based approach
docker-compose up -d
```

### Development remains the same
```bash
# Still works as before
docker-compose -f docker-compose.dev.yml up
```

## 📁 File Migration

### Configuration Files
Your existing files work without changes:
```bash
# Keep your existing configuration
# config.toml - No changes needed
# storage.db - No changes needed

# Optional: Create development-specific files
cp config.toml config.dev.toml
cp storage.db storage.dev.db
```

### Environment Variables
```bash
# If you used GitHub tokens before
export GITHUB_TOKEN=your_token

# Works the same way
docker-compose up -d
```

## ⚙️ Configuration Migration

### Custom UI Repository
If you were using custom UI settings:

**Before:**
```yaml
# In docker-compose.custom-ui.yml
args:
  - UI_REPO_URL=https://github.com/yourusername/teldrive-ui.git
  - UI_REPO_BRANCH=main
```

**After:**
```yaml
# In docker-compose.yml
args:
  - UI_REPO_URL=https://github.com/yourusername/teldrive-ui.git
  - UI_REPO_BRANCH=main
```

### Build Arguments
All build arguments remain the same:
- `UI_REPO_URL` - Your UI repository URL
- `UI_REPO_BRANCH` - Branch to build from
- `GITHUB_TOKEN` - For private repositories
- `BUILD_MODE` - production/development

## 🚀 Benefits of Migration

### ✅ Simplified Approach
- **2 clear modes** instead of 5 confusing variants
- **Production vs Development** - obvious choice
- **Less maintenance** - fewer files to manage

### ✅ Better Performance
- **Optimized builds** - focus on 2 modes instead of 5
- **Better caching** - streamlined Docker layers
- **Faster development** - improved dev workflow

### ✅ Enhanced Control
- **Source-based builds** - always latest changes
- **No registry dependency** - build everything locally
- **Full customization** - modify any component

## 🔧 Troubleshooting Migration

### Build Time Concerns
**Issue**: "Builds take longer now"
```bash
# Solution: Builds are optimized and cached
# First build: ~4 minutes
# Subsequent builds: ~2 minutes (with cache)

# For faster development rebuilds
docker-compose -f docker-compose.dev.yml up --build
```

### Missing Registry Images
**Issue**: "I prefer pre-built images"
```bash
# Solution: Source builds provide more control
# Benefits:
# - Always latest UI changes
# - Full customization capability  
# - No dependency on registry availability
# - Consistent build environment
```

### Offline Mode Removed
**Issue**: "I need offline deployment"
```bash
# Solution: Build once, save image
docker-compose build
docker save teldrive_teldrive:latest > teldrive-image.tar

# On offline machine
docker load < teldrive-image.tar
docker-compose up -d
```

## 📋 Migration Checklist

### Pre-Migration
- [ ] Backup your `config.toml`
- [ ] Backup your `storage.db`
- [ ] Note your custom UI repository settings
- [ ] Stop existing deployment: `docker-compose -f <old-file> down`

### Migration Steps
- [ ] Pull latest TelDrive code: `git pull`
- [ ] Update UI repository settings in `docker-compose.yml` if needed
- [ ] Test new deployment: `docker-compose up`
- [ ] Verify UI loads correctly
- [ ] Verify data persistence (sessions, files)

### Post-Migration
- [ ] Remove old compose files (already done in new version)
- [ ] Update any scripts/documentation referencing old files
- [ ] Update CI/CD pipelines if applicable
- [ ] Share new commands with team members

## 🆘 Common Migration Issues

### Issue: Build Fails
```bash
# Clean rebuild
docker-compose build --no-cache

# Check build logs
docker-compose build --progress=plain
```

### Issue: UI Not Loading
```bash
# Verify UI repository access
git clone https://github.com/cdgmx/teldrive-ui.git

# Check build args
docker-compose config
```

### Issue: Data Loss
```bash
# Verify volume mounts
docker-compose config | grep volumes

# Check file permissions
ls -la storage.db config.toml
```

### Issue: Port Conflicts
```bash
# Check if ports are in use
netstat -tulpn | grep :8080

# Stop conflicting services
docker-compose -f <old-file> down
```

## 📞 Getting Help

### Documentation
- **Build Guide**: [TELDRIVE_BUILD_GUIDE.md](TELDRIVE_BUILD_GUIDE.md)
- **Main Documentation**: https://teldrive-docs.pages.dev

### Support Channels
- **GitHub Issues**: https://github.com/tgdrive/teldrive/issues
- **Discussions**: https://github.com/tgdrive/teldrive/discussions

### Quick Commands Reference
```bash
# Production deployment
docker-compose up -d

# Development environment  
docker-compose -f docker-compose.dev.yml up

# View logs
docker-compose logs -f

# Clean rebuild
docker-compose build --no-cache

# Stop deployment
docker-compose down
```

---

**Migration complete!** You now have a simplified, more maintainable TelDrive deployment. 🎉