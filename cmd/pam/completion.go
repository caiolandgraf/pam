package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caiolandgraf/pam/internal/config"
	"github.com/caiolandgraf/pam/internal/styles"
)

func (a *App) handleCompletion() {
	args := os.Args[2:]
	if len(args) == 0 {
		printError("Usage: pam completion <bash|zsh|fish> [--install]")
	}

	install := false
	shell := ""
	for _, arg := range args {
		switch strings.ToLower(arg) {
		case "--install":
			install = true
		case "bash", "zsh", "fish":
			shell = strings.ToLower(arg)
		}
	}

	if shell == "" {
		printError("Usage: pam completion <bash|zsh|fish> [--install]")
	}

	script := ""
	switch shell {
	case "bash":
		script = bashCompletionScript
	case "zsh":
		script = zshCompletionScript
	case "fish":
		script = fishCompletionScript
	default:
		printError("Unsupported shell: %s. Use bash, zsh, or fish.", shell)
	}

	if !install {
		fmt.Print(script)
		return
	}

	destPath := completionInstallPath(shell)
	if destPath == "" {
		printError("Could not determine completion install path for %s", shell)
	}

	// Ensure parent directory exists
	dir := filepath.Dir(destPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		printError("Could not create directory %s: %v", dir, err)
	}

	if err := os.WriteFile(destPath, []byte(script), 0644); err != nil {
		printError("Could not write completion file: %v", err)
	}

	fmt.Println(styles.Success.Render("✓ Completion installed successfully"))
	fmt.Println(styles.Faint.Render("  File: " + destPath))
	fmt.Println()

	switch shell {
	case "bash":
		fmt.Println(styles.Faint.Render("  Reload with: source " + destPath))
		fmt.Println(styles.Faint.Render("  Or restart your terminal."))
	case "zsh":
		zshDir := filepath.Dir(destPath)
		fmt.Println(styles.Faint.Render("  Make sure your ~/.zshrc contains:"))
		fmt.Println()
		fmt.Println("    fpath+=(" + zshDir + ")")
		fmt.Println("    autoload -Uz compinit && compinit")
		fmt.Println()
		fmt.Println(styles.Faint.Render("  Then restart your terminal or run:"))
		fmt.Println(
			styles.Faint.Render("    rm -f ~/.zcompdump; source ~/.zshrc"),
		)
	case "fish":
		fmt.Println(
			styles.Faint.Render(
				"  Fish loads completions automatically from this path.",
			),
		)
		fmt.Println(
			styles.Faint.Render(
				"  Open a new terminal or run: source " + destPath,
			),
		)
	}
}

// completionInstallPath returns the standard persistent path for shell completion scripts.
func completionInstallPath(shell string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	switch shell {
	case "bash":
		// ~/.local/share/bash-completion/completions/pam
		return filepath.Join(
			home,
			".local",
			"share",
			"bash-completion",
			"completions",
			"pam",
		)
	case "zsh":
		// ~/.zsh/completions/_pam  (user adds fpath+=~/.zsh/completions to .zshrc)
		return filepath.Join(home, ".zsh", "completions", "_pam")
	case "fish":
		// ~/.config/fish/completions/pam.fish
		return filepath.Join(home, ".config", "fish", "completions", "pam.fish")
	default:
		return ""
	}
}

// handleCompleteTables is an internal hidden command called by shell completion scripts.
// It prints one table name per line so the shell can offer them as completions.
func (a *App) handleCompleteTables() {
	tables := a.fetchTableNames()
	for _, t := range tables {
		fmt.Println(t)
	}
}

// handleCompleteConnections is an internal hidden command called by shell completion scripts.
func (a *App) handleCompleteConnections() {
	if a.config.Connections == nil {
		return
	}
	for name := range a.config.Connections {
		fmt.Println(name)
	}
}

// handleCompleteQueries is an internal hidden command called by shell completion scripts.
func (a *App) handleCompleteQueries() {
	if a.config.CurrentConnection == "" {
		return
	}
	connYAML, ok := a.config.Connections[a.config.CurrentConnection]
	if !ok || connYAML.Queries == nil {
		return
	}
	for name := range connYAML.Queries {
		fmt.Println(name)
	}
}

func (a *App) fetchTableNames() []string {
	if a.config.CurrentConnection == "" {
		return nil
	}

	connYAML, ok := a.config.Connections[a.config.CurrentConnection]
	if !ok {
		return nil
	}

	conn := config.FromConnectionYaml(connYAML)
	if err := conn.Open(); err != nil {
		return nil
	}
	defer conn.Close()

	tables, err := conn.GetTables()
	if err != nil {
		return nil
	}

	return tables
}

