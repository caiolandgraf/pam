package table

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) moveUp() Model {
	if m.selectedRow > 0 {
		m.selectedRow--
		if m.selectedRow < m.offsetY {
			m.offsetY = m.selectedRow
		}
	}
	return m
}

func (m Model) moveDown() Model {
	if m.selectedRow < m.numRows()-1 {
		m.selectedRow++
		if m.selectedRow >= m.offsetY+m.visibleRows {
			m.offsetY = m.selectedRow - m.visibleRows + 1
		}
	}
	return m
}

func (m Model) moveLeft() Model {
	if m.selectedCol > 0 {
		m.selectedCol--
		if m.selectedCol < m.offsetX {
			m.offsetX = m.selectedCol
		}
	}
	return m
}

func (m Model) moveRight() Model {
	if m.selectedCol < m.numCols()-1 {
		m.selectedCol++
		if m.selectedCol >= m.offsetX+m.visibleCols {
			m.offsetX = m.selectedCol - m.visibleCols + 1
		}
	}
	return m
}

func (m Model) jumpToFirstCol() Model {
	m.selectedCol = 0
	m.offsetX = 0
	return m
}

func (m Model) jumpToLastCol() Model {
	m.selectedCol = m.numCols() - 1
	if m.visibleCols < m.numCols() {
		m.offsetX = m.numCols() - m.visibleCols
	}
	return m
}

func (m Model) jumpToFirstRow() Model {
	m.selectedRow = 0
	m.offsetY = 0
	return m
}

func (m Model) jumpToLastRow() Model {
	m.selectedRow = m.numRows() - 1
	m.offsetY = m.numRows() - m.visibleRows
	return m
}

func (m Model) pageUp() Model {
	m.selectedRow -= m.visibleRows
	if m.selectedRow < 0 {
		m.selectedRow = 0
	}
	m.offsetY = m.selectedRow
	return m
}

func (m Model) pageDown() Model {
	m.selectedRow += m.visibleRows
	if m.selectedRow >= m.numRows() {
		m.selectedRow = m.numRows() - 1
	}
	if m.selectedRow >= m.offsetY+m.visibleRows {
		m.offsetY = m.selectedRow - m.visibleRows + 1
	}
	return m
}

func (m Model) toggleVisualMode() (Model, tea.Cmd) {
	m.visualMode = !m.visualMode
	
	if m.visualMode {
		m.visualStartRow = m.selectedRow
		m.visualStartCol = m.selectedCol
	}
	
	return m, nil
}

func (m Model) getSelectionBounds() (minRow, maxRow, minCol, maxCol int) {
	if !m.visualMode {
		return m.selectedRow, m.selectedRow, m.selectedCol, m.selectedCol
	}
	
	// Multi-cell selection
	minRow = min(m.visualStartRow, m.selectedRow)
	maxRow = max(m.visualStartRow, m.selectedRow)
	minCol = min(m.visualStartCol, m.selectedCol)
	maxCol = max(m.visualStartCol, m.selectedCol)
	
	return
}

func (m Model) isCellInSelection(row, col int) bool {
	minRow, maxRow, minCol, maxCol := m.getSelectionBounds()
	return row >= minRow && row <= maxRow && col >= minCol && col <= maxCol
}

func (m Model) copySelection() (Model, tea.Cmd) {
	minRow, maxRow, minCol, maxCol := m.getSelectionBounds()
	
	var allRows [][]string
	
	if m.visualMode {
		headerRow := make([]string, 0)
		for col := minCol; col <= maxCol; col++ {
			headerRow = append(headerRow, m.columns[col])
		}
		allRows = append(allRows, headerRow)
	}
	
	for row := minRow; row <= maxRow; row++ {
		dataRow := make([]string, 0)
		for col := minCol; col <= maxCol; col++ {
			dataRow = append(dataRow, m.data[row][col])
		}
		allRows = append(allRows, dataRow)
	}
	
	numCols := maxCol - minCol + 1
	colWidths := make([]int, numCols)
	
	for _, row := range allRows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}
	
	var result strings.Builder
	
	for rowIdx, row := range allRows {
		for colIdx, cell := range row {
			paddedCell := fmt.Sprintf("%-*s", colWidths[colIdx], cell)
			result.WriteString(paddedCell)
			
			if colIdx < len(row)-1 {
				result.WriteString("  ")
			}
		}
		
		if rowIdx < len(allRows)-1 {
			result.WriteString("\n")
		}
	}
	
	content := result.String()
	clipboard.WriteAll(content)
	
	m.visualMode = false
	m.blinkCopiedCell = true
	
	return m, func() tea.Msg {
		time.Sleep(200 * time.Millisecond)
		return blinkMsg{}
	}
}
