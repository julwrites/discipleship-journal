# AI Agent Instructions

You are an expert Software Engineer working on this project. Your primary responsibility is to implement features and fixes while strictly adhering to the **Task Documentation System**.

## Core Philosophy
This repository uses the agent-harness conventions.

## Documentation Reference
*   **Architecture**: Refer to `docs/architecture/` for system design.
*   **Features**: Refer to `docs/features/` for feature specifications.
*   **Security**: Refer to `docs/security/` for risk assessments and mitigations.

## Code Style & Standards

### General
*   **Testing Rigor**: Always update tests when modifying code. Ensure all tests pass before finishing a task.
*   **Infrastructure Maintenance**: Keep local deployment scripts (e.g., `docker-compose.yml`) synchronized with production configuration changes.
*   Follow the existing patterns in the codebase.
*   Ensure all new code is covered by tests (if testing infrastructure exists).
*   **Commit Messages**: Use descriptive commit messages. Format: `type(scope): description`.

### Frontend (React/Vite)
*   Use TypeScript for all new code.
*   Use Functional Components and Hooks.
*   Use Tailwind CSS for styling.
*   Use Shadcn/UI for UI components.
*   Follow standard directory structure: `components`, `pages`, `hooks`, `services`, `types`.

### Backend (Go)
*   Follow standard Go idioms (Effective Go).
*   Use `chi` for routing.
*   Use `pgx` for PostgreSQL interactions.
*   Keep handler logic separate from business logic (services).
*   Use environment variables for configuration.

## PR Review Methodology
When performing a PR review, follow this "Human-in-the-loop" process to ensure depth and efficiency.

### 1. Preparation
1.  **Create Task**: `python3 scripts/tasks.py create review "Review PR #<N>: <Title>"`
2.  **Fetch Details**: Use `gh` to get the PR context.
    *   `gh pr view <N>`
    *   `gh pr diff <N>`

### 2. Analysis & Planning (The "Review Plan")
**Do not review line-by-line yet.** Instead, analyze the changes and document a **Review Plan** in the task file (or present it for approval).

Your plan must include:
*   **High-Level Summary**: Purpose, new APIs, breaking changes.
*   **Dependency Check**: New libraries, maintenance status, security.
*   **Impact Assessment**: Effect on existing code/docs.
*   **Focus Areas**: Prioritized list of files/modules to check.
*   **Suggested Comments**: Draft comments for specific lines.
    *   Format: `File: <path> | Line: <N> | Comment: <suggestion>`
    *   Tone: Friendly, suggestion-based ("Consider...", "Nit: ...").

### 3. Execution
Once the human approves the plan and comments:
1.  **Pending Review**: Create a pending review using `gh`.
    *   `COMMIT_SHA=$(gh pr view <N> --json headRefOid -q .headRefOid)`
    *   `gh api repos/{owner}/{repo}/pulls/{N}/reviews -f commit_id="$COMMIT_SHA"`
2.  **Batch Comments**: Add comments to the pending review.
    *   `gh api repos/{owner}/{repo}/pulls/{N}/comments -f body="..." -f path="..." -f commit_id="$COMMIT_SHA" -F line=<L> -f side="RIGHT"`
3.  **Submit**:
    *   `gh pr review <N> --approve --body "Summary..."` (or `--request-changes`).

### 4. Close Task
*   Update task status to `completed`.

