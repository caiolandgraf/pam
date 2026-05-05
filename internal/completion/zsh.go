package completion

func GenerateZsh() string {
	return `#compdef pam

_pam() {
    local -a completions
    completions=(${(f)"$(pam __complete $words[2,-1])"})

    if [[ -n "$completions" ]]; then
        _describe 'pam' completions
    else
        _files  # Fallback to files
    fi
}

# Register completion function (works when sourced in .zshrc or installed as file)
compdef _pam pam
`
}
