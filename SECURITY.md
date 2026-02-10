# Security Policy

## Overview

Security is a top priority for ENVM. As an environment variable management tool that handles sensitive data, we take security vulnerabilities seriously and appreciate the community's help in responsibly disclosing any security issues.

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

We recommend always using the latest stable version to ensure you have the latest security patches and updates.

## Reporting a Vulnerability

**Please DO NOT report security vulnerabilities through public GitHub issues.**

### Reporting Process

If you discover a security vulnerability, please report it by:

1. **Email**: Send details to **hezronnyamboga01@gmail.com** or create a private security advisory on GitHub
2. **GitHub Security Advisory**: Use GitHub's [private vulnerability reporting feature](https://github.com/envm-org/envm/security/advisories/new)

### What to Include in Your Report

To help us understand and resolve the issue quickly, please include:

- **Type of vulnerability** (e.g., authentication bypass, SQL injection, information disclosure)
- **Affected component(s)** (API, CLI, specific module)
- **Affected version(s)**
- **Step-by-step instructions** to reproduce the vulnerability
- **Proof of concept** (code, screenshots, or video demonstration)
- **Potential impact** of the vulnerability
- **Suggested fix** (if you have one)
- **Any relevant logs or error messages**
- **Your contact information** for follow-up questions

### Example Report Format

```
Title: [Brief description of vulnerability]

Severity: [Critical/High/Medium/Low]

Component: [API/CLI/Database/etc.]

Version: [Affected version(s)]

Description:
[Detailed description of the vulnerability]

Steps to Reproduce:
1. [Step 1]
2. [Step 2]
3. [Step 3]

Impact:
[Description of potential impact]

Proof of Concept:
[Code, screenshots, or commands demonstrating the vulnerability]

Suggested Fix:
[Your recommendation, if any]

Reporter: [Your name/handle (optional)]
Contact: [Your email for follow-up]
```

## Response Timeline

- **Initial Response**: Within 48 hours
- **Confirmation**: Within 5 business days
- **Status Updates**: Every 7 days until resolution
- **Fix Release**: Depends on severity and complexity
  - Critical: Within 7 days
  - High: Within 14 days
  - Medium: Within 30 days
  - Low: Next scheduled release

## Disclosure Policy

We follow **Coordinated Vulnerability Disclosure**:

1. You report the vulnerability privately
2. We acknowledge receipt within 48 hours
3. We work with you to understand and reproduce the issue
4. We develop and test a fix
5. We release a security patch
6. After the patch is available and users have had time to update (typically 7-14 days), we publish:
   - Security advisory on GitHub
   - CVE (if applicable)
   - Credit to the reporter (with your permission)

### Public Disclosure Timeline

- **Critical/High**: 90 days from initial report or immediately after a fix is available (whichever comes first)
- **Medium**: 120 days from initial report
- **Low**: 180 days from initial report

## Security Best Practices for Users

### For Self-Hosted Deployments

1. **Keep ENVM Updated**
   - Regularly update to the latest version
   - Subscribe to release notifications
   - Monitor security advisories

2. **Secure Your Deployment**
   - Use HTTPS/TLS for all connections
   - Enable authentication and authorization
   - Use strong, unique passwords
   - Implement network segmentation
   - Restrict API access to trusted networks
   - Enable audit logging

3. **Encryption**
   - Ensure encryption at rest is enabled
   - Use strong encryption keys
   - Rotate encryption keys regularly
   - Store keys securely (e.g., using a vault)

4. **Database Security**
   - Use strong database credentials
   - Limit database access to necessary services only
   - Enable database encryption
   - Regular database backups with encryption

5. **Environment Configuration**
   - Never commit `.env` files or secrets to version control
   - Use environment-specific configurations
   - Limit environment variable exposure
   - Regularly rotate credentials

6. **Access Control**
   - Implement role-based access control (RBAC)
   - Follow principle of least privilege
   - Regularly audit user access and permissions
   - Remove inactive users promptly

