package initui

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/caiolandgraf/pam/internal/db"
	"github.com/caiolandgraf/pam/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Field indices
const (
	fieldName     = iota // Connection alias
	fieldType            // DB type selector
	fieldHost            // host / file path for sqlite
	fieldPort            // port
	fieldUser            // username
	fieldPassword        // password (masked)
	fieldDatabase        // database / schema name
	fieldTotal           // sentinel
)

// dbDefaults holds the default port for each DB type
var dbDefaults = map[string]string{
	"postgres":   "5432",
	"mysql":      "3306",
	"sqlserver":  "1433",
	"clickhouse": "9000",
	"oracle":     "1521",
	"firebird":   "3050",
	"sqlite":     "",
}

// fieldsForType returns which field indices are visible for a given db type.
// SQLite only needs name + type + host (file path).
func fieldsForType(dbType string) []int {
	switch dbType {
	case "sqlite":
		return []int{fieldName, fieldType, fieldHost}
	default:
		return []int{
			fieldName,
			fieldType,
			fieldHost,
			fieldPort,
			fieldUser,
			fieldPassword,
			fieldDatabase,
		}
	}
}

type fieldState struct {
	value  string
	cursor int
}

type InitInputModel struct {
	fields      [fieldTotal]fieldState
	dbTypes     []string
	cursorIndex int // index within the *visible* field list
	aborted     bool
	showPass    bool // toggle password visibility
}

func NewInitInputModel(name, dbType, connString string) InitInputModel {
	types := db.GetSupportedDBTypes()

	if dbType == "" && connString != "" {
		dbType = db.InferDBType(connString)
	}

	m := InitInputModel{
		dbTypes:  types,
		showPass: false,
	}

	m.fields[fieldName] = fieldState{value: name, cursor: len(name)}

	if dbType == "" {
		dbType = types[0]
	}
	m.fields[fieldType] = fieldState{value: dbType}

	// Pre-fill default port
	if port, ok := dbDefaults[dbType]; ok {
		m.fields[fieldPort] = fieldState{value: port, cursor: len(port)}
	}

	// If a full connString was provided, try to parse it into individual fields
	if connString != "" {
		m.parseConnString(connString)
	}

	return m
}

// parseConnString tries to decompose an existing connection string into fields.
func (m *InitInputModel) parseConnString(connString string) {
	u, err := url.Parse(connString)
	if err != nil {
		return
	}
	if host := u.Hostname(); host != "" {
		m.fields[fieldHost] = fieldState{value: host, cursor: len(host)}
	}
	if port := u.Port(); port != "" {
		m.fields[fieldPort] = fieldState{value: port, cursor: len(port)}
	}
	if u.User != nil {
		user := u.User.Username()
		m.fields[fieldUser] = fieldState{value: user, cursor: len(user)}
		if pass, ok := u.User.Password(); ok {
			m.fields[fieldPassword] = fieldState{value: pass, cursor: len(pass)}
		}
	}
	if dbName := strings.TrimPrefix(u.Path, "/"); dbName != "" {
		m.fields[fieldDatabase] = fieldState{value: dbName, cursor: len(dbName)}
	}
}

// buildConnString assembles a connection string from the individual fields.
func (m *InitInputModel) buildConnString() string {
	dbType := m.fields[fieldType].value
	host := m.fields[fieldHost].value
	port := m.fields[fieldPort].value
	user := m.fields[fieldUser].value
	pass := m.fields[fieldPassword].value
	database := m.fields[fieldDatabase].value

	if host == "" {
		return ""
	}

	switch dbType {
	case "sqlite":
		return host

	case "postgres":
		u := &url.URL{Scheme: "postgres"}
		if user != "" || pass != "" {
			u.User = url.UserPassword(user, pass)
		}
		if port != "" {
			u.Host = host + ":" + port
		} else {
			u.Host = host
		}
		u.Path = "/" + database
		return u.String()

	case "mysql":
		// DSN format: user:pass@tcp(host:port)/dbname
		auth := ""
		if user != "" || pass != "" {
			auth = url.QueryEscape(user) + ":" + url.QueryEscape(pass) + "@"
		}
		hp := host
		if port != "" {
			hp = host + ":" + port
		}
		return fmt.Sprintf("%stcp(%s)/%s", auth, hp, database)

	case "sqlserver":
		u := &url.URL{Scheme: "sqlserver"}
		if user != "" || pass != "" {
			u.User = url.UserPassword(user, pass)
		}
		if port != "" {
			u.Host = host + ":" + port
		} else {
			u.Host = host
		}
		if database != "" {
			q := u.Query()
			q.Set("database", database)
			u.RawQuery = q.Encode()
		}
		return u.String()

	case "clickhouse":
		u := &url.URL{Scheme: "clickhouse"}
		if user != "" || pass != "" {
			u.User = url.UserPassword(user, pass)
		}
		if port != "" {
			u.Host = host + ":" + port
		} else {
			u.Host = host
		}
		u.Path = "/" + database
		return u.String()

	case "oracle":
		// user/pass@host:port/service
		return fmt.Sprintf("%s/%s@%s:%s/%s", user, pass, host, port, database)

	case "firebird":
		return fmt.Sprintf("%s:%s@%s:%s/%s", user, pass, host, port, database)

	default:
		u := &url.URL{Scheme: dbType}
		if user != "" || pass != "" {
			u.User = url.UserPassword(user, pass)
		}
		if port != "" {
			u.Host = host + ":" + port
		} else {
			u.Host = host
		}
		u.Path = "/" + database
		return u.String()
	}
}