const bashCompletionScript = `# pam bash completion
# Install: pam completion bash --install
#   or:    eval "$(pam completion bash)"
#   or:    pam completion bash > ~/.local/share/bash-completion/completions/pam

_pam_completions() {
    local cur prev words cword
    _init_completion -n = || return

    local commands="init switch use add save remove delete query run list ls edit info status test tables t explore table-view tv disconnect clear unset explain history help completion"

    # If completing the first argument (the command)
    if [[ ${cword} -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
        return
    fi

    local cmd="${words[1]}"

    case "${cmd}" in
        switch|use)
            local connections
            connections=$(pam __complete_connections 2>/dev/null)
            COMPREPLY=( $(compgen -W "${connections}" -- "${cur}") )
            return
            ;;
        run)
            # Complete query names for run
            if [[ ${cword} -eq 2 ]] && [[ "${cur}" != -* ]]; then
                local queries
                queries=$(pam __complete_queries 2>/dev/null)
                COMPREPLY=( $(compgen -W "${queries}" -- "${cur}") )
                return
            fi
            COMPREPLY=( $(compgen -W "--edit -e --last -l" -- "${cur}") )
            return
            ;;
        remove|delete)
            if [[ ${cword} -eq 2 ]]; then
                local queries
                queries=$(pam __complete_queries 2>/dev/null)
                COMPREPLY=( $(compgen -W "${queries}" -- "${cur}") )
                return
            fi
            ;;
        tables|t|explore)
            if [[ ${cword} -eq 2 ]] && [[ "${cur}" != -* ]]; then
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                COMPREPLY=( $(compgen -W "${tables}" -- "${cur}") )
                return
            fi
            COMPREPLY=( $(compgen -W "--oneline -o" -- "${cur}") )
            return
            ;;
        table-view|tv)
            if [[ ${cword} -eq 2 ]] && [[ "${cur}" != -* ]]; then
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                COMPREPLY=( $(compgen -W "${tables}" -- "${cur}") )
                return
            fi
            ;;
        query)
            # Complete --table= with table names
            if [[ "${cur}" == --table=* ]]; then
                local prefix="${cur%%=*}="
                local typed="${cur#*=}"
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                local completions=()
                while IFS= read -r t; do
                    [[ -n "${t}" ]] && completions+=("${prefix}${t}")
                done <<< "${tables}"
                COMPREPLY=( $(compgen -W "${completions[*]}" -- "${cur}") )
                compopt -o nospace
                return
            fi
            if [[ "${cur}" == -t=* ]]; then
                local prefix="-t="
                local typed="${cur#*=}"
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                local completions=()
                while IFS= read -r t; do
                    [[ -n "${t}" ]] && completions+=("${prefix}${t}")
                done <<< "${tables}"
                COMPREPLY=( $(compgen -W "${completions[*]}" -- "${cur}") )
                compopt -o nospace
                return
            fi
            # After --table or -t flag, complete with table names
            if [[ "${prev}" == "--table" ]] || [[ "${prev}" == "-t" ]]; then
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                COMPREPLY=( $(compgen -W "${tables}" -- "${cur}") )
                return
            fi
            COMPREPLY=( $(compgen -W "--table -t --edit -e" -- "${cur}") )
            return
            ;;
        explain)
            if [[ ${cword} -eq 2 ]] && [[ "${cur}" != -* ]]; then
                local tables
                tables=$(pam __complete_tables 2>/dev/null)
                COMPREPLY=( $(compgen -W "${tables}" -- "${cur}") )
                return
            fi
            COMPREPLY=( $(compgen -W "--depth -d" -- "${cur}") )
            return
            ;;
        list)
            if [[ ${cword} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "connections queries" -- "${cur}") )
                return
            fi
            ;;
        edit)
            if [[ ${cword} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "config queries" -- "${cur}") )
                return
            fi
            ;;
        info)
            if [[ ${cword} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "tables views" -- "${cur}") )
                return
            fi
            ;;
        help)
            if [[ ${cword} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
                return
            fi
            ;;
        completion)
            if [[ ${cword} -eq 2 ]]; then
                COMPREPLY=( $(compgen -W "bash zsh fish" -- "${cur}") )
                return
            fi
            ;;
    esac
}

complete -F _pam_completions pam
`