7. **Monitoring and Logging**
   - Enable comprehensive audit logging
   - Monitor for suspicious activities
   - Set up alerts for security events
   - Regularly review logs

### For CLI Users

1. **Keep CLI Updated**
   ```bash
   envm version
   envm upgrade  # If available
   ```

2. **Secure Local Storage**
   - Protect your `.envm` files with appropriate file permissions
   - Don't share `.envm` files containing production secrets
   - Use encrypted `.envm` files for sensitive data

3. **Authentication**
   - Use strong API tokens
   - Rotate tokens regularly
   - Never share authentication tokens
   - Store tokens securely (use system keychain if available)

4. **Network Security**
   - Only sync with trusted servers
   - Verify server certificates
   - Use VPN when accessing from untrusted networks

## Known Security Considerations

### Encryption

- ENVM uses **AES-256** encryption for stored environment variables
- Encryption at rest is enabled by default
- Transport layer security (TLS 1.3) is required for all network communications
- Keys are derived using secure key derivation functions (e.g., PBKDF2 or Argon2)

### Authentication

- JWT tokens expire after configurable period (default: 24 hours)
- Refresh tokens are rotated on each use
- Failed login attempts are rate-limited
- Account lockout after consecutive failed attempts

### Database Security

- All database queries use prepared statements (protection against SQL injection)
- Database credentials are stored as environment variables, never in code
- Supports connection encryption to database

### Audit Logging

- All sensitive operations are logged
- Logs include: timestamp, user, action, resource, IP address
- Logs are immutable and tamper-resistant
- Configurable log retention periods

## Security Features

### Current Implementation

- [x] Encryption at rest (AES-256)
- [x] TLS/HTTPS for all communications
- [x] JWT-based authentication
- [x] Role-based access control (RBAC)
- [x] Audit logging
- [x] Rate limiting
- [x] Input validation and sanitization
- [x] Prepared statements for database queries
- [x] Secure password hashing (bcrypt)

### Planned Features

- [ ] Integration with HashiCorp Vault
- [ ] AWS Secrets Manager integration
- [ ] Azure Key Vault integration
- [ ] Secret rotation automation
- [ ] Advanced threat detection
- [ ] Security scanning in CI/CD
- [ ] Automated security testing

## Compliance

ENVM is designed with compliance in mind, considering:

- **GDPR**: Data privacy and user rights
- **SOC 2**: Security controls and procedures
- **HIPAA**: Healthcare data protection (when applicable)
- **PCI DSS**: Payment card data security (when applicable)

Note: Compliance adherence depends on proper configuration and deployment practices.

## Security Tools and Scanning

### Development

- Static analysis: `golangci-lint` with security checks (gosec)
- Dependency scanning: `go mod` with vulnerability database
- Secret scanning: Pre-commit hooks to prevent secret commits

### CI/CD Pipeline

- Automated security testing
- Container image scanning
- Dependency vulnerability scanning
- SAST (Static Application Security Testing)
- DAST (Dynamic Application Security Testing) on staging

## Bug Bounty Program

We currently do not have a formal bug bounty program, but we deeply appreciate security researchers' efforts. We recognize contributors in our security advisories and release notes (with permission).

## Hall of Fame

We thank the following security researchers for responsibly disclosing vulnerabilities:

<!-- This section will be updated as vulnerabilities are reported and fixed -->

*No vulnerabilities have been reported yet.*

## Contact

- **Security Email**: hezronnyamboga01@gmail.com
- **GitHub Security Advisories**: https://github.com/envm-org/envm/security/advisories
- **Project Team**: security@envm-org.dev

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [CWE (Common Weakness Enumeration)](https://cwe.mitre.org/)
- [CVE (Common Vulnerabilities and Exposures)](https://cve.mitre.org/)
- [NIST Cybersecurity Framework](https://www.nist.gov/cyberframework)

---

**Thank you for helping keep ENVM and our community safe!** 🔒

Last Updated: February 10, 2026