func (m InitInputModel) Init() tea.Cmd {
	return nil
}

// visibleFields returns the ordered list of field indices for the current db type.
func (m *InitInputModel) visibleFields() []int {
	return fieldsForType(m.fields[fieldType].value)
}

// currentField returns the field index at the current cursor position.
func (m *InitInputModel) currentField() int {
	vf := m.visibleFields()
	if m.cursorIndex >= len(vf) {
		return vf[len(vf)-1]
	}
	return vf[m.cursorIndex]
}

func (m InitInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			m.aborted = true
			return m, tea.Quit

		case "enter":
			if m.isComplete() {
				return m, tea.Quit
			}
			m.moveToNextEmpty()
			return m, nil

		case "tab", "down":
			vf := m.visibleFields()
			if m.cursorIndex < len(vf)-1 {
				m.cursorIndex++
			}

		case "shift+tab", "up":
			if m.cursorIndex > 0 {
				m.cursorIndex--
			}

		case "right":
			cf := m.currentField()
			if cf == fieldType {
				m.cycleDbType(1)
			} else {
				fs := &m.fields[cf]
				if fs.cursor < len(fs.value) {
					fs.cursor++
				}
			}

		case "left":
			cf := m.currentField()
			if cf == fieldType {
				m.cycleDbType(-1)
			} else {
				fs := &m.fields[cf]
				if fs.cursor > 0 {
					fs.cursor--
				}
			}

		case "ctrl+a":
			cf := m.currentField()
			if cf != fieldType {
				m.fields[cf].cursor = 0
			}

		case "ctrl+e":
			cf := m.currentField()
			if cf != fieldType {
				fs := &m.fields[cf]
				fs.cursor = len(fs.value)
			}

		case "ctrl+p":
			m.showPass = !m.showPass

		case "backspace":
			m.handleBackspace()

		default:
			m.handleInput(msg.String())
		}
	}

	return m, nil
}

func (m *InitInputModel) handleBackspace() {
	cf := m.currentField()
	if cf == fieldType {
		return
	}
	fs := &m.fields[cf]
	if fs.cursor > 0 {
		fs.value = fs.value[:fs.cursor-1] + fs.value[fs.cursor:]
		fs.cursor--
	}
}

func (m *InitInputModel) handleInput(ch string) {
	if ch == "" {
		return
	}
	// Strip bracketed-paste wrappers that some terminals inject
	if len(ch) > 1 && strings.HasPrefix(ch, "[") && strings.HasSuffix(ch, "]") {
		ch = ch[1 : len(ch)-1]
	}

	cf := m.currentField()
	if cf == fieldType {
		// Number shortcut to pick db type
		if len(ch) == 1 && ch >= "1" && ch <= "9" {
			idx := int(ch[0] - '1')
			if idx < len(m.dbTypes) {
				m.setDbType(m.dbTypes[idx])
			}
		}
		return
	}

	fs := &m.fields[cf]
	fs.value = fs.value[:fs.cursor] + ch + fs.value[fs.cursor:]
	fs.cursor += len(ch)
}

func (m *InitInputModel) setDbType(t string) {
	m.fields[fieldType].value = t

	// Reset port to the new default only if the current value is still
	// one of the known defaults (i.e. the user hasn't customised it)
	if def, ok := dbDefaults[t]; ok {
		current := m.fields[fieldPort].value
		isDefault := current == ""
		if !isDefault {
			for _, v := range dbDefaults {
				if v == current {
					isDefault = true
					break
				}
			}
		}
		if isDefault {
			m.fields[fieldPort] = fieldState{value: def, cursor: len(def)}
		}
	}

	// Keep cursor index in bounds after the visible field list may have changed
	vf := m.visibleFields()
	if m.cursorIndex >= len(vf) {
		m.cursorIndex = len(vf) - 1
	}
}

func (m *InitInputModel) cycleDbType(dir int) {
	types := m.dbTypes
	if len(types) == 0 {
		return
	}
	current := m.fields[fieldType].value
	idx := 0
	for i, t := range types {
		if t == current {
			idx = i
			break
		}
	}
	newIdx := (idx + dir + len(types)) % len(types)
	m.setDbType(types[newIdx])
}

func (m *InitInputModel) isComplete() bool {
	if strings.TrimSpace(m.fields[fieldName].value) == "" {
		return false
	}
	if m.fields[fieldType].value == "" {
		return false
	}
	vf := m.visibleFields()
	for _, f := range vf {
		if f == fieldType || f == fieldPassword {
			continue // type is already checked; password is optional
		}
		if strings.TrimSpace(m.fields[f].value) == "" {
			return false
		}
	}
	return true
}

