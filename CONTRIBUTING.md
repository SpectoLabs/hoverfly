# Contributing to Hoverfly

Thanks for your interest in contributing to Hoverfly! This guide will help you get started.

## Getting Started

### Prerequisites

- **Go 1.26.1+** — install from [golang.org/dl](https://golang.org/dl)
- **Ruby** and **Python** — needed for some middleware tests
  ```bash
  # macOS
  brew install ruby python
  ```

### Clone and Build

```bash
git clone https://github.com/SpectoLabs/hoverfly.git
# or your fork: git clone https://github.com/<your_username>/hoverfly.git
cd hoverfly
make build
```

Binaries are output to the `target/` directory.

### Running Tests

```bash
# All tests (unit + functional + vet)
make test

# Hoverfly unit tests only
make hoverfly-test

# Hoverctl unit tests only
make hoverctl-test

# Functional tests (requires built binaries in target/)
make hoverfly-functional-test
make hoverctl-functional-test

# Run a single test
cd core && go test -v ./... -run TestFunctionName
```

Tests use the [Ginkgo/Gomega](https://onsi.github.io/ginkgo/) BDD framework. Functional tests in `functional-tests/` spin up actual Hoverfly instances.

## How to Contribute

### Finding Something to Work On

- Look for issues labelled **[good first issue](https://github.com/SpectoLabs/hoverfly/labels/good%20first%20issue)** — these are beginner-friendly.
- Browse **[help wanted](https://github.com/SpectoLabs/hoverfly/labels/help%20wanted)** for tasks where maintainers welcome contributions.
- If you have an idea for a new feature, open a [feature request](https://github.com/SpectoLabs/hoverfly/issues/new?template=feature_request.md) first to discuss it before writing code.

### Workflow

1. **Fork** the repository.
2. **Create a feature branch** on your fork (`git checkout -b my-feature`).
3. **Make your changes** and add tests where appropriate.
4. **Run the tests** (`make test`) and ensure they pass.
5. **Format and vet** your code:
   ```bash
   make fmt
   make vet
   ```
6. **Commit** with a clear message describing what and why.
7. **Open a pull request** against the `master` branch.

### Pull Request Guidelines

- Keep PRs focused — one logical change per PR.
- Include a description of **what** you changed, **why**, and **how to test** it.
- If your PR fixes an issue, reference it (e.g. `Fixes #123`).
- Be responsive to review feedback.

### Code Style

- Follow standard Go conventions (`gofmt`, `go vet`).
- Keep functions focused and reasonably sized.
- Write tests for new functionality — both unit tests and functional tests where applicable.

## Asking Questions

If you have questions about the codebase, how something works, or whether a change would be welcome:

- Open a [GitHub Discussion](https://github.com/SpectoLabs/hoverfly/discussions)
- Use the [question issue template](https://github.com/SpectoLabs/hoverfly/issues/new?template=question-about-Hoverfly.md)

## License

By contributing to Hoverfly, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).
