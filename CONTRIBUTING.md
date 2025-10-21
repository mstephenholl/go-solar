# Contributing to go-solar

Thank you for your interest in contributing to go-solar! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Code Quality](#code-quality)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)

## Getting Started

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-solar.git
   cd go-solar
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/mstephenholl/go-solar.git
   ```

## Development Setup

### Prerequisites

- Go 1.21 or later
- Make (optional, but recommended)
- Git

### Install Development Tools

```bash
make install-tools
```

This will install:
- golangci-lint
- staticcheck
- goimports
- govulncheck

### Verify Your Setup

```bash
make ci-quick
```

This runs formatting checks, vet, and tests.

## Making Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-new-calculation`
- `fix/correct-polar-cases`
- `docs/update-examples`
- `refactor/improve-performance`

### Code Style

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names
- Add comments for exported functions
- Document complex algorithms
- Keep functions focused and small

### Writing Code

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following these guidelines:
   - Write clear, concise code
   - Add tests for new functionality
   - Update documentation as needed
   - Follow existing code patterns

3. Commit your changes:
   ```bash
   git add .
   git commit -m "feat: add new calculation method"
   ```

### Commit Message Format

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(elevation): add support for custom elevation angles
fix(hourangle): correct polar night calculation
docs(readme): update installation instructions
test(sunrise): add edge case tests for polar regions
```

## Testing

### Run All Tests

```bash
make test
```

### Run Tests with Coverage

```bash
make coverage
```

This generates a coverage report in `coverage/coverage.html`.

### Run Benchmarks

```bash
make test-bench
```

### Writing Tests

- Use table-driven tests where appropriate
- Test edge cases (polar regions, extreme dates, etc.)
- Aim for 100% coverage on new code
- Use descriptive test names
- Add example tests for documentation

**Example:**

```go
func TestNewFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    float64
        expected float64
    }{
        {"zero", 0.0, 0.0},
        {"positive", 5.0, 25.0},
        {"negative", -5.0, 25.0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := NewFunction(tt.input)
            if !AlmostEqual(got, tt.expected, 1e-10) {
                t.Errorf("got %v, want %v", got, tt.expected)
            }
        })
    }
}
```

## Code Quality

### Format Code

```bash
make fmt
```

### Run Linters

```bash
make lint
```

### Run Security Checks

```bash
make security
```

### Pre-commit Checklist

Before committing, ensure:

```bash
make pre-commit
```

This runs:
- Code formatting
- `go vet`
- Tests

## Submitting Changes

### Pull Request Process

1. Update your branch with the latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/master
   ```

2. Push your changes:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create a Pull Request on GitHub

4. Ensure all CI checks pass

5. Wait for review and address feedback

### Pull Request Guidelines

- Provide a clear description of the changes
- Reference any related issues
- Include test results if applicable
- Update documentation if needed
- Ensure CI passes
- Keep PRs focused (one feature/fix per PR)

### PR Description Template

```markdown
## Description
Brief description of what this PR does

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
Describe the tests you ran and how to reproduce

## Checklist
- [ ] Tests pass locally
- [ ] Added tests for new functionality
- [ ] Updated documentation
- [ ] Code follows project style
- [ ] No breaking changes (or documented)
```

## Release Process

Releases are automated using GoReleaser:

1. Update version in relevant files
2. Create and push a new tag:
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```
3. GitHub Actions will automatically create a release

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

## Questions?

- Open an issue for bug reports or feature requests
- Check existing issues before creating new ones
- Be respectful and constructive in discussions

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professional communication

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to go-solar! 🌞