func (m *InitInputModel) moveToNextEmpty() {
	vf := m.visibleFields()
	for i, f := range vf {
		if f == fieldType || f == fieldPassword {
			continue
		}
		if strings.TrimSpace(m.fields[f].value) == "" {
			m.cursorIndex = i
			return
		}
	}
}

// View renders the TUI.
func (m InitInputModel) View() string {
	var b strings.Builder

	b.WriteString(styles.Title.Render("Initialize new connection"))
	b.WriteString("\n\n")

	vf := m.visibleFields()
	for i, f := range vf {
		focused := i == m.cursorIndex
		switch f {
		case fieldType:
			m.renderTypeDropdown(&b, focused)
		case fieldPassword:
			m.renderPasswordField(&b, focused)
		default:
			label := fieldLabel(f, m.fields[fieldType].value)
			fs := m.fields[f]
			m.renderField(&b, label, fs.value, fs.cursor, focused)
		}
	}

	// Live preview of the assembled connection string
	if preview := m.buildConnString(); preview != "" {
		b.WriteString("\n")
		b.WriteString(styles.Faint.Render("  conn › "))
		b.WriteString(
			lipgloss.NewStyle().
				Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
				Italic(true).
				Render(preview),
		)
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(
		styles.Faint.Render(
			"↑/↓ Tab: navigate  ←/→: cursor / cycle type  Ctrl+P: show/hide password  Enter: confirm  Esc: cancel",
		),
	)

	return b.String()
}

func fieldLabel(f int, dbType string) string {
	switch f {
	case fieldName:
		return "Connection name"
	case fieldHost:
		if dbType == "sqlite" {
			return "File path      "
		}
		return "Host           "
	case fieldPort:
		return "Port           "
	case fieldUser:
		return "Username       "
	case fieldPassword:
		return "Password       "
	case fieldDatabase:
		return "Database       "
	}
	return "Field          "
}

func (m InitInputModel) renderField(
	b *strings.Builder,
	label, value string,
	cursorPos int,
	focused bool,
) {
	if focused {
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
			Bold(true).
			Render("  " + label + " › ")

		before := value[:cursorPos]
		after := value[cursorPos:]
		inputBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(before + "▏" + after)

		b.WriteString(prompt + inputBox + "\n")
	} else {
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render("  " + label + "   ")

		inputBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(value)

		b.WriteString(prompt + inputBox + "\n")
	}
}

func (m InitInputModel) renderPasswordField(b *strings.Builder, focused bool) {
	fs := m.fields[fieldPassword]
	display := fs.value
	cursorPos := fs.cursor
	if !m.showPass && len(display) > 0 {
		display = strings.Repeat("•", len(display))
	}
	label := fieldLabel(fieldPassword, "")
	m.renderField(b, label, display, cursorPos, focused)
}

func (m InitInputModel) renderTypeDropdown(b *strings.Builder, focused bool) {
	const label = "Database type  "
	current := m.fields[fieldType].value
	if current == "" {
		current = "<select>"
	}

	if focused {
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
			Bold(true).
			Render("  " + label + " › ")

		inputBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(current + "  ◀ ▶")

		b.WriteString(prompt + inputBox + "\n")

		// Show numbered list of available types
		b.WriteString(styles.Faint.Render("    "))
		for i, t := range m.dbTypes {
			if i > 0 {
				b.WriteString(styles.Faint.Render("  "))
			}
			if t == m.fields[fieldType].value {
				b.WriteString(
					lipgloss.NewStyle().
						Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
						Bold(true).
						Render(fmt.Sprintf("[%d] %s", i+1, t)),
				)
			} else {
				b.WriteString(styles.Faint.Render(fmt.Sprintf("%d:%s", i+1, t)))
			}
		}
		b.WriteString("\n")
	} else {
		prompt := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render("  " + label + "   ")

		inputBox := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Muted)).
			Render(current)

		b.WriteString(prompt + inputBox + "\n")
	}
}

func (m InitInputModel) GetName() string {
	return m.fields[fieldName].value
}

func (m InitInputModel) GetDBType() string {
	return m.fields[fieldType].value
}

func (m InitInputModel) GetConnString() string {
	return m.buildConnString()
}

func (m InitInputModel) WasAborted() bool {
	return m.aborted
}

func CollectInitParameters(
	name, dbType, connString string,
) (string, string, string, error) {
	model := NewInitInputModel(name, dbType, connString)
	program := tea.NewProgram(model)

	finalModel, err := program.Run()
	if err != nil {
		return "", "", "", err
	}

	inputModel := finalModel.(InitInputModel)
	if inputModel.WasAborted() {
		return "", "", "", ErrAborted
	}

	return inputModel.GetName(), inputModel.GetDBType(), inputModel.GetConnString(), nil
}

var ErrAborted = fmt.Errorf("init input aborted")
