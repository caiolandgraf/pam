package table

import "sort"

func (m *Model) toggleMarkRow(row int) {
	if row < 0 || row >= m.numRows() {
		return
	}
	if m.markedRows == nil {
		m.markedRows = make(map[int]bool)
	}
	if m.markedRows[row] {
		delete(m.markedRows, row)
	} else {
		m.markedRows[row] = true
	}
}

func (m *Model) clearMarkedRows() {
	if m.markedRows == nil {
		return
	}
	for k := range m.markedRows {
		delete(m.markedRows, k)
	}
}

func (m Model) hasMarkedRows() bool {
	return len(m.markedRows) > 0
}

func (m Model) markedCount() int {
	return len(m.markedRows)
}

func (m Model) isRowMarked(row int) bool {
	if m.markedRows == nil {
		return false
	}
	return m.markedRows[row]
}

func (m Model) getMarkedRows() []int {
	if m.markedRows == nil || len(m.markedRows) == 0 {
		return nil
	}
	rows := make([]int, 0, len(m.markedRows))
	for r := range m.markedRows {
		rows = append(rows, r)
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i] < rows[j] })
	return rows
}
