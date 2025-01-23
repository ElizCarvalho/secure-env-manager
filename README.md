# Secure ENV Manager

A CLI tool for securely managing .env files using AES-GCM encryption.

## Authentication Strategy

The application uses a simple but effective authentication strategy:

1. **Environment Variables**: Initial credentials are set through environment variables:
   - `SECURE_ENV_INITIAL_USER`: Username for authentication
   - `SECURE_ENV_INITIAL_PASS`: Password for authentication and encryption

2. **Command Line Authentication**: Users must provide credentials when running the program:
```bash
make run USER=username PASS=password
```

3. **Validation Process**:
   - Credentials are validated against the environment variables
   - The password is also used as the encryption key for AES-GCM
   - All team members must use the same credentials to ensure file sharing

4. **Security Measures**:
   - Credentials are never stored in plain text
   - The configuration file is protected against modifications
   - Each session requires re-authentication

### Why Pre-configured Credentials?

The decision to use pre-configured credentials instead of a user registration system was made for several reasons:

1. **Team Synchronization**: Since the tool is designed for team use, having shared credentials ensures all team members can encrypt/decrypt files using the same key.
2. **Simplified Key Management**: Using a single set of credentials eliminates the complexity of managing multiple encryption keys and sharing files between users.
3. **Security by Process**: The security model relies on proper credential distribution through secure team channels, rather than a potentially vulnerable user registration system.
4. **Reduced Attack Surface**: By not implementing user registration and management, we eliminate potential security vulnerabilities associated with these features.

## Running the Program

1. Clone the repository:
```bash
git clone [repo-url]
cd secure-env-manager
```

2. Run the program:
```bash
make run USER=username PASS=password
```

**Important:** 
- Use the team's shared credentials. Everyone should use the same credentials to ensure they can decrypt each other's files.
- The `config.json` file is protected against changes (read-only) to ensure credential integrity.
- Obtain the credentials from your team before using the program.

## How to Use

The program offers two main operations:

1. Encrypt .env file:
```bash
make run USER=username PASS=password
# Choose option 1
# Enter project name (e.g., myproject)
# Enter path to .env file to encrypt (e.g., ~/dev/myproject/.env)
```

2. Decrypt file:
```bash
make run USER=username PASS=password
# Choose option 2
# Enter project name
# Enter path to file to update (e.g., ~/dev/myproject/.env)
```

## Uninstallation

To completely remove the program:
```bash
make uninstall
```

## File Structure

```
secure-env-manager/
├── cmd/
│   └── secure-env/     # Application entry point
│       ├── main.go
│       └── main_test.go
├── internal/
│   ├── config/         # Application configuration
│   │   ├── config.go
│   │   └── config_test.go
│   ├── crypto/         # Encryption functions
│   │   ├── crypto.go
│   │   └── crypto_test.go
│   └── handler/        # Command handlers
│       ├── handler.go
│       └── handler_test.go
├── test-output/        # Test coverage reports
├── .cache/            # Go cache (ignored by Git)
├── Dockerfile         # Container configuration
├── Makefile          # Automation commands
├── go.mod            # Go dependencies
├── go.sum            # Dependencies checksums
└── README.md         # This file
```

## Security

- Team-shared credentials
- AES-GCM encryption for .env files
- Password also used as encryption key
- Decrypted files automatically ignored by Git
- Credentials file protected against changes

## Make Commands

```bash
make run USER=username PASS=password  # Run the program
make uninstall                       # Remove the program
make help                           # Show help
```

## Requirements

- Docker 

## Contributing

We love contributions from our community! If you want to contribute to the project:

1. Read our [Contributing Guide](CONTRIBUTING.md) to understand our development process
2. Learn about our semantic commit pattern that we use to automate releases
3. Fork the project and submit your Pull Request!

All contributions are welcome, from documentation fixes to new features. 