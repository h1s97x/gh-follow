# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x     | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of GH-Follow seriously. If you have discovered a security vulnerability, we appreciate your help in disclosing it to us in a responsible manner.

### How to Report

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via GitHub's private vulnerability reporting feature:

1. Go to the [Security Advisories](https://github.com/h1s97x/gh-follow/security/advisories) page
2. Click "Report a vulnerability"
3. Fill out the form with details about the vulnerability

### What to Include

Please include the following information in your report:

- Type of vulnerability
- Full path of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit)
- Any special configuration required to reproduce the issue
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the vulnerability

### Response Timeline

- We will acknowledge your report within 48 hours
- We will provide a more detailed response within 7 days
- We will keep you informed of the progress towards a fix

### Security Best Practices

When using GH-Follow, please follow these security guidelines:

1. **Token Security**: Never share your GitHub token or commit it to version control
2. **File Permissions**: The tool sets file permissions to 0600 for sensitive files
3. **Gist Sync**: When using Gist sync, use private Gists for sensitive data
4. **Regular Updates**: Keep the tool updated to the latest version

## Security Features

GH-Follow implements the following security measures:

- **Secure Token Handling**: Tokens are retrieved via `gh auth token` and never stored by the tool
- **Encrypted Storage**: All local files use restricted permissions (0600)
- **No Hardcoded Secrets**: The codebase contains no hardcoded credentials or API keys
- **Input Validation**: All user inputs are validated before processing

## Known Security Considerations

- **Local Storage**: Follow lists are stored locally in JSON format. Ensure your system is secure.
- **Gist Sync**: If enabled, data is stored in GitHub Gists. Use private Gists for sensitive data.
- **GitHub Token**: The tool requires a GitHub token with `user` scope for sync operations.

Thank you for helping keep GH-Follow and its users safe!
