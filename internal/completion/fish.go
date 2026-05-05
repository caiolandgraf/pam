package completion

func GenerateFish() string {
	return `# Fish completion for pam
complete -c pam -f -a "(pam __complete (commandline -opc)[2..-1])"
`
}