const zshCompletionScript = `#compdef pam
# pam zsh completion
# Install: pam completion zsh --install
#   or:    eval "$(pam completion zsh)"
#   or:    pam completion zsh > "${fpath[1]}/_pam"

_pam_tables() {
    local -a tables
    tables=(${(f)"$(pam __complete_tables 2>/dev/null)"})
    _describe 'table' tables
}

_pam_connections() {
    local -a conns
    conns=(${(f)"$(pam __complete_connections 2>/dev/null)"})
    _describe 'connection' conns
}

_pam_queries() {
    local -a queries
    queries=(${(f)"$(pam __complete_queries 2>/dev/null)"})
    _describe 'query' queries
}

_pam() {
    local -a commands
    commands=(
        'init:Create a new database connection'
        'switch:Switch active connection'
        'use:Switch active connection'
        'add:Save a new named query'
        'save:Save a new named query'
        'remove:Remove a saved query'
        'delete:Remove a saved query'
        'query:Run a SQL query against a table'
        'run:Execute a saved query'
        'list:List connections or queries'
        'ls:List connections'
        'edit:Open config or queries in editor'
        'info:Show tables or views'
        'status:Show current connection'
        'test:Test current connection'
        'tables:List or query database tables'
        't:List or query database tables'
        'explore:Explore database schema'
        'table-view:View and edit table structure'
        'tv:View and edit table structure'
        'disconnect:Disconnect from current database'
        'clear:Disconnect from current database'
        'unset:Disconnect from current database'
        'explain:Visualize table relationships'
        'history:Show query history'
        'help:Show help'
        'completion:Generate shell completion script'
    )

    if (( CURRENT == 2 )); then
        _describe 'command' commands
        return
    fi

    local cmd="${words[2]}"

    case "${cmd}" in
        switch|use)
            _pam_connections
            ;;
        run)
            if (( CURRENT == 3 )); then
                _pam_queries
            else
                _arguments \
                    '--edit[Edit query before running]' \
                    '-e[Edit query before running]' \
                    '--last[Run last query]' \
                    '-l[Run last query]'
            fi
            ;;
        remove|delete)
            if (( CURRENT == 3 )); then
                _pam_queries
            fi
            ;;
        tables|t|explore)
            if (( CURRENT == 3 )); then
                _pam_tables
            else
                _arguments \
                    '--oneline[Display one table per line]' \
                    '-o[Display one table per line]'
            fi
            ;;
        table-view|tv)
            if (( CURRENT == 3 )); then
                _pam_tables
            fi
            ;;
        query)
            _arguments \
                '--table=[Target table name]:table:_pam_tables' \
                '-t=[Target table name]:table:_pam_tables' \
                '--edit[Edit SQL in editor]' \
                '-e[Edit SQL in editor]'
            ;;
        explain)
            if (( CURRENT == 3 )); then
                _pam_tables
            else
                _arguments \
                    '--depth[Relationship depth]:depth:' \
                    '-d[Relationship depth]:depth:'
            fi
            ;;
        list)
            if (( CURRENT == 3 )); then
                local -a sub
                sub=('connections:List all connections' 'queries:List saved queries')
                _describe 'subcommand' sub
            fi
            ;;
        edit)
            if (( CURRENT == 3 )); then
                local -a sub
                sub=('config:Edit configuration' 'queries:Edit queries')
                _describe 'subcommand' sub
            fi
            ;;
        info)
            if (( CURRENT == 3 )); then
                local -a sub
                sub=('tables:Show tables' 'views:Show views')
                _describe 'subcommand' sub
            fi
            ;;
        help)
            if (( CURRENT == 3 )); then
                _describe 'command' commands
            fi
            ;;
        completion)
            if (( CURRENT == 3 )); then
                local -a shells
                shells=('bash' 'zsh' 'fish')
                _describe 'shell' shells
            fi
            ;;
    esac
}

_pam "$@"
`

