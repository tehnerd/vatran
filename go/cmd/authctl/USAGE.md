# authctl - Katran Authentication User Management

A CLI tool for managing users in the Katran authentication database.

## Building

```bash
go build -o authctl ./go/cmd/authctl/
```

## Commands

### Add a User

```bash
# Interactive (prompts for password securely)
authctl add -username admin -db /var/lib/katran/auth.db

# Non-interactive (for scripts)
authctl add -username admin -password 'secretpass' -db /var/lib/katran/auth.db

# With custom bcrypt cost
authctl add -username admin -bcrypt-cost 14 -db /var/lib/katran/auth.db
```

**Options:**
- `-username` (required): Username to create
- `-password`: Password (prompts if not provided)
- `-db`: Path to SQLite database (default: `/var/lib/katran/auth.db`)
- `-bcrypt-cost`: bcrypt cost factor (default: 12)

**Notes:**
- Password must be at least 8 characters
- When prompting, password confirmation is required

### List Users

```bash
authctl list -db /var/lib/katran/auth.db
```

**Example output:**
```
ID  USERNAME  CREATED              UPDATED
--  --------  -------              -------
1   admin     2024-01-15 10:30:00  2024-01-15 10:30:00
2   operator  2024-01-16 14:22:00  2024-01-16 14:22:00

Total: 2 user(s)
```

**Options:**
- `-db`: Path to SQLite database (default: `/var/lib/katran/auth.db`)

### Change Password

```bash
# Interactive (prompts for new password)
authctl passwd -username admin -db /var/lib/katran/auth.db

# Non-interactive
authctl passwd -username admin -password 'newpassword' -db /var/lib/katran/auth.db
```

**Options:**
- `-username` (required): Username to update
- `-password`: New password (prompts if not provided)
- `-db`: Path to SQLite database (default: `/var/lib/katran/auth.db`)
- `-bcrypt-cost`: bcrypt cost factor (default: 12)

### Delete a User

```bash
# With confirmation prompt
authctl delete -username olduser -db /var/lib/katran/auth.db

# Skip confirmation (for scripts)
authctl delete -username olduser -force -db /var/lib/katran/auth.db
```

**Options:**
- `-username` (required): Username to delete
- `-db`: Path to SQLite database (default: `/var/lib/katran/auth.db`)
- `-force`: Skip confirmation prompt

### Help

```bash
authctl help
```

## Database Location

The default database path is `/var/lib/katran/auth.db`. This should match the `database_path` in your Katran server configuration:

```yaml
server:
  auth:
    enabled: true
    database_path: "/var/lib/katran/auth.db"
```

## Security Notes

- Passwords are hashed using bcrypt before storage
- The default bcrypt cost of 12 provides good security/performance balance
- Increase cost for higher security (slower hashing): `-bcrypt-cost 14`
- Password input is hidden when using interactive prompts
- Avoid passing passwords via command line in production (visible in process list)
