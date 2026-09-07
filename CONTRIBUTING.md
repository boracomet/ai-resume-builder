# Contributing to AI Resume Builder

Thank you for your interest in contributing! This project welcomes bug reports, feature ideas, documentation improvements, and pull requests.

## Getting Started

1. Fork the repository and clone your fork locally.
2. Install [Go 1.22+](https://go.dev/dl/) and ensure Chrome/Chromium is available for PDF export.
3. Copy environment variables: `cp .env.example .env`
4. Run the server: `go run ./cmd/server`
5. Open [http://localhost:8080](http://localhost:8080)

## Development Workflow

1. Create a branch from `main`:

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and keep commits focused and descriptive.

3. Run tests before opening a pull request:

   ```bash
   go build ./...
   go test ./...
   ```

4. Push your branch and open a pull request against `main`.

## Pull Request Guidelines

- Describe **what** changed and **why**.
- Link related issues when applicable.
- Keep changes scoped — one logical change per PR when possible.
- Ensure `go build ./...` and `go test ./...` pass locally (CI runs the same checks).

## Code Style

- Follow existing patterns in the codebase (naming, file layout, handler structure).
- Prefer small, readable diffs over large refactors unless discussed first.
- Add tests for non-trivial behavior when it makes sense.

## Reporting Issues

When filing a bug report, include:

- Steps to reproduce
- Expected vs. actual behavior
- Go version, OS, and whether you use Docker or local dev
- Relevant logs or screenshots

## Questions

Open a [GitHub Discussion](https://github.com/boracomet/ai-resume-builder/discussions) or issue if you are unsure about an approach before investing significant time.