const fishCompletionScript = `# pam fish completion
# Install: pam completion fish --install
#   or:    pam completion fish | source
#   or:    pam completion fish > ~/.config/fish/completions/pam.fish

# Disable file completions for pam
complete -c pam -f

# Helper functions
function __pam_tables
    pam __complete_tables 2>/dev/null
end

function __pam_connections
    pam __complete_connections 2>/dev/null
end

function __pam_queries
    pam __complete_queries 2>/dev/null
end

function __pam_no_subcommand
    set -l cmd (commandline -opc)
    test (count $cmd) -eq 1
end

function __pam_using_command
    set -l cmd (commandline -opc)
    test (count $cmd) -gt 1; and test "$cmd[2]" = "$argv[1]"
end

function __pam_using_any_command
    set -l cmd (commandline -opc)
    if test (count $cmd) -le 1
        return 1
    end
    for c in $argv
        if test "$cmd[2]" = "$c"
            return 0
        end
    end
    return 1
end

function __pam_needs_table_after_flag
    set -l cmd (commandline -opc)
    set -l last $cmd[-1]
    test "$last" = "--table" -o "$last" = "-t"
end

# Subcommands
complete -c pam -n __pam_no_subcommand -a init -d 'Create a new database connection'
complete -c pam -n __pam_no_subcommand -a switch -d 'Switch active connection'
complete -c pam -n __pam_no_subcommand -a use -d 'Switch active connection'
complete -c pam -n __pam_no_subcommand -a add -d 'Save a new named query'
complete -c pam -n __pam_no_subcommand -a save -d 'Save a new named query'
complete -c pam -n __pam_no_subcommand -a remove -d 'Remove a saved query'
complete -c pam -n __pam_no_subcommand -a delete -d 'Remove a saved query'
complete -c pam -n __pam_no_subcommand -a query -d 'Run a SQL query against a table'
complete -c pam -n __pam_no_subcommand -a run -d 'Execute a saved query'
complete -c pam -n __pam_no_subcommand -a list -d 'List connections or queries'
complete -c pam -n __pam_no_subcommand -a ls -d 'List connections'
complete -c pam -n __pam_no_subcommand -a edit -d 'Open config or queries in editor'
complete -c pam -n __pam_no_subcommand -a info -d 'Show tables or views'
complete -c pam -n __pam_no_subcommand -a status -d 'Show current connection'
complete -c pam -n __pam_no_subcommand -a test -d 'Test current connection'
complete -c pam -n __pam_no_subcommand -a tables -d 'List or query database tables'
complete -c pam -n __pam_no_subcommand -a t -d 'List or query database tables'
complete -c pam -n __pam_no_subcommand -a explore -d 'Explore database schema'
complete -c pam -n __pam_no_subcommand -a table-view -d 'View and edit table structure'
complete -c pam -n __pam_no_subcommand -a tv -d 'View and edit table structure'
complete -c pam -n __pam_no_subcommand -a disconnect -d 'Disconnect from database'
complete -c pam -n __pam_no_subcommand -a clear -d 'Disconnect from database'
complete -c pam -n __pam_no_subcommand -a unset -d 'Disconnect from database'
complete -c pam -n __pam_no_subcommand -a explain -d 'Visualize table relationships'
complete -c pam -n __pam_no_subcommand -a history -d 'Show query history'
complete -c pam -n __pam_no_subcommand -a help -d 'Show help'
complete -c pam -n __pam_no_subcommand -a completion -d 'Generate shell completion'

# switch / use — complete connection names
complete -c pam -n '__pam_using_any_command switch use' -a '(__pam_connections)'

# run — complete query names and flags
complete -c pam -n '__pam_using_command run' -a '(__pam_queries)'
complete -c pam -n '__pam_using_command run' -l edit -s e -d 'Edit query before running'
complete -c pam -n '__pam_using_command run' -l last -s l -d 'Run last query'

# remove / delete — complete query names
complete -c pam -n '__pam_using_any_command remove delete' -a '(__pam_queries)'

# tables / t / explore — complete table names
complete -c pam -n '__pam_using_any_command tables t explore' -a '(__pam_tables)'
complete -c pam -n '__pam_using_any_command tables t' -l oneline -s o -d 'One table per line'

# table-view / tv — complete table names
complete -c pam -n '__pam_using_any_command table-view tv' -a '(__pam_tables)'

# query — flags and table completion
complete -c pam -n '__pam_using_command query; and __pam_needs_table_after_flag' -a '(__pam_tables)'
complete -c pam -n '__pam_using_command query' -l table -s t -d 'Target table' -r -a '(__pam_tables)'
complete -c pam -n '__pam_using_command query' -l edit -s e -d 'Edit SQL in editor'

# explain — complete table names
complete -c pam -n '__pam_using_command explain' -a '(__pam_tables)'
complete -c pam -n '__pam_using_command explain' -l depth -s d -d 'Relationship depth' -r

# list subcommands
complete -c pam -n '__pam_using_command list' -a 'connections queries'

# edit subcommands
complete -c pam -n '__pam_using_command edit' -a 'config queries'

# info subcommands
complete -c pam -n '__pam_using_command info' -a 'tables views'

# help — complete command names
complete -c pam -n '__pam_using_command help' -a 'init switch use add save remove delete query run list ls edit info status test tables t explore table-view tv disconnect clear unset explain history help completion'

# completion shells
complete -c pam -n '__pam_using_command completion' -a 'bash zsh fish'
`
