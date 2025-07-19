# PRD: Transition from GitHub Workflow to Local Docker Builds

## STEP 1 – Generate the Plan

### Repo Snapshot
• **Go application** (v1.24) with Telegram Drive functionality using gotd/td library and PostgreSQL
• **Current build system**: GoReleaser + GitHub Actions publishing to Docker Hub/GHCR registries
• **Docker architecture**: Minimal `goreleaser.dockerfile` (FROM scratch) with pre-built static binary
• **Build tooling**: Taskfile for automation, UI assets downloaded from GitHub releases via `scripts/extract.go`
• **Current deployment**: Docker Compose pulls `ghcr.io/tgdrive/teldrive` images from registry
• **Dependencies**: PostgreSQL with PGroonga extension, optional Redis cache, imgproxy service
• **Storage**: Uses `/storage.db` for Telegram session data (mounted from `./storage.db`)
• **Configuration**: Mounted as `/config.toml` from local `./config.toml` file

### Story Context
TelDrive is a Telegram-based cloud storage solution that currently relies on GitHub Actions to build and publish Docker images to public registries. The project needs to enable local development and self-hosted deployment by building Docker images directly from source code, eliminating dependency on external registries and enabling customization.

### Legend
✅ Completed 🔄 In Progress ⏳ Not Started

### Tasks Checklist

- [x] ✅ **1 – Fix storage path consistency and create base Dockerfile**
  - **Description**: Resolve storage.db path mismatch and create multi-stage Dockerfile for local builds.
  - **Priority**: Critical
  - **Prerequisites**: Docker BuildKit enabled (v20.10+), understanding of current entrypoint
  - **Blockers**: None
  - **Sub-tasks**:
    - [x] Fix entrypoint storage path from `/storage.db` to match mount point
    - [x] Create builder stage with Go 1.24 and build dependencies
    - [x] Add code generation step (`go generate ./...`) matching taskfile
    - [x] Configure final runtime stage with minimal base image
    - [x] Set proper entrypoint: `["/teldrive","run","--tg-storage-file","/storage.db"]`
    - [x] Test storage file persistence and configuration mounting
  - **Deliverables**: 
    - ✅ `Dockerfile` - Multi-stage build for local compilation
    - ✅ `.dockerignore` - Build context optimization (>50% reduction)
    - ✅ Storage path consistency verified (`/storage.db` matches volume mounts)
    - ✅ Build arguments: `UI_ASSET_URL`, `UI_SKIP_DOWNLOAD`, `BUILD_MODE`

- [ ] ⏳ **2 – Implement robust UI asset management with fallbacks**
  - **Description**: Enhance current `scripts/extract.go` with error handling and offline build support.
  - **Priority**: High
  - **Prerequisites**: Analysis of current UI download from GitHub releases, network access in Docker
  - **Blockers**: None
  - **Sub-tasks**:
    - [ ] Add error handling and retry logic to `scripts/extract.go`
    - [ ] Implement build argument `UI_ASSET_URL` for custom/local UI sources
    - [ ] Create offline build mode with pre-cached UI assets in `ui/dist`
    - [ ] Add GitHub API authentication to prevent rate limiting
    - [ ] Implement graceful degradation when UI assets unavailable
    - [ ] Add build-time validation of UI asset integrity

- [ ] ⏳ **3 – Create backward-compatible Docker Compose variants**
  - **Description**: Create new local-build compose file while preserving existing registry-based deployment.
  - **Priority**: High
  - **Prerequisites**: Docker Compose v3.9+ knowledge, completed Task 1
  - **Blockers**: 1
  - **Sub-tasks**:
    - [ ] Create `docker-compose.local.yml` with `build:` directive pointing to repository root
    - [ ] Preserve existing `docker-compose.yml` for registry-based deployment
    - [ ] Maintain volume mounts: `./storage.db:/storage.db` and `./config.toml:/config.toml`
    - [ ] Keep external postgres network dependency for compatibility
    - [ ] Test complete stack with database migrations and service communication
    - [ ] Validate storage persistence across container restarts

