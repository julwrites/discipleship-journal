# Task Documentation System

This project uses a strict **Task Documentation System** to track all work.

## Core Philosophy
**"If it's not documented in `docs/tasks/`, it didn't happen."**

## Quick Start
We use the `scripts/tasks` wrapper to manage tasks.

```bash
# List all tasks
./scripts/tasks list

# Find the next task to work on
./scripts/tasks next

# Create a new task
./scripts/tasks create features "Title of my feature"

# Update task status
./scripts/tasks update [TASK_ID] in_progress
```

## Documentation
For the full guide on how to use this system, including file formats and workflows, please read:
👉 [**GUIDE.md**](./GUIDE.md)

## Directory Structure
*   `foundation/`: Core architecture and setup
*   `infrastructure/`: Services, adapters, platform code
*   `domain/`: Business logic, use cases
*   `presentation/`: UI, state management
*   `features/`: End-to-end feature implementation
*   `migration/`: Refactoring, upgrades
*   `testing/`: Testing infrastructure
*   `review/`: Code reviews and PR analysis
*   `security/`: Security reviews and tasks
