# Features

## Query Management

Save, organize, and execute your SQL queries with ease.

```bash
# Add queries with auto-incrementing IDs
pam add daily_report "SELECT * FROM sales WHERE date = CURRENT_DATE"
pam add user_count "SELECT COUNT(*) FROM users"
pam add employees "SELECT TOP 10 * FROM employees ORDER BY last_name"

# Add parameterized queries with :param|default syntax
pam add emp_by_salary "SELECT * FROM employees WHERE salary > :min_sal|30000"
pam add search_users "SELECT * FROM users WHERE name LIKE :name|P% AND status = :status|active"

# When creating queries with params and not default, pam will prompt you for the param value every time you run the query
pam add search_by_name "SELECT * FROM employees where first_name = :name"

# Run parameterized queries with named parameters (order doesn't matter!)
pam run emp_by_salary --min_sal 50000
pam run search_users --name Michael --status active
# Or use positional args (must match SQL order)
pam run search_users Michael active

# List all saved queries
pam list

# Search for specific queries
pam list emp    # Finds queries with 'emp' in name or SQL
pam list employees --oneline # displays each query in one line

# Run by name or ID
pam run daily_report
pam run 2

# Edit query before running (great for testing parameter values)
pam run emp_by_salary --edit
```

<img width="1188" height="714" alt="image" src="https://github.com/user-attachments/assets/016c7a61-ace4-49cc-9375-564ee6089899" />

## TUI Table Viewer

Navigate query results with Vim-style keybindings, update cells in-place, delete rows and copy data

<img width="1173" height="709" alt="image" src="https://github.com/user-attachments/assets/3959011b-532f-4374-a86d-a39217cd39f0" />

**Key Features:**
- Syntax-highlighted SQL display
- Column type indicators
- Primary key markers
- Live cell editing
- Visual selection mode

## Connection Switching

Manage multiple database connections and switch between them instantly.

```bash
# List all connections
pam list connections
pam use production
```
Display current connection and check if it is reachable
```
pam status
```
<div align=center>
  <img width="523" height="582" alt="image" src="https://github.com/user-attachments/assets/4046f6cd-376e-45c0-bcfd-20484e34470b" />
</div>

## Database Exploration

Explore your database schema and visualize relationships between tables.

```bash
# List all tables and views in multi-column format
pam explore

# Query a table directly
pam explore employees --limit 100

# Open tables in the results view, use Enter to query everything in the table
pam tables

# Visualize foreign key relationships
pam explain employees
pam explain employees --depth 2    # Show relationships 2 levels deep
```

<img width="860" height="139" alt="image" src="https://github.com/user-attachments/assets/4cea0f4d-d3b9-4173-8b42-6ee6b289cc7b" />

**Note:** The `pam explain` command is currently a work in progress and may change in future versions.

---

## Editor Integration

PAM uses your `$EDITOR` environment variable for editing queries and UPDATE/DELETE statements.

<div align=center>
  <img width="448" height="238" alt="image" src="https://github.com/user-attachments/assets/f416f41a-8ec3-4a35-86e7-0bba6596f75f" />
</div>

```bash
# Set your preferred editor (example in bash)
export EDITOR=vim
export EDITOR=nano
export EDITOR=code
```

You can also use the editor to edit queries before running them

```bash
# Edit existing query before running
pam run daily_report --edit

# Create and run a new query on the fly
pam run

# Re-run the last executed query
pam run --last

# Edit all queries at once
pam edit queries

# Edit a specific query
pam edit recent_users
```

## Interactive Shell

Run queries in an interactive REPL with persistent connection, history, and multi-line support.

```bash
pam shell          # or: pam repl
```

**Example session:**
```bash
pam@mydb> select * from users limit 5;
pam@mydb> list_users --status active
pam@mydb> --last
pam@mydb> list user
pam@mydb> status
```

**Meta-commands:**

| Command | Description |
|---------|-------------|
| `exit`, `quit`, `\q` | Exit the shell |
| `help`, `\h` | Show help |
| `status` | Show connection info |
| `list`, `ls`, `\l` | List queries or connections |
| `tables`, `\dt` | List tables in the current database |
| `config` | Open the PAM config file in `$EDITOR` |

Multi-line: type SQL without trailing `;` to continue. End with `;` or press Enter on blank line to execute.