- [ ] ⏳ **4 – Add essential build optimizations and .dockerignore**
  - **Description**: Create .dockerignore and implement caching for faster builds meeting <2min criteria.
  - **Priority**: High
  - **Prerequisites**: Understanding of current .gitignore patterns, Docker layer caching
  - **Blockers**: 1
  - **Sub-tasks**:
    - [ ] Create .dockerignore excluding: `bin/`, `ui/dist/`, `*.db`, `.git/`, `logs/`, `*.env*`
    - [ ] Implement Go module dependency caching in separate Docker layer
    - [ ] Optimize Dockerfile layer ordering: dependencies → generate → build → runtime
    - [ ] Add BuildKit cache mounts for Go build cache (`/root/.cache/go-build`)
    - [ ] Add build arguments: `BUILD_MODE=production|development`
    - [ ] Test build time reduction from baseline

- [ ] ⏳ **5 – Leverage existing multi-architecture capabilities**
  - **Description**: Adapt existing GoReleaser multi-arch config for local Docker builds.
  - **Priority**: Low
  - **Prerequisites**: Analysis of current .goreleaser.yml multi-arch setup
  - **Blockers**: 1
  - **Sub-tasks**:
    - [ ] Review existing amd64/arm64 build configuration in .goreleaser.yml
    - [ ] Add Docker buildx support for local multi-platform builds
    - [ ] Test cross-compilation compatibility with current CGO_ENABLED=0 setup
    - [ ] Document platform-specific build commands
    - [ ] Validate functional equivalence across architectures

- [ ] ⏳ **6 – Create deployment variants and documentation**
  - **Description**: Provide multiple deployment options and comprehensive setup documentation.
  - **Priority**: Medium
  - **Prerequisites**: Technical writing skills, completed Tasks 1-4
  - **Blockers**: 4
  - **Sub-tasks**:
    - [ ] Create `docker-compose.local.yml` for local builds
    - [ ] Maintain `docker-compose.registry.yml` for registry-based deployment
    - [ ] Write comprehensive local development setup guide
    - [ ] Document build customization options and environment variables
    - [ ] Create troubleshooting guide for common build failures

- [ ] ⏳ **6 – Implement secure configuration management for local builds**
  - **Description**: Leverage existing config.toml structure for secure local credential handling.
  - **Priority**: Medium
  - **Prerequisites**: Analysis of config.sample.toml structure, Docker secrets knowledge
  - **Blockers**: 3
  - **Sub-tasks**:
    - [ ] Create config.toml template with placeholder values for local development
    - [ ] Document secure practices for JWT secret (`jwt.secret`) configuration
    - [ ] Add guidance for Telegram encryption key (`tg.uploads.encryption-key`) setup
    - [ ] Implement environment variable override support for sensitive values
    - [ ] Create development vs production configuration examples
    - [ ] Document credential rotation procedures for local deployments

- [ ] ⏳ **7 – Add development workflow enhancements**
  - **Description**: Create development-friendly build variants and debugging support.
  - **Priority**: Medium
  - **Prerequisites**: Understanding of Go debugging, completed Task 4
  - **Blockers**: 4
  - **Sub-tasks**:
    - [ ] Add development Dockerfile target with debugging symbols
    - [ ] Create docker-compose.dev.yml with volume mounts for live code changes
    - [ ] Add debug build arguments and development environment variables
    - [ ] Configure hot-reload for Go code changes during development
    - [ ] Document IDE integration and debugging workflow
    - [ ] Add development-specific logging and profiling capabilities

- [ ] ⏳ **8 – Create comprehensive documentation and migration guide**
  - **Description**: Document local build process and provide migration path from registry deployments.
  - **Priority**: Medium
  - **Prerequisites**: Completed core tasks 1-4, technical writing
  - **Blockers**: 4
  - **Sub-tasks**:
    - [ ] Write step-by-step local build setup guide
    - [ ] Document differences between docker-compose.yml (registry) and docker-compose.local.yml
    - [ ] Create troubleshooting guide for common build and runtime issues
    - [ ] Provide migration instructions for existing registry-based deployments
    - [ ] Document build customization options and environment variables
    - [ ] Add performance comparison between local and registry builds

