package configui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/caiolandgraf/pam/internal/config"
	"github.com/caiolandgraf/pam/internal/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ── Tabs ──────────────────────────────────────────────────────────────────────

const (
	tabGeneral = iota
	tabUI
	tabTheme
	tabCount
)

var tabLabels = [tabCount]string{"General", "UI Display", "Theme"}

// ── Available colour schemes (must match styles.GetScheme keys) ───────────────

var schemeNames = []string{
	"default", "dracula", "gruvbox", "solarized", "nord", "monokai",
	"black-metal", "black-metal-gorgoroth", "vesper",
	"catppuccin-mocha", "tokyo-night", "rose-pine", "terracotta",
}

// ── ui-visibility fields metadata ─────────────────────────────────────────────

type boolField struct {
	label string
	get   func(*config.UIVisibility) bool
	set   func(*config.UIVisibility, bool)
}

var uiFields = []boolField{
	{
		"Show query name",
		func(u *config.UIVisibility) bool { return u.QueryName },
		func(u *config.UIVisibility, v bool) { u.QueryName = v },
	},
	{
		"Show query SQL",
		func(u *config.UIVisibility) bool { return u.QuerySQL },
		func(u *config.UIVisibility, v bool) { u.QuerySQL = v },
	},
	{
		"Show type icons",
		func(u *config.UIVisibility) bool { return u.TypeDisplay },
		func(u *config.UIVisibility, v bool) { u.TypeDisplay = v },
	},
	{
		"Show key icons",
		func(u *config.UIVisibility) bool { return u.KeyIcons },
		func(u *config.UIVisibility, v bool) { u.KeyIcons = v },
	},
	{
		"Show cell preview",
		func(u *config.UIVisibility) bool { return u.FooterCellContent },
		func(u *config.UIVisibility, v bool) { u.FooterCellContent = v },
	},
	{
		"Show footer stats",
		func(u *config.UIVisibility) bool { return u.FooterStats },
		func(u *config.UIVisibility, v bool) { u.FooterStats = v },
	},
	{
		"Show footer keymaps",
		func(u *config.UIVisibility) bool { return u.FooterKeymaps },
		func(u *config.UIVisibility, v bool) { u.FooterKeymaps = v },
	},
}

// ── General tab fields ────────────────────────────────────────────────────────

const (
	genConn = iota
	genRowLimit
	genColWidth
	genHistSize
	genCount
)

var genLabels = [genCount]string{
	"Current connection",
	"Default row limit",
	"Default column width",
	"History size",
}

// ── Model ─────────────────────────────────────────────────────────────────────

type Model struct {
	cfg      *config.Config
	tab      int
	cursor   int
	width    int
	saved    bool
	quitting bool

	// General
	connNames []string
	connIndex int
	rowLimit  int
	colWidth  int
	histSize  int

	// Theme
	schemeIndex int

	// UI (live copy of the struct)
	ui config.UIVisibility
}

// New initialises the model from an existing Config.
func New(cfg *config.Config) Model {
	// collect & sort connection names
	names := make([]string, 0, len(cfg.Connections))
	for k := range cfg.Connections {
		names = append(names, k)
	}
	sort.Strings(names)

	// find current connection index
	connIdx := 0
	for i, n := range names {
		if n == cfg.CurrentConnection {
			connIdx = i
			break
		}
	}

	// find current scheme index
	schemeIdx := 0
	for i, s := range schemeNames {
		if s == cfg.ColorScheme {
			schemeIdx = i
			break
		}
	}

	histSize := cfg.History.Size
	if histSize == 0 {
		histSize = 100
	}

	return Model{
		cfg:         cfg,
		connNames:   names,
		connIndex:   connIdx,
		rowLimit:    cfg.DefaultRowLimit,
		colWidth:    cfg.DefaultColumnWidth,
		histSize:    histSize,
		schemeIndex: schemeIdx,
		ui:          cfg.UIVisibility,
	}
}

