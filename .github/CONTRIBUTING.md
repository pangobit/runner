# Contributing to Runner

Thank you for considering contributing to Runner! This document outlines the process for contributing to this project.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How Can I Contribute?

### Reporting Bugs

- **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/pangobit/runner/issues).
- If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/pangobit/runner/issues/new). Be sure to include a **title and clear description**, as much relevant information as possible, and a **code sample** or an **executable test case** demonstrating the expected behavior that is not occurring.

### Suggesting Enhancements

- Open a new issue with a clear title and detailed description of the suggested enhancement.
- Provide specific examples and explanations of how this enhancement would be useful.

### Pull Requests

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Development Setup

1. Fork and clone the repository
2. Install Go (version 1.20+)
3. Run `go mod download` to install dependencies
4. Make your changes
5. Run tests with `go test ./...`
6. Ensure code passes `go vet ./...`

## Coding Standards

- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Write [good commit messages](https://chris.beams.io/posts/git-commit/)
- Include appropriate tests
- Update documentation as needed

## Security

For security issues, please see our [Security Policy](SECURITY.md).

## License

By contributing, you agree that your contributions will be licensed under the project's license.

Thank you for your contributions!
