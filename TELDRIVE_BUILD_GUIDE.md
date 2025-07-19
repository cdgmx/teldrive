# 🚀 TelDrive Build Guide - Simplified Deployment

This guide covers the **simplified two-mode deployment approach** for TelDrive: **Production** and **Development**.

## 📋 Table of Contents

- [Quick Start](#-quick-start)
- [Production Deployment](#-production-deployment)
- [Development Environment](#-development-environment)
- [Configuration](#-configuration)
- [Troubleshooting](#-troubleshooting)
- [Migration Guide](#-migration-guide)

## ⚡ Quick Start

### Production (Recommended)
```bash
# Clone the repository
git clone https://github.com/tgdrive/teldrive.git
cd teldrive

# Start production deployment
docker-compose up -d
```

### Development
```bash
# Clone the repository
git clone https://github.com/tgdrive/teldrive.git
cd teldrive

# Start development environment
docker-compose -f docker-compose.dev.yml up
```

---

## 🔄 **Docker Compose Variants**

### **Comparison Overview**
| Variant | Purpose | Build Time | Internet | Your Code | Production | UI Source |
|---------|---------|------------|----------|-----------|------------|-----------|
| **Local** | Custom builds | 3-5 min | Yes | ✅ Backend + UI | ✅ | Your repo |
| **Registry** | Official images | 30 sec | Yes | ❌ | ✅ | Official |
| **Custom UI** | UI customization | 4-6 min | Yes | ✅ UI only | ✅ | Your repo |
| **Dev** | Development | 2-4 min | Yes | ✅ Backend + UI | ❌ | Your repo |
| **Offline** | Air-gapped | 1-2 min | No | ✅ | ✅ | Cached |

---

### **1. 🏗️ docker-compose.local.yml** (Local Build)

**Purpose**: Build everything from source with your custom UI

```bash
docker-compose -f docker-compose.local.yml up
```

**What it builds:**
- ✅ TelDrive backend from your local source code
- ✅ UI from your GitHub repository (`https://github.com/cdgmx/teldrive-ui.git`)
- ✅ Uses all your local modifications

**Configuration:**
```yaml
build:
  args:
    - UI_SOURCE_TYPE=build
    - UI_REPO_URL=https://github.com/cdgmx/teldrive-ui.git
    - UI_REPO_BRANCH=main
    - BUILD_MODE=production
```

**When to use:**
- You've modified TelDrive source code
- You want your custom UI
- Self-hosting with full customization
- Testing local changes

---

### **2. 📦 docker-compose.registry.yml** (Registry/Pre-built)

**Purpose**: Use official TelDrive images from GitHub Container Registry

```bash
docker-compose -f docker-compose.registry.yml up
```

**What it uses:**
- ✅ Pre-built TelDrive image (`ghcr.io/tgdrive/teldrive:latest`)
- ✅ Official UI from TelDrive team
- ✅ No compilation needed

**Configuration:**
```yaml
image: ghcr.io/tgdrive/teldrive:latest
# No build process - just pulls image
```

**When to use:**
- Quick setup without customization
- Production with official releases
- Stable, tested versions
- Fastest deployment

---

### **3. 🎨 docker-compose.custom-ui.yml** (Custom UI Build)

**Purpose**: Official TelDrive backend with your custom UI

```bash
docker-compose -f docker-compose.custom-ui.yml up
```

**What it combines:**
- ✅ Official TelDrive backend (stable)
- ✅ Your custom UI from GitHub
- ✅ Best of both worlds

**Configuration:**
```yaml
build:
  args:
    - UI_SOURCE_TYPE=build
    - UI_REPO_URL=https://github.com/cdgmx/teldrive-ui.git
    - BUILD_MODE=production
```

**When to use:**
- UI customization only
- Keep backend stable
- Production with custom branding
- Showcase UI modifications

---

### **4. 🛠️ docker-compose.dev.yml** (Development)

**Purpose**: Full development environment with debugging tools

```bash
docker-compose -f docker-compose.dev.yml up
```

**What it includes:**
- ✅ Development build with debugging symbols
- ✅ Live code reloading (source mounted as volume)
- ✅ Debugging support (port 2345)
- ✅ Database access (port 5432)
- ✅ Redis cache (port 6379)
- ✅ Adminer DB management (port 8081)

**Configuration:**
```yaml
build:
  target: builder  # Development build
  args:
    - UI_SOURCE_TYPE=build
    - BUILD_MODE=development
volumes:
  - '.:/app:rw'  # Live code mounting
ports:
  - "2345:2345"  # Debugger
  - "5432:5432"  # PostgreSQL
  - "6379:6379"  # Redis
  - "8081:8080"  # Adminer
```

**When to use:**
- Active development
- Debugging backend issues
- UI development and testing
- Database inspection

---

## 🎯 **UI Build Options**

### **UI Build Methods**

#### **1. Build from Source** (`UI_SOURCE_TYPE=build`)
Builds UI directly from your GitHub repository during Docker build.

```bash
docker build --build-arg UI_SOURCE_TYPE=build \
             --build-arg UI_REPO_URL=https://github.com/cdgmx/teldrive-ui.git \
             --build-arg UI_REPO_BRANCH=main \
             -t teldrive:custom-ui .
```

**Process:**
1. Installs Node.js and pnpm in Docker
2. Clones your UI repository
3. Runs `pnpm install`
4. Runs `pnpm run build`
5. Copies `dist/*` to TelDrive

**Advantages:**
- ✅ Always uses latest code from your repository
- ✅ No dependency on GitHub releases
- ✅ Full control over UI build process
- ✅ Supports private repositories

#### **2. Download Pre-built** (`UI_SOURCE_TYPE=download`)
Downloads pre-built UI from GitHub releases (fallback method).

```bash
docker build --build-arg UI_SOURCE_TYPE=download \
             --build-arg UI_ASSET_URL=https://github.com/tgdrive/teldrive-ui/releases/download/latest/teldrive-ui.zip \
             -t teldrive:downloaded-ui .
```

**When to use:**
- Fallback when build fails
- Using official UI releases
- Faster builds (no compilation)

#### **3. Skip/Offline Mode** (`UI_SOURCE_TYPE=skip`)
Uses cached UI assets for offline builds.

```bash
# Pre-cache UI assets first
mkdir -p ui/dist
# Copy your built UI to ui/dist/

# Build offline
docker build --build-arg UI_SOURCE_TYPE=skip -t teldrive:offline .
```

**When to use:**
- Air-gapped environments
- Corporate networks
- Repeated builds with same UI
- CI/CD with cached assets

---

### **UI Build Arguments Reference**

| Argument | Options | Description | Example |
|----------|---------|-------------|---------|
| `UI_SOURCE_TYPE` | `build`, `download`, `skip` | How to obtain UI assets | `build` |
| `UI_REPO_URL` | GitHub URL | Your UI repository URL | `https://github.com/cdgmx/teldrive-ui.git` |
| `UI_REPO_BRANCH` | Branch name | Branch to build from | `main`, `develop` |
| `UI_ASSET_URL` | Download URL | Fallback download URL | GitHub releases URL |
| `GITHUB_TOKEN` | Token | For private repos or rate limiting | `ghp_xxxx` |
| `BUILD_MODE` | `production`, `development` | Build optimization level | `production` |

---

## ⚙️ **Configuration**

### **Required Files**

#### **1. storage.db**
```bash
# Create empty file for Telegram session storage
touch storage.db
```

#### **2. config.toml**
```bash
# Copy sample and customize
cp config.sample.toml config.toml
```

**Essential configuration:**
```toml
[db]
data-source = "postgres://teldriveadmin:password@db:5432/teldrive"

[jwt]
secret = "your-jwt-secret-change-this-in-production"
allowed-users = ["YourTelegramUsername"]

[tg]
storage-file = "/storage.db"

[tg.uploads]
encryption-key = "your-32-char-encryption-key-here"
```

### **Environment Variables**

#### **For Private UI Repositories:**
```bash
export GITHUB_TOKEN=your_github_token
docker-compose -f docker-compose.local.yml up
```

#### **For Custom UI Repository:**
```bash
# Override in docker-compose.yml or use build args
docker build --build-arg UI_REPO_URL=https://github.com/your-username/custom-ui.git
```

---

## 🛠️ **Development Workflow**

### **UI Development**
```bash
# 1. Start development environment
docker-compose -f docker-compose.dev.yml up

# 2. Make changes to your UI repository
# 3. Rebuild with latest changes
docker-compose -f docker-compose.dev.yml up --build

# 4. Access services:
# - TelDrive: http://localhost:8080
# - Database: localhost:5432 (adminer: http://localhost:8081)
# - Redis: localhost:6379
```

### **Backend Development**
```bash
# 1. Start with live code mounting
docker-compose -f docker-compose.dev.yml up

# 2. Code changes are automatically reflected
# 3. For debugging, attach to port 2345

# 4. Database management via Adminer:
# http://localhost:8081
# Server: db, User: teldriveadmin, Password: password
```

### **Testing Local Changes**
```bash
# Test backend changes
docker-compose -f docker-compose.local.yml up --build

# Test UI changes only
docker-compose -f docker-compose.custom-ui.yml up --build

# Test with cached UI (faster)
docker-compose -f docker-compose.local.yml build --build-arg UI_SOURCE_TYPE=skip
```

---

## 🚀 **Production Deployment**

### **Option 1: Full Custom Build**
```bash
# Build with your modifications
docker build --build-arg UI_SOURCE_TYPE=build \
             --build-arg UI_REPO_URL=https://github.com/cdgmx/teldrive-ui.git \
             -t teldrive:production .

# Deploy
docker run -d \
  --name teldrive \
  -p 8080:8080 \
  -v ./storage.db:/storage.db \
  -v ./config.toml:/config.toml \
  teldrive:production
```

### **Option 2: Docker Compose Production**
```bash
# Using custom UI compose
docker-compose -f docker-compose.custom-ui.yml up -d

# Or full local build
docker-compose -f docker-compose.local.yml up -d

# Or official images
docker-compose -f docker-compose.registry.yml up -d
```

### **Option 3: Air-gapped Deployment**
```bash
# 1. Build UI assets (with internet)
git clone https://github.com/cdgmx/teldrive-ui.git
cd teldrive-ui
pnpm install && pnpm run build
cp -r dist/* ../teldrive/ui/dist/

# 2. Build TelDrive (no internet needed)
cd ../teldrive
docker build --build-arg UI_SOURCE_TYPE=skip -t teldrive:offline .

# 3. Deploy
docker run -d \
  --name teldrive \
  -p 8080:8080 \
  -v ./storage.db:/storage.db \
  -v ./config.toml:/config.toml \
  teldrive:offline
```

---

## 🔒 **Security Considerations**

### **Private Repository Access**
```bash
# For private UI repositories
export GITHUB_TOKEN=your_personal_access_token

docker build --build-arg GITHUB_TOKEN=$GITHUB_TOKEN \
             --build-arg UI_SOURCE_TYPE=build \
             --build-arg UI_REPO_URL=https://github.com/private/ui-repo.git \
             -t teldrive:private .
```

### **Production Security**
```toml
# config.toml - Use strong secrets
[jwt]
secret = "use-a-strong-random-secret-here"

[tg.uploads]
encryption-key = "32-character-encryption-key-here"
```

### **File Permissions**
```bash
# Ensure proper permissions
chmod 600 config.toml  # Protect config file
chmod 666 storage.db   # Allow container access
```

---

## 🐛 **Troubleshooting**

### **Build Issues**

#### **UI Build Failures**
```bash
# Check if your UI repo builds locally
git clone https://github.com/cdgmx/teldrive-ui.git
cd teldrive-ui
pnpm install
pnpm run build  # Should create dist/ folder

# Fallback to download mode
docker build --build-arg UI_SOURCE_TYPE=download \
             --build-arg UI_ASSET_URL=https://github.com/tgdrive/teldrive-ui/releases/download/latest/teldrive-ui.zip
```

#### **Docker Build Errors**
```bash
# Clean build (remove cached layers)
docker-compose -f docker-compose.local.yml build --no-cache

# Check detailed build logs
docker build --progress=plain --no-cache .

# Check available space
docker system df
docker system prune  # Clean up if needed
```

### **Runtime Issues**

#### **Database Connection Errors**
```bash
# Check database status
docker-compose -f docker-compose.local.yml ps
docker-compose -f docker-compose.local.yml logs db

# Verify database is healthy
docker-compose -f docker-compose.local.yml exec db pg_isready -U teldriveadmin
```

#### **Storage File Issues**
```bash
# Check storage.db
ls -la storage.db
file storage.db  # Should show "empty" or "data"

# Recreate if it's a directory
rm -rf storage.db && touch storage.db

# Fix permissions
chmod 666 storage.db
```

#### **UI Not Loading**
```bash
# Check if UI was built/copied correctly
docker-compose -f docker-compose.local.yml exec teldrive ls -la /app/ui/

# Check for UI assets
docker-compose -f docker-compose.local.yml exec teldrive find /app -name "index.html"

# Rebuild with verbose output
docker-compose -f docker-compose.local.yml build --progress=plain
```

### **Network Issues**

#### **GitHub Access Problems**
```bash
# Test GitHub connectivity
curl -I https://github.com/cdgmx/teldrive-ui.git

# Use GitHub token for rate limiting
export GITHUB_TOKEN=your_token
docker build --build-arg GITHUB_TOKEN=$GITHUB_TOKEN

# Switch to offline mode
docker build --build-arg UI_SOURCE_TYPE=skip
```

#### **Port Conflicts**
```bash
# Check what's using port 8080
lsof -i :8080
netstat -tulpn | grep 8080

# Use different port
docker run -p 8081:8080 teldrive:local
```

### **Performance Issues**

#### **Slow Builds**
```bash
# Use cached UI for faster rebuilds
docker build --build-arg UI_SOURCE_TYPE=skip

# Enable BuildKit for better caching
export DOCKER_BUILDKIT=1
docker build .

# Use multi-stage build caching
docker build --target builder  # Stop at build stage
```

#### **High Memory Usage**
```bash
# Limit build memory
docker build --memory=2g .

# Use smaller base images
# (Already using alpine in Dockerfile)

# Clean up build cache
docker builder prune
```

---

## 📊 **Performance Comparison**

### **Build Times**
| Variant | Cold Build | Cached Build | Download Size | Final Image |
|---------|------------|--------------|---------------|-------------|
| Registry | 30s | 10s | ~60MB | ~60MB |
| Local (download UI) | 3-4min | 1-2min | ~10MB UI | ~80MB |
| Local (build UI) | 4-6min | 2-3min | ~5MB deps | ~85MB |
| Custom UI | 4-6min | 2-3min | ~5MB deps | ~85MB |
| Dev | 3-5min | 1-2min | ~5MB deps | ~120MB |
| Offline | 1-2min | 30s | 0MB | ~75MB |

### **Resource Usage**
| Variant | CPU (Build) | Memory (Build) | CPU (Runtime) | Memory (Runtime) |
|---------|-------------|----------------|---------------|------------------|
| Registry | Low | Low | Low | ~100MB |
| Local | High | Medium | Low | ~120MB |
| Custom UI | High | Medium | Low | ~120MB |
| Dev | Medium | High | Medium | ~200MB |
| Offline | Medium | Low | Low | ~100MB |

---

## 🎯 **Recommendations**

### **Choose Registry if:**
- ✅ Quick setup needed
- ✅ No customization required
- ✅ Production stability priority
- ✅ Official UI is sufficient

### **Choose Local if:**
- ✅ Backend modifications made
- ✅ Custom UI required
- ✅ Full control needed
- ✅ Self-hosting with customization

### **Choose Custom UI if:**
- ✅ UI customization only
- ✅ Stable backend preferred
- ✅ Custom branding needed
- ✅ Production with UI changes

### **Choose Dev if:**
- ✅ Active development
- ✅ Debugging required
- ✅ Database management needed
- ✅ UI/backend development

### **Choose Offline if:**
- ✅ Corporate/restricted networks
- ✅ Air-gapped environments
- ✅ Repeated builds
- ✅ CI/CD with caching

---

## 📝 **Your UI Repository Requirements**

### **Required Structure**
```
your-teldrive-ui/
├── package.json          # Must have "build" script
├── pnpm-lock.yaml        # pnpm lockfile
├── src/                  # Source files
├── public/               # Static assets
└── dist/                 # Build output (created by pnpm run build)
```

### **Required Scripts in package.json**
```json
{
  "scripts": {
    "build": "vite build",  // or your build command
    "dev": "vite dev",      // for development
    "preview": "vite preview"
  }
}
```

### **Build Process**
The Docker build will:
1. Clone your repository
2. Run `pnpm install`
3. Run `pnpm run build`
4. Copy `dist/*` to TelDrive's UI location

---

## 🔗 **Useful Commands**

### **Quick Commands Reference**
```bash
# Build with your custom UI
docker-compose -f docker-compose.custom-ui.yml up

# Development environment
docker-compose -f docker-compose.dev.yml up

# Official images (fastest)
docker-compose -f docker-compose.registry.yml up

# Offline build
docker build --build-arg UI_SOURCE_TYPE=skip -t teldrive:offline .

# Clean rebuild
docker-compose -f docker-compose.local.yml build --no-cache

# View logs
docker-compose -f docker-compose.local.yml logs -f teldrive

# Shell access
docker-compose -f docker-compose.local.yml exec teldrive sh

# Database access
docker-compose -f docker-compose.local.yml exec db psql -U teldriveadmin -d teldrive
```

### **Maintenance Commands**
```bash
# Update images
docker-compose -f docker-compose.registry.yml pull

# Clean up
docker system prune
docker volume prune

# Backup database
docker-compose -f docker-compose.local.yml exec db pg_dump -U teldriveadmin teldrive > backup.sql

# Restore database
docker-compose -f docker-compose.local.yml exec -T db psql -U teldriveadmin teldrive < backup.sql
```

---

This guide covers all aspects of building and deploying TelDrive with various configurations. Choose the method that best fits your needs and environment!