// ── tea.Model interface ───────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		// quit without saving
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		// save & quit
		case "ctrl+s", "s":
			m.applyToConfig()
			if err := m.cfg.Save(); err == nil {
				m.saved = true
			}
			return m, tea.Quit

		// switch tab
		case "tab", "shift+tab":
			if msg.String() == "tab" {
				m.tab = (m.tab + 1) % tabCount
			} else {
				m.tab = (m.tab - 1 + tabCount) % tabCount
			}
			m.cursor = 0
			return m, nil

		// navigate within tab
		case "up", "k":
			m.moveCursor(-1)
			return m, nil
		case "down", "j":
			m.moveCursor(1)
			return m, nil

		// change value
		case "left", "h", "-":
			m.changeValue(-1)
			return m, nil
		case "right", "l", "+", " ", "enter":
			m.changeValue(+1)
			return m, nil
		}
	}
	return m, nil
}

func (m *Model) moveCursor(delta int) {
	switch m.tab {
	case tabGeneral:
		m.cursor = clamp(m.cursor+delta, 0, genCount-1)
	case tabUI:
		m.cursor = clamp(m.cursor+delta, 0, len(uiFields)-1)
	case tabTheme:
		m.cursor = 0
	}
}

func (m *Model) changeValue(delta int) {
	switch m.tab {

	case tabGeneral:
		switch m.cursor {
		case genConn:
			if len(m.connNames) == 0 {
				return
			}
			m.connIndex = mod(m.connIndex+delta, len(m.connNames))
		case genRowLimit:
			m.rowLimit = clamp(m.rowLimit+delta*100, 10, 100_000)
		case genColWidth:
			m.colWidth = clamp(m.colWidth+delta, 5, 80)
		case genHistSize:
			m.histSize = clamp(m.histSize+delta*10, 0, 10_000)
		}

	case tabUI:
		f := uiFields[m.cursor]
		cur := f.get(&m.ui)
		f.set(&m.ui, !cur)

	case tabTheme:
		m.schemeIndex = mod(m.schemeIndex+delta, len(schemeNames))
	}
}

func (m *Model) applyToConfig() {
	if len(m.connNames) > 0 {
		m.cfg.CurrentConnection = m.connNames[m.connIndex]
	}
	m.cfg.DefaultRowLimit = m.rowLimit
	m.cfg.DefaultColumnWidth = m.colWidth
	m.cfg.History.Size = m.histSize
	m.cfg.ColorScheme = schemeNames[m.schemeIndex]
	m.cfg.UIVisibility = m.ui
}

// ── View ──────────────────────────────────────────────────────────────────────

func (m Model) View() string {
	if m.quitting && !m.saved {
		return styles.Faint.Render("  Config unchanged.\n")
	}
	if m.saved {
		return styles.Success.Render("  ✓ Configuration saved.\n")
	}

	w := m.width
	if w < 60 {
		w = 72
	}

	var b strings.Builder

	// ── header ──
	b.WriteString("\n")
	b.WriteString(styles.Title.Render("  🗂  PAM Configuration"))
	b.WriteString("\n")
	b.WriteString(styles.Separator.Render(strings.Repeat("─", min(w-2, 68))))
	b.WriteString("\n\n")

	// ── tabs ──
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// ── content ──
	switch m.tab {
	case tabGeneral:
		b.WriteString(m.renderGeneral())
	case tabUI:
		b.WriteString(m.renderUI())
	case tabTheme:
		b.WriteString(m.renderTheme())
	}

	// ── footer ──
	b.WriteString("\n")
	b.WriteString(m.renderFooter())
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderTabs() string {
	var parts []string
	for i, label := range tabLabels {
		if i == m.tab {
			parts = append(parts, lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
				Underline(true).
				Render("[ "+label+" ]"))
		} else {
			parts = append(parts, styles.Faint.Render("[ "+label+" ]"))
		}
	}
	return "  " + strings.Join(parts, "  ")
}

// ── General tab ───────────────────────────────────────────────────────────────

func (m Model) renderGeneral() string {
	var b strings.Builder

	accent := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent))
	muted := styles.Faint

	rows := []struct {
		label string
		value string
		hint  string
	}{
		{
			genLabels[genConn],
			connValue(m.connNames, m.connIndex),
			"← → to cycle",
		},
		{
			genLabels[genRowLimit],
			fmt.Sprintf("%d", m.rowLimit),
			"← / → in steps of 100",
		},
		{
			genLabels[genColWidth],
			fmt.Sprintf("%d chars", m.colWidth),
			"← / →",
		},
		{
			genLabels[genHistSize],
			fmt.Sprintf("%d", m.histSize),
			"← / → in steps of 10",
		},
	}

	for i, row := range rows {
		cursor := "  "
		labelStyle := muted
		valStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(styles.ActiveScheme.Normal))

		if i == m.cursor {
			cursor = accent.Render("▸ ")
			labelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
				Bold(true)
			valStyle = accent
		}

		label := labelStyle.Render(fmt.Sprintf("%-24s", row.label))
		val := valStyle.Render(row.value)
		hint := muted.Render("  " + row.hint)

		b.WriteString(cursor + label + " " + val + hint + "\n")
	}

	return b.String()
}

