# Contributing to Secure ENV Manager

We love your input! We want to make contributing to Secure ENV Manager as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

1. Fork the repo and create your branch from `main`
2. If you've added code that should be tested, add tests
3. If you've changed APIs, update the documentation
4. Ensure the test suite passes
5. Make sure your code follows our coding standards
6. Issue that pull request!

## Semantic Commit Messages

We use semantic commits to automate our release process. Each commit message should be structured as follows:

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

### Types

- `feat`: A new feature (increments MINOR version)
- `fix`: A bug fix (increments PATCH version)
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
- `ci`: Changes to our CI configuration files and scripts

### Breaking Changes

When a commit includes a breaking change, it MUST:
1. Have a `!` after the type/scope: `feat!: remove support for Node 6`
2. Include `BREAKING CHANGE:` in the footer:

```
feat!: remove support for Node 6

BREAKING CHANGE: use JavaScript features not available in Node 6.
```

### Examples

```bash
feat: add support for encrypted files
fix: resolve memory leak in encryption process
docs: update installation instructions
style: format code according to new style guide
refactor: restructure encryption module
perf: improve file reading performance
test: add tests for decryption process
build: update dependencies
ci: add GitHub Actions workflow
```

## Release Process

We use standard-version to manage our releases. The process is automated:

```bash
# For bug fixes
npm run release:patch    # 1.0.0 -> 1.0.1

# For new features
npm run release:minor    # 1.0.0 -> 1.1.0

# For breaking changes
npm run release:major    # 1.0.0 -> 2.0.0
```

The release process will:
1. Update version numbers
2. Generate CHANGELOG.md
3. Create a git tag
4. Push changes and tags to GitHub

## Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the CHANGELOG.md with notes on your changes
3. The PR will be merged once you have the sign-off of at least one maintainer

## Any Questions?

Feel free to open an issue with your question or suggestion. We're always happy to help! 