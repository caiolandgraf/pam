# Shell Completion

PAM provides intelligent tab completion for bash, zsh, and fish. Completions are **dynamic** - they automatically include your saved queries and connections.

## Installation

The quickest way to install completions to the standard path for your shell:

```bash
pam completion --install
```

This automatically detects your shell and writes the completion file to the appropriate location (`~/.config/fish/completions/pam.fish`, `~/.config/pam/pam.bash`, etc.).

Alternatively, you can install manually:

<details>
<summary>Bash</summary>

```bash
# Temporary (current session)
source <(pam completion bash)

# Permanent (add to ~/.bashrc)
echo 'eval "$(pam completion bash)"' >> ~/.bashrc
```
</details>

<details>
<summary>Zsh</summary>

```bash
# Temporary (current session)
autoload -U compinit && compinit
source <(pam completion zsh)

# Permanent (add to ~/.zshrc)
echo 'autoload -U compinit && compinit' >> ~/.zshrc
echo 'eval "$(pam completion zsh)"' >> ~/.zshrc
```
</details>

<details>
<summary>Fish</summary>

```bash
# Fish loads completions automatically from ~/.config/fish/completions/
pam completion fish > ~/.config/fish/completions/pam.fish
# Restart your shell or run: exec fish
```
</details>

## Usage

After installation, press TAB to autocomplete:

```bash
pam [TAB]              # List all commands
pam run [TAB]          # List queries from current connection
pam switch [TAB]       # List connection names
pam info [TAB]         # List: table, view
pam list [TAB]         # List: queries
pam edit [TAB]         # List queries to edit
```

**Note:** Query completion only shows queries from your current connection. Use `pam switch <connection>` to change connections.

> Shell completion is included in all release builds. Use `pam completion --install` for the fastest setup.
