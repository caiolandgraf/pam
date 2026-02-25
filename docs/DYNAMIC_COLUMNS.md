# Dynamic Column Width Implementation

## Overview

The table viewer now features **dynamic column widths** that automatically adapt to your terminal size and data content. This provides a much better user experience compared to fixed-width columns, as tables now utilize the full available screen space efficiently.

## Features

### ✨ Automatic Sizing
- **Content-aware**: Column widths are calculated based on actual data (sampling up to 100 rows)
- **Header-aware**: Takes into account column names, type icons, primary key indicators, and sort arrows
- **Terminal-aware**: Automatically uses all available horizontal space in your terminal

### 📏 Smart Constraints
- **Minimum width**: 8 characters (ensures readability)
- **Maximum width**: 50 characters (prevents single columns from dominating)
- **Proportional distribution**: Extra space is distributed intelligently among columns

### 🔄 Responsive Behavior
- **Live resizing**: When you resize your terminal, column widths recalculate automatically
- **Adaptive scrolling**: Horizontal scrolling adjusts to show the maximum number of columns that fit
- **Graceful degradation**: If terminal is too narrow, prioritizes showing at least one column

## How It Works

### 1. Width Calculation Algorithm

The `calculateColumnWidths()` function performs the following steps:

```
For each column:
  1. Calculate header width (including icons: type, PK, sort)
  2. Sample up to 100 rows of data
  3. Find the maximum content width
  4. Apply min/max constraints (8-50 chars)

Total needed = sum of all widths + borders between columns
```

### 2. Space Distribution

**If content exceeds available space:**
- Scale down all columns proportionally
- Respect minimum width constraint
- Trim largest columns first if still too wide

**If extra space is available:**
- Distribute evenly among all columns
- Respect maximum width constraint per column
- Add remainder to first columns

### 3. Visible Columns Calculation

The `calculateVisibleColumns()` function:
- Starts from the current horizontal offset (`offsetX`)
- Accumulates column widths until available space is exhausted
- Accounts for border characters between columns
- Ensures at least one column is always visible

## Implementation Details

### New Model Fields

```go
type Model struct {
    // ... existing fields ...
    cellWidth    int   // Fallback width for edge cases
    columnWidths []int // Dynamic width for each column
}
```

### Key Functions

#### `calculateColumnWidths()`
**Purpose**: Computes optimal width for each column based on content and available space

**Constants**:
- `minWidth = 8`: Minimum readable width
- `maxWidth = 50`: Maximum width to prevent dominance
- `borderWidth = 1`: Space for column separators

**Process**:
1. Analyze content (headers + sample rows)
2. Apply constraints
3. Fit to available terminal width
4. Distribute excess space

#### `calculateVisibleColumns()`
**Purpose**: Determines how many columns can fit in the current viewport

**Logic**:
- Iterates from `offsetX` to end of columns
- Accumulates widths until terminal width exceeded
- Returns count of visible columns

### Integration Points

#### Window Resize Handler
```go
func (m Model) handleWindowResize(msg tea.WindowSizeMsg) Model {
    m.width = msg.Width
    m.height = msg.Height
    
    // Trigger dynamic width calculation
    m.calculateColumnWidths()
    
    // ... rest of resize logic ...
}
```

#### View Rendering
The `view.go` functions now check for dynamic widths:

```go
// Use dynamic width if available, otherwise fallback
colWidth := m.cellWidth
if j < len(m.columnWidths) {
    colWidth = m.columnWidths[j]
}

content := formatCell(m.data[rowIndex][j], colWidth)
```

## Performance Considerations

### Sampling Strategy
- Only samples **first 100 rows** for width calculation
- Balances accuracy with performance
- For datasets with millions of rows, avoids scanning entire table

### Recalculation Triggers
Width recalculation only occurs on:
- Initial table load
- Terminal resize events
- **NOT** on navigation or data updates (uses cached widths)

### Memory Efficiency
- `columnWidths` slice is O(n) where n = number of columns
- Typically very small (most tables have < 50 columns)
- No significant memory overhead

## Edge Cases Handled

### Very Narrow Terminals
- Ensures at least 1 column is always visible
- May truncate content with ellipsis (…)
- Horizontal scrolling still works

### Very Wide Terminals
- Distributes extra space among columns
- Respects maximum width constraint
- Creates pleasant, readable layout

### No Data
- Uses header information only
- Falls back to sensible defaults
- Handles empty result sets gracefully

### Mixed Content Types
- Numeric columns tend to be narrower
- Text columns can expand up to max width
- Type icons provide visual hints

## Configuration

### Fallback Width
The `default_column_width` config option is still respected as a fallback:

```yaml
default_column_width: 15
```

This is used when:
- `columnWidths` hasn't been initialized yet
- Edge cases where dynamic calculation fails
- Backward compatibility

### Future Enhancements
Potential improvements for future versions:

- [ ] Per-column width preferences in config
- [ ] Smart detection of UUID/email columns (fixed widths)
- [ ] User-adjustable column widths during runtime
- [ ] Remember column widths per query/table
- [ ] Weighted distribution (give more space to important columns)

## Testing

To test dynamic column behavior:

1. **Resize test**: Run a query and resize your terminal window
2. **Wide data test**: Query a table with very long text fields
3. **Many columns test**: Query a table with 20+ columns
4. **Narrow terminal test**: Shrink terminal to minimum width

Example test query:
```sql
SELECT 
    id,
    name,
    email,
    description,
    created_at,
    updated_at,
    status
FROM users
LIMIT 100;
```

## Migration Notes

### Breaking Changes
None! This is a backward-compatible enhancement.

### Behavioral Changes
- Columns now take up full terminal width (instead of leaving empty space)
- Column widths vary based on content (instead of being uniform)
- Horizontal scrolling may show different number of columns depending on their widths

### Configuration Impact
- `default_column_width` is now a fallback only
- Existing configs continue to work without changes
- No migration required

## Troubleshooting

### Columns seem too narrow
- Check your terminal width (should be at least 80 chars)
- Verify `minWidth` constant in code (currently 8)
- Consider if your data has many columns competing for space

### Columns seem too wide
- Check `maxWidth` constant in code (currently 50)
- This is intentional to prevent single columns from dominating

### Performance issues
- Width calculation uses sampling (100 rows max)
- Should be negligible even for large datasets
- Only recalculates on terminal resize

## Credits

This feature was implemented to provide a modern, responsive TUI experience that adapts to user's terminal environment, maximizing readability and usability of the table viewer.