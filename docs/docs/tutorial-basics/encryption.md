---
sidebar_position: 4
---

# Encryption

Learn how to protect sensitive values with encryption.

## Why Encryption?

Environment variables often contain sensitive data:
- Database credentials
- API keys and tokens
- OAuth secrets
- Payment gateway credentials

ENVM uses **AES-256-GCM** encryption to protect these values at rest.

## Setting Up Encryption

### Initialize Encryption

When you first encrypt a value, ENVM will prompt you to set up a master password:

```bash
envm set API_SECRET="my-secret-key" --encrypt
# Enter master password: ********
# Confirm master password: ********
# ✓ Encryption initialized
# ✓ API_SECRET encrypted and saved
```

### Master Password

The master password is used to derive the encryption key. It is:
- **Never stored** - you must remember it
- **Required** to decrypt values
- **Shareable** with team members who need access

:::caution Important
If you lose your master password, encrypted values cannot be recovered. Store it securely!
:::

## Encrypting Values

### Encrypt on Set

```bash
# Encrypt a new variable
envm set DATABASE_PASSWORD="super-secret" --encrypt

# Encrypt multiple values
envm set AWS_ACCESS_KEY="AKIA..." AWS_SECRET_KEY="..." --encrypt
```

### Encrypt Existing Variables

```bash
# Encrypt an existing variable
envm encrypt DATABASE_PASSWORD

# Encrypt all variables in a group
envm encrypt --group=secrets

# Encrypt by pattern
envm encrypt "AWS_*"
```

### Encrypted Values in `.envm`

Encrypted values are stored with the `ENC:` prefix:

```bash
[database]
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_PASSWORD=ENC:v1:aes256gcm:base64encodedciphertext...

[aws]
AWS_ACCESS_KEY=ENC:v1:aes256gcm:base64encodedciphertext...
AWS_SECRET_KEY=ENC:v1:aes256gcm:base64encodedciphertext...
```

## Decrypting Values

### View Decrypted Value

```bash
# Decrypt and display a value (requires password)
envm decrypt DATABASE_PASSWORD
# Enter master password: ********
# DATABASE_PASSWORD: super-secret
```

### Decrypt for Use

When running commands, encrypted values are automatically decrypted:

```bash
envm run -- node app.js
# All encrypted values are decrypted in the environment
```

### Permanent Decryption

To permanently remove encryption from a value:

```bash
# Decrypt and store as plain text
envm decrypt DATABASE_PASSWORD --permanent
# Warning: This will store the value in plain text. Continue? [y/N]
```

## Key Management

### Rotating Keys

Periodically rotate your encryption key:

```bash
envm key rotate
# Enter current master password: ********
# Enter new master password: ********
# Confirm new master password: ********
# ✓ All encrypted values re-encrypted with new key
```

### Backing Up Keys

Export a backup of your encryption key:

```bash
# Export encrypted key backup
envm key export --output=key-backup.enc
# This file is encrypted with your master password

# Import key backup
envm key import key-backup.enc
```

## Team Sharing

### Sharing the Master Password

For team environments, share the master password securely:

1. Use a password manager (1Password, Bitwarden, etc.)
2. Share through encrypted channels
3. Use ENVM's team features (with sync server)

### Per-Environment Keys

Use different keys for different environments:

```bash
# Set environment-specific key
envm key set --env=production

# Different keys for dev vs production
ENVM_MASTER_PASSWORD=dev-password envm list     # Development
ENVM_MASTER_PASSWORD=prod-password envm list    # Production
```

## Environment Variable

Store the master password in an environment variable for automation:

```bash
# Set the master password
export ENVM_MASTER_PASSWORD="your-master-password"

# Now commands work without prompts
envm set API_KEY="secret" --encrypt
envm run -- npm start
```

:::warning CI/CD Security
In CI/CD pipelines, use your platform's secret management:
- GitHub Actions: `secrets.ENVM_MASTER_PASSWORD`
- GitLab CI: CI/CD Variables
- CircleCI: Environment Variables
:::

## Encryption Algorithms

ENVM supports multiple encryption algorithms:

| Algorithm | Description | Default |
|-----------|-------------|---------|
| `aes-256-gcm` | AES-256 with Galois/Counter Mode | ✓ |
| `chacha20-poly1305` | ChaCha20 with Poly1305 MAC | |

Configure in `.envm.config.yaml`:

```yaml
encryption:
  algorithm: aes-256-gcm
  key_derivation: argon2id
```

## Next Steps

- [Running Commands](./running-commands) - Use encrypted variables in your app
- [Multi-Environment](./multi-environment) - Manage different environments