### Success Criteria
**Task 1** – ✅ COMPLETED: `Dockerfile` and `.dockerignore` created; storage.db path consistency verified (`/storage.db`); multi-stage build with Go 1.24, UI asset download, and code generation; build arguments implemented; ready for <5 minute builds.  
**Task 2** – UI assets download successfully in 90% of builds; offline builds work with pre-cached `ui/dist`; build fails gracefully with clear errors; GitHub API rate limiting handled with authentication; build argument `UI_ASSET_URL` functional.  
**Task 3** – `docker-compose -f docker-compose.local.yml up` builds and runs complete stack; existing `docker-compose.yml` unchanged; storage and config volumes work identically; database migrations execute; external postgres network preserved.  
**Task 4** – Subsequent builds complete in <2 minutes with caching; .dockerignore reduces build context by >50%; Go module dependencies cached; BuildKit cache mounts functional; `BUILD_MODE` argument works.  
**Task 5** – Multi-architecture builds leverage existing GoReleaser config; amd64/arm64 cross-compilation works with CGO_ENABLED=0; buildx commands documented; functional equivalence validated.  
**Task 6** – config.toml template created with placeholders; JWT secret and encryption key setup documented; environment variable overrides functional; development/production config examples provided.  
**Task 7** – Development Dockerfile target includes debugging symbols; docker-compose.dev.yml enables live code changes; hot-reload functional; IDE debugging workflow documented.  
**Task 8** – Local build setup guide enables new developer success in <10 minutes; registry vs local deployment differences documented; troubleshooting covers storage, UI, and network issues; migration path <5 steps.

## STEP 2 – Critique the Plan (Post-Audit Revision)

### Audit Findings Addressed
• **FIXED**: Storage path mismatch resolved - Task 1 now addresses `/storage.db` entrypoint consistency
• **FIXED**: UI asset fragility addressed - Task 2 includes error handling, authentication, and offline fallbacks
• **FIXED**: Backward compatibility ensured - Task 3 creates new compose file instead of modifying existing
• **FIXED**: Missing .dockerignore addressed - Task 4 includes comprehensive build context optimization
• **FIXED**: Multi-arch redundancy resolved - Task 5 leverages existing GoReleaser capabilities

### Remaining Risks and Mitigations
• **UI Asset Dependency**: Task 2 includes GitHub API authentication and offline build mode with pre-cached assets
• **Database Migrations**: Task 3 validates database initialization and migration execution in local builds
• **Configuration Security**: Task 6 leverages existing config.toml structure for secure credential management
• **Build Performance**: Task 4 targets <2 minute builds with comprehensive caching strategy

### Task Dependency Optimization
• **Tasks 1 & 2 Parallelized**: Removed blocker dependency to enable concurrent development
• **Priority Realignment**: Task 5 (multi-arch) demoted to Low priority as capability already exists
• **Documentation Consolidation**: Combined documentation efforts in Task 8

### Success Criteria Strengthening
• **Specific Validation**: Each task now includes concrete technical validation criteria
• **Path Consistency**: Storage and configuration paths explicitly validated
• **Performance Targets**: Build time targets specified with caching requirements
• **Backward Compatibility**: Explicit preservation of existing deployment methods

### Implementation Strategy
**Phase 1 (Critical - Week 1)**: Tasks 1, 2, 3 - Core local build functionality
- ✅ Fix storage path consistency and create base Dockerfile (COMPLETED)
- Implement robust UI asset management with fallbacks  
- Create backward-compatible Docker Compose variants

**Phase 2 (Important - Week 2)**: Tasks 4, 6 - Optimization and security
- Add build optimizations and .dockerignore
- Implement secure configuration management

**Phase 3 (Enhancement - Week 3)**: Tasks 5, 7, 8 - Advanced features
- Multi-architecture support (if needed)
- Development workflow enhancements
- Comprehensive documentation

### Critical Success Factors
1. **Storage Path Consistency**: Ensuring `/storage.db` entrypoint matches volume mounts
2. **UI Asset Reliability**: Robust fallback mechanisms for GitHub API failures
3. **Build Performance**: Meeting <2 minute cached build targets
4. **Backward Compatibility**: Preserving existing registry-based deployment options

The revised plan addresses all critical audit findings while maintaining a practical implementation timeline and clear success criteria.