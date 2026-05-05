# Commands

## Connection Management

| Command | Description | Example |
|---------|-------------|---------|
| `init <name> <type> <conn-string> [schema]` | Create new database connection | `pam create mydb postgres "postgresql://..."` |
| `use/switch <name>` | Switch to a different connection | `pam use production` |
| `status` | Show current active connection | `pam status` |
| `list connections` | List all configured connections | `pam list connections` |

## Query Operations

| Command | Description | Example |
|---------|-------------|---------|
| `add <name> [sql]` | Add a new saved query | `pam add users "SELECT * FROM users"` |
| `remove <name\|id>` | Remove a saved query | `pam remove users` or `pam remove 3` |
| `list queries` | List all saved queries | `pam list queries` |
| `list queries --oneline` | lists each query in one line | `pam list -o` |
| `list queries <searchterm>` | lists queries containing search term | `pam list employees` |
| `run <name\|id\|sql>` | Execute a query | `pam run users` or `pam run 2` |
| `run` | Create and run a new query | `pam run` |
| `run --edit` | Edit query before running | `pam run users --edit` |
| `run --last`, `-l` | Re-run last executed query | `pam run --last` |
| `run --param` | run with named params | `pam run --name PAM` |
| `shell` | Interactive query REPL (alias: `repl`) | `pam shell` |


## Database Exploration

| Command | Description | Example |
|---------|-------------|---------|
| `explore` | List all tables and views in multi-column format | `pam explore` |
| `explore <table> [-l N]` | Query a table with optional row limit | `pam explore employees --limit 100` |
| `explain <table> [-d N] [-c]` | Visualize foreign key relationships | `pam explain employees --depth 2` |
| `tables` | List all tables in using the results view, access with Enter| `pam tables` |

## Configuration

| Command | Description | Example |
|---------|-------------|---------|
| `config` | Edit main configuration file | `pam config` |
| `edit` | Edit all queries for current connection | `pam edit` |
| `edit <name\|id>` | Edit a single named query | `pam edit 3` |
| `remove --connection <name>` | Remove a db connection | `pam remove --conection dev4`` |
| `help [command]` | Show help information | `pam help run` |

