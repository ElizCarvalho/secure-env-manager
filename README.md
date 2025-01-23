# Secure ENV Manager

A CLI tool for securely managing .env files using AES-GCM encryption.

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