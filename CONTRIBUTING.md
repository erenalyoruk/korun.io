# THIS IS NOT COMPLETE!

# 🤝 Contributing to Korun.io

We're thrilled you're interested in contributing to Korun.io! Your efforts, no matter how small, help us build a more secure and developer-friendly secret management platform.

This document outlines the guidelines and best practices for contributing to Korun.io.

## Table of Contents

- [THIS IS NOT COMPLETE!](#this-is-not-complete)
- [🤝 Contributing to Korun.io](#-contributing-to-korunio)
  - [Table of Contents](#table-of-contents)
  - [Code of Conduct](#code-of-conduct)
  - [How to Contribute](#how-to-contribute)
    - [Reporting Bugs](#reporting-bugs)
    - [Suggesting Features](#suggesting-features)
    - [Writing Code](#writing-code)
    - [Improving Documentation](#improving-documentation)
  - [Getting Started: Setting up Your Development Environment](#getting-started-setting-up-your-development-environment)
    - [Prerequisites](#prerequisites)
    - [Forking the Repository](#forking-the-repository)
    - [Cloning Your Fork](#cloning-your-fork)
  - [Contribution Guidelines](#contribution-guidelines)
    - [Branching Strategy](#branching-strategy)
    - [Commit Messages](#commit-messages)
    - [Pull Request (PR) Process](#pull-request-pr-process)
    - [Code Style \& Quality](#code-style--quality)
    - [Testing](#testing)
  - [Project Structure Overview](#project-structure-overview)
  - [Need Help?](#need-help)

## <a name="code-of-conduct"></a>Code of Conduct

Please note that this project is released with a [Contribute Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating in this project, you agree to abide by its term.

## <a name="how-to-contribute"></a>How to Contribute

There are many ways to contribute, not just by writing code:

### <a name="reporting-bugs"></a>Reporting Bugs

If you find a bug, please open an issue on our [GitHub Issues page](https://github.com/erenalyoruk/korun.io/issues) and provide:

- A clear and concise description of the bug.
- Steps to reproduce the behavior.
- Expected behavior.
- Actual behavior.
- Screenshots or error messages (if applicable).
- Your environment details (OS, browser, Go version, Node.js version, etc.).

### <a name="suggesting-features"></a>Suggesting Features

We love new ideas! If you have a feature request or a suggestion for improvement, please open an issue on our [GitHub Issues page](https://github.com/erenalyoruk/korun.io/issues) and include:

- A clear and concise description of the feature.
- The problem it solves.
- Any potential solutions or design ideas you have.

### <a name="writing-code"></a>Writing Code

This is where the magic happens! Whether it's fixing a bug, implementing a new feature, or refactoring existing code, all code contributions are highly valued. See the [Contribution Guidelines](#4-contribution-guidelines) section for details on how to submit your code.

### <a name="improving-documentation"></a>Improving Documentation

Clear and comprehensive documentation is key for any successful project. If you find errors, omissions, or areas for improvement in our README.md, CONTRIBUTING.md, or any inline code comments, please open a PR.

## <a name="getting-started-setting-up-your-development-environment"></a>Getting Started: Setting up Your Development Environment

### <a name="prerequisites"></a>Prerequisites

TODO: prerequisities

### <a name="forking-the-repository"></a>Forking the Repository

1. Go to the [Korun.io GitHub repository](https://github.com/erenalyoruk/korun.io).
2. Click the "Fork" button in the top right corner.

### <a name="cloning-your-fork"></a>Cloning Your Fork

```bash
git clone https://github.com/YOUR_USERNAME/korun.io.git
cd korun.io
```

## <a name="contribution-guidelines"></a>Contribution Guidelines

### <a name="branching-strategy"></a>Branching Strategy

We use a feature-branch workflow.

1. **Create a new branch** from `main`:
   ```bash
   git checkout main
   git pull origin main
   git checkout -b feature/my-awesome-feature # or bugfix/fix-issue-123
   ```
2. **Work on your feature/bugfix**.
3. **Push your branch** to your fork.

### <a name="commit-messages"></a>Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) specification for clear and concise commit history.

Examples:

- `feat: add user registration endpoint`
- `feat!: send an email to the customer when a product is shipped`
  - Has `!` at the end to draw attention to breaking change.
- `fix(auth-service): correct password hashing bug`
- `docs: update contributing guide`
- `refactor(secrets-service): extract encryption logic`
- `chore: update go dependencies`

### <a name="pull-request-pr-process"></a>Pull Request (PR) Process

- **Submit a Pull Request** to the `main` branch of the `korun.io` repository.
- **Provide a clear title** and detailed description of your changes.
	- Link to any relevant issues (e.g., `Closes #123`).
	- Explain what problem your PR solves and how.
	- Include screenshots or GIFs for UI changes.
- **Ensure your PR passes all CI checks**.
- **Be responsive to feedback** during the code review process. We might suggest changes to improve code quality, performance, or adhere to project standards.

### <a name="code-style--quality"></a>Code Style & Quality

- **Go**:
	- Follow `go fmt` and `go vet`. We recommend using `golangci-lint` locally.
	- Write clear, idiomatic Go code.
	- Keep functions small and focused.
- **TypeScript**/**React**:
	- Adhere to our ESLint and Prettier configurations. Run `npm run lint` and `npm run format`.
	- Use functional components with Hooks.
	- Prioritize readability and maintainability.
- **General**:
	- **Performance**: Be mindful of performance implications, especially in core data paths.
	- **Security**: Always consider security implications. Avoid hardcoding secrets, validate all inputs, and sanitize outputs.
	- **Modularity**: Design components to be loosely coupled and highly cohesive.
	- **Error Handling:** Handle errors gracefully and provide informative error messages.
	- **Logging**: Use structured logging (log/slog or a similar library) for important events and errors.

### <a name="testing"></a>Testing

- **Unit Tests**: Write unit tests for all new functions and components.
- **Integration Tests**: Add integration tests for service interactions where appropriate.
- **End-to-End Tests**: Contributions involving major features should ideally be accompanied by or consider how they'd be covered by E2E tests.
- Ensure all existing tests pass before submitting a PR.

## <a name="project-structure-overview"></a>Project Structure Overview

Familiarize yourself with the overall project structure:

```
korun.io/
├── .github/                  # GitHub Actions CI/CD workflows
├── backend/                  # Go microservices
│   ├── api-gateway/
│   ├── auth-service/
│   ├── secrets-service/
│   ├── audit-service/
│   ├── notification-service/
│   ├── webhook-service/
│   └── shared/               # Shared Go modules (events, messaging, models, utils)
├── frontend/                 # Next.js web dashboard
│   └── web-app/
├── cli/                      # Go CLI tool
│   └── korun-cli/
├── infrastructure/           # Dockerfiles, Kubernetes manifests, Terraform
│   ├── docker/
│   ├── kubernetes/
│   └── terraform/
├── docs/                     # Additional documentation
├── docker-compose.yml        # Local development setup
├── go.mod                    # Root Go module
├── .editorconfig             # Editor configuration
├── CODE_OF_CONDUCT.md        # Code of conduct statement
├── CONTRIBUTING.md           # this file
└── README.md
```

## <a name="need-help"></a>Need Help?

If you have any questions or get stuck, don't hesitate to:

- Open an issue on GitHub.

Thank you for helping us build Korun.io!
