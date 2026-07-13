<!--
### Sync Impact Report
- Version change: [CONSTITUTION_VERSION] (template) -> 1.0.0
- List of modified principles:
  - PRINCIPLE_1: [PRINCIPLE_1_NAME] -> I. Code Quality & Standards
  - PRINCIPLE_2: [PRINCIPLE_2_NAME] -> II. Focused Microservice Design
  - PRINCIPLE_3: [PRINCIPLE_3_NAME] -> III. Test-Driven Development (TDD)
  - PRINCIPLE_4: [PRINCIPLE_4_NAME] -> IV. Integration & Contract Testing
  - PRINCIPLE_5: [PRINCIPLE_5_NAME] -> V. Observability & Simplicity
- Added sections:
  - Constraints & Technology Stack
  - Development & Review Workflow
- Removed sections: None
- Templates requiring updates:
  - .specify/templates/plan-template.md (✅ updated)
  - .specify/templates/spec-template.md (✅ updated)
  - .specify/templates/tasks-template.md (✅ updated)
- Follow-up TODOs: None
-->
# Proxmox DNS Sync Constitution

## Core Principles

### I. Code Quality & Standards
Code MUST be clean, well-formatted, and subject to automatic static analysis and linting.
All code MUST exhibit strict type safety and aim for high test coverage (>90%). Comments and
docstrings MUST be preserved and maintained.
Rationale: Code quality ensures long-term maintainability, readability, and prevents technical
debt in a synchronization service where correctness is critical.

### II. Focused Microservice Design
The project MUST be structured as a set of highly focused, single-responsibility microservices
or modules. Each module MUST have a clear boundary, minimal external dependencies, and expose
well-defined APIs or CLI interfaces. Inter-service coupling MUST be minimized.
Rationale: Keeping components isolated and small makes them easier to understand, test, modify,
and deploy independently, especially when interfacing with external Proxmox APIs.

### III. Test-Driven Development (TDD) (NON-NEGOTIABLE)
Strict Test-Driven Development (TDD) workflow MUST be followed. Implementation tasks MUST
start with defining unit and integration tests, demonstrating their failure, writing the
minimal code required to pass, and then refactoring (Red-Green-Refactor).
Rationale: Prevents regression, ensures robust test suites, and guarantees that code design is
driven by consumption and validation requirements.

### IV. Integration & Contract Testing
Integration and contract tests MUST verify boundary conditions, API schemas, and communication
protocols (e.g., DNS providers, Proxmox APIs) under both success and failure states. External
service integrations MUST be fully simulated/mocked for unit testing and verified via
integration testing.
Rationale: Sync operations rely on external systems, so verifying contracts and integration
behaviors prevents silent failures.

### V. Observability & Simplicity
Every service MUST produce structured JSON log output with severity levels, contextual metadata,
and support tracing/metrics where necessary. Keep the code simple, avoiding pre-mature
optimization (YAGNI).
Rationale: Sync daemons run in the background; structured observability is vital to trace errors
and understand state synchronization without manual intervention.

## Constraints & Technology Stack

- **Language & Runtime**: To be defined in individual feature specifications, but must support
  static typing/linting.
- **Configuration**: All configuration MUST be loaded from environment variables or structured
  files, following the 12-Factor App methodology.
- **Error Handling**: Failures MUST be handled gracefully with clear, user-friendly error messages
  and stack traces logged at the appropriate level.

## Development & Review Workflow

- **Pre-commit Checks**: Run local linters, typecheckers, and test suites prior to pushing code.
- **Review Gates**: No code change shall be merged without passing all automated tests and
  static analysis.

## Governance

This constitution governs all design decisions and code implementations in the `proxmox-dns-sync`
repository. Amendments require documentation updates, a version bump, and team approval.

**Version**: 1.0.0 | **Ratified**: 2026-07-13 | **Last Amended**: 2026-07-13
