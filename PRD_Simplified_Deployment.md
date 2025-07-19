# PRD: Simplified TelDrive Deployment Strategy

## STEP 1 – Generate the Plan

### Repo Snapshot
• **Current Status**: Multiple Docker Compose variants creating complexity
• **Goal**: Simplify to two clear deployment modes: Production and Development
• **Approach**: Remove registry dependency, offline mode, and backward compatibility
• **Focus**: Custom implementation as the primary production method

### Story Context
TelDrive currently has multiple deployment variants (local, registry, custom-ui, dev, offline) which creates confusion and maintenance overhead. We want to simplify this to just two clear modes: a Production deployment that builds everything from source with custom UI, and a Development environment for active development work.

### Legend
✅ Completed 🔄 In Progress ⏳ Not Started

### Tasks Checklist

- [ ] ⏳ **1 – Consolidate to Production Docker Compose**
  - **Description**: Create single production compose that builds TelDrive + custom UI from source
  - **Priority**: High
  - **Prerequisites**: Current docker-compose variants understanding
  - **Blockers**: None
  - **Sub-tasks**:
    - [ ] Create `docker-compose.yml` as the main production deployment
    - [ ] Remove `docker-compose.local.yml`, `docker-compose.registry.yml`, `docker-compose.custom-ui.yml`
    - [ ] Configure production build with custom UI as default
    - [ ] Set optimal production build arguments and settings
    - [ ] Test production deployment functionality

- [ ] ⏳ **2 – Streamline Development Environment**
  - **Description**: Simplify development compose to focus on development workflow
  - **Priority**: High
  - **Prerequisites**: Understanding of development needs
  - **Blockers**: None
  - **Sub-tasks**:
    - [ ] Keep `docker-compose.dev.yml` as the only development variant
    - [ ] Optimize for development workflow (live reload, debugging)
    - [ ] Remove unnecessary development complexity
    - [ ] Ensure fast rebuild cycles for development
    - [ ] Configure development-specific services (Redis, Adminer)

- [ ] ⏳ **3 – Remove Offline and Registry Support**
  - **Description**: Clean up Dockerfile and scripts to remove unused build modes
  - **Priority**: Medium
  - **Prerequisites**: Completed Tasks 1-2
  - **Blockers**: 1, 2
  - **Sub-tasks**:
    - [ ] Remove `UI_SOURCE_TYPE=download` and `UI_SOURCE_TYPE=skip` from Dockerfile
    - [ ] Simplify UI build process to only support building from source
    - [ ] Remove registry image references and fallback mechanisms
    - [ ] Clean up build arguments to only essential ones
    - [ ] Remove offline mode documentation and examples

- [ ] ⏳ **4 – Update Documentation for Simplified Approach**
  - **Description**: Rewrite documentation to reflect the simplified two-mode approach
  - **Priority**: Medium
  - **Prerequisites**: Completed Tasks 1-3
  - **Blockers**: 3
  - **Sub-tasks**:
    - [ ] Update `TELDRIVE_BUILD_GUIDE.md` to show only Production and Development
    - [ ] Create clear "Production vs Development" comparison
    - [ ] Remove references to registry, offline, and multiple variants
    - [ ] Add simplified quick start guide
    - [ ] Update troubleshooting for the two-mode approach

- [ ] ⏳ **5 – Optimize Build Performance**
  - **Description**: Focus on making the two remaining modes as fast and efficient as possible
  - **Priority**: Medium
  - **Prerequisites**: Completed Tasks 1-3
  - **Blockers**: 3
  - **Sub-tasks**:
    - [ ] Optimize Dockerfile layer caching for production builds
    - [ ] Implement efficient development rebuild strategy
    - [ ] Add BuildKit optimizations for both modes
    - [ ] Configure optimal resource usage for each mode
    - [ ] Test and benchmark build performance improvements

- [ ] ⏳ **6 – Create Migration Guide**
  - **Description**: Help users transition from old multi-variant approach to simplified approach
  - **Priority**: Low
  - **Prerequisites**: Completed Tasks 1-4
  - **Blockers**: 4
  - **Sub-tasks**:
    - [ ] Document migration from `docker-compose.local.yml` to `docker-compose.yml`
    - [ ] Explain changes from registry-based to source-based deployment
    - [ ] Provide command equivalents for old deployment methods
    - [ ] Create troubleshooting for migration issues
    - [ ] Add FAQ for the simplified approach

### Success Criteria
**Task 1** – Single `docker-compose.yml` successfully builds and runs TelDrive with custom UI; old compose files removed; production deployment works identically to previous custom-ui variant.  
**Task 2** – `docker-compose.dev.yml` provides efficient development environment; live reload functional; debugging capabilities maintained; development services properly configured.  
**Task 3** – Dockerfile only supports building UI from source; no registry or offline fallbacks; build process simplified and streamlined; unused build arguments removed.  
**Task 4** – Documentation clearly explains Production vs Development modes only; quick start guide updated; no references to removed variants; troubleshooting covers simplified approach.  
**Task 5** – Production builds complete in <4 minutes; development rebuilds in <2 minutes; BuildKit optimizations implemented; resource usage optimized for each mode.  
**Task 6** – Migration guide enables smooth transition from old approach; command equivalents provided; common migration issues documented and resolved.

## STEP 2 – Critique the Plan

### Simplified Approach Benefits
• **Reduced Complexity**: Two clear modes instead of five variants eliminates confusion
• **Maintenance Efficiency**: Less code to maintain, fewer edge cases to handle
• **User Experience**: Clear choice between "Production" and "Development"
• **Performance Focus**: Optimization efforts concentrated on two modes instead of five

### Potential Risks and Mitigations
• **Breaking Changes**: Removing registry support may affect existing users
  - *Mitigation*: Provide clear migration guide and command equivalents
• **Reduced Flexibility**: Some users may prefer registry-based deployment
  - *Mitigation*: Source-based builds provide more control and customization
• **Build Time Increase**: All deployments now require building from source
  - *Mitigation*: Focus on build optimization and caching strategies

### Task Dependencies and Logic
• **Sequential Approach**: Tasks 1-2 can run in parallel, then 3 depends on both
• **Documentation Last**: Task 4 waits for implementation completion
• **Performance Focus**: Task 5 optimizes the simplified implementation
• **User Support**: Task 6 helps with transition

### Success Criteria Validation
• **Measurable Targets**: Build time targets specified for both modes
• **Functional Validation**: Each task has clear deliverables and verification steps
• **User Impact**: Migration guide ensures existing users can transition smoothly

### Implementation Strategy
**Phase 1 (Week 1)**: Tasks 1, 2 - Core simplification
- Create production docker-compose.yml
- Streamline development environment

**Phase 2 (Week 2)**: Task 3 - Cleanup
- Remove unused build modes and complexity

**Phase 3 (Week 3)**: Tasks 4, 5, 6 - Documentation and optimization
- Update documentation
- Optimize performance
- Create migration guide

### Critical Success Factors
1. **Clear Mode Distinction**: Production and Development must have obvious differences and use cases
2. **Performance Maintenance**: Simplified approach should not sacrifice build performance
3. **User Migration**: Existing users must be able to transition without data loss
4. **Documentation Quality**: Clear guidance on when to use each mode

The simplified approach reduces complexity while maintaining all essential functionality, focusing on the two most important use cases: production deployment and active development.