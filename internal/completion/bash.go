package completion

func GenerateBash() string {
	return `# Bash completion for pam
_pam_complete() {
    local cur prev words cword
    _init_completion || return

    # Call pam __complete with current arguments
    local completions
    completions=$(pam __complete "${words[@]:1}")

    # Filter completions based on current word
    COMPREPLY=($(compgen -W "$completions" -- "$cur"))
}

complete -F _pam_complete pam
`
}