func connValue(names []string, idx int) string {
	if len(names) == 0 {
		return "(no connections)"
	}
	return names[idx]
}

// ── UI Display tab ────────────────────────────────────────────────────────────

func (m Model) renderUI() string {
	var b strings.Builder

	accent := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent))
	muted := styles.Faint
	on := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Success)).
		Bold(true)
	off := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Error))

	for i, f := range uiFields {
		cursor := "  "
		labelStyle := muted
		if i == m.cursor {
			cursor = accent.Render("▸ ")
			labelStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(styles.ActiveScheme.Primary)).
				Bold(true)
		}

		val := f.get(&m.ui)
		var toggle string
		if val {
			toggle = on.Render("  ✓  on ")
		} else {
			toggle = off.Render("  ✗  off")
		}

		b.WriteString(
			cursor + labelStyle.Render(
				fmt.Sprintf("%-26s", f.label),
			) + toggle + "\n",
		)
	}

	b.WriteString("\n")
	b.WriteString(muted.Render("  Space or Enter to toggle"))
	b.WriteString("\n")
	return b.String()
}

// ── Theme tab ─────────────────────────────────────────────────────────────────

func (m Model) renderTheme() string {
	var b strings.Builder

	scheme := styles.GetScheme(schemeNames[m.schemeIndex])
	accent := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ActiveScheme.Accent))
	muted := styles.Faint

	// scheme selector
	prev := muted.Render("← ")
	name := accent.Render(schemeNames[m.schemeIndex])
	next := muted.Render(" →")
	b.WriteString(
		"  " + styles.Faint.Render(
			fmt.Sprintf("%-14s", "Color scheme"),
		) + "  " + prev + name + next + "\n\n",
	)

	// colour preview swatches using the selected scheme
	swatches := []struct {
		label string
		color string
	}{
		{"Primary", scheme.Primary},
		{"Accent ", scheme.Accent},
		{"Success", scheme.Success},
		{"Error  ", scheme.Error},
		{"Muted  ", scheme.Muted},
		{"Normal ", scheme.Normal},
	}

	b.WriteString("  Preview:\n\n")
	for i := 0; i < len(swatches); i += 2 {
		left := swatches[i]
		swatch := func(label, color string) string {
			block := lipgloss.NewStyle().
				Background(lipgloss.Color(color)).
				Foreground(lipgloss.Color(color)).
				Render("  ██  ")
			hex := muted.Render(color)
			return fmt.Sprintf("    %s %s %s", muted.Render(label), block, hex)
		}
		b.WriteString(swatch(left.label, left.color))
		if i+1 < len(swatches) {
			right := swatches[i+1]
			b.WriteString("    " + swatch(right.label, right.color))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// ── Footer ────────────────────────────────────────────────────────────────────

func (m Model) renderFooter() string {
	key := func(k, desc string) string {
		return styles.TableHeader.Render(k) + styles.Faint.Render(desc)
	}
	parts := []string{
		key("ctrl+s", "/s save"),
		key("q", "/esc quit"),
		key("tab", " next tab"),
		key("↑↓", " field"),
		key("←→", " value"),
	}
	return styles.Separator.Render(strings.Repeat("─", 48)) + "\n  " +
		strings.Join(parts, "  ") + "\n"
}

// ── helpers ───────────────────────────────────────────────────────────────────

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func mod(v, n int) int {
	return ((v % n) + n) % n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Run starts the config TUI and returns whether the config was saved.
func Run(cfg *config.Config) (bool, error) {
	m := New(cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	result, err := p.Run()
	if err != nil {
		return false, err
	}
	final, ok := result.(Model)
	if !ok {
		return false, nil
	}
	return final.saved, nil
}
