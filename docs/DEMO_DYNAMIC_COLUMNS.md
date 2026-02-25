# Dynamic Columns Demo Guide

## Quick Demo

This guide demonstrates the new dynamic column width feature in PAM's table viewer.

## What You'll See

### Before: Fixed Width Columns
Previously, all columns had a fixed width (default 15 characters), which meant:
- Small terminal = horizontal scrolling required even for narrow data
- Large terminal = wasted empty space on the right
- Inconsistent readability across different data types

### After: Dynamic Width Columns
Now columns adapt intelligently:
- ✅ Full terminal width is utilized
- ✅ Columns sized based on actual content
- ✅ Automatic resizing when terminal changes
- ✅ Smart min/max constraints for readability

## Live Demo Steps

### Step 1: Basic Query with Dynamic Widths

```bash
# If you have a test database configured
./pam run your_query_name

# Or run a direct query
./pam tables users
```

**What to observe:**
- Columns automatically fill your terminal width
- Wider data (like emails or descriptions) get more space
- Narrow data (like IDs or ages) use less space
- No wasted empty space on the right

### Step 2: Resize Your Terminal

While viewing a table:
1. Make your terminal **narrower** → columns shrink proportionally
2. Make your terminal **wider** → columns expand to use available space
3. The table re-renders automatically!

**What to observe:**
- Instant adaptation to new terminal size
- Horizontal scroll adjusts to show maximum fitting columns
- No manual refresh needed

### Step 3: Navigate Horizontally

Use `h` (left) and `l` (right) to scroll through columns:

**What to observe:**
- Different columns may have different widths
- Visible column count changes based on their individual widths
- Status bar shows current position

### Step 4: Sort Columns

Press `f` on different columns to toggle sort:

**What to observe:**
- Sort indicators (↑↓•) appear in column headers
- Column width includes space for sort indicators
- Table re-queries with ORDER BY clause

## Example Scenarios

### Scenario A: Many Narrow Columns
```sql
SELECT id, age, status, active, score, rank FROM users;
```
**Result**: All columns visible at once, each using minimal space

### Scenario B: Mix of Wide and Narrow
```sql
SELECT id, name, email, description, created_at FROM users;
```
**Result**: 
- `id` → narrow (~8 chars)
- `name` → medium (~15-20 chars)
- `email` → wider (~25-30 chars)
- `description` → max width (50 chars)
- `created_at` → medium (~20 chars)

### Scenario C: Very Wide Data
```sql
SELECT id, very_long_json_column, another_long_text FROM data;
```
**Result**:
- Horizontal scroll enabled
- Each visible column uses its optimal width
- Navigate with `h`/`l` keys

## Visual Comparison

### Fixed Width (Old Behavior)
```
┌──────────────┬──────────────┬──────────────┬──────────────┐
│ id           │ name         │ email        │ status       │
├──────────────┼──────────────┼──────────────┼──────────────┤
│ 1            │ Alice        │ alice@exa... │ active       │
│ 2            │ Bob          │ bob@examp... │ inactive     │
└──────────────┴──────────────┴──────────────┴──────────────┘
                                                              [wasted space →]
```

### Dynamic Width (New Behavior)
```
┌────┬──────────────┬───────────────────────────┬──────────┐
│ id │ name         │ email                     │ status   │
├────┼──────────────┼───────────────────────────┼──────────┤
│ 1  │ Alice        │ alice@example.com         │ active   │
│ 2  │ Bob          │ bob@example.com           │ inactive │
└────┴──────────────┴───────────────────────────┴──────────┘
[uses full terminal width]
```

## Testing Different Terminal Sizes

### Narrow Terminal (80 columns)
- Shows 2-4 columns at once (depending on content)
- Horizontal scroll for additional columns
- Minimum column width: 8 characters

### Medium Terminal (120 columns)
- Shows 4-6 columns at once
- Good balance of width and coverage
- Recommended size for most use cases

### Wide Terminal (200+ columns)
- Shows most/all columns at once
- Columns expand to fill space (up to max 50 chars each)
- Extra space distributed proportionally

## Advanced Features

### Content Sampling
The algorithm samples **up to 100 rows** to determine optimal widths:
- Fast even with millions of rows
- Accurate for typical data distributions
- Avoids performance overhead

### Type-Aware Sizing
Column type indicators affect width calculation:
- `α` (text) → tends to be wider
- `№` (integer) → tends to be narrower
- `◷` (timestamp) → fixed typical width
- `⚿` (primary key) → adds indicator space

### Smart Truncation
When content exceeds column width:
- Text is truncated with `…` ellipsis
- Press `Enter` on cell to view full content
- Detail view shows complete untruncated data

## Tips & Tricks

### Maximize Readability
1. Resize terminal to ~120-150 columns for best experience
2. Use `f` to sort by important columns
3. Press `Enter` to view long content in detail view

### Work with Many Columns
1. Use `h`/`l` to navigate horizontally
2. Press `0` to jump to first column
3. Press `$` to jump to last column
4. Check status bar for position: `[row/col]`

### Optimize for Your Data
- Wide text fields automatically get more space
- Numeric IDs stay compact
- JSON columns capped at max width (use Enter for full view)

## Troubleshooting

### "Columns too narrow to read"
- Increase terminal width
- Reduce number of columns in query
- Check `default_column_width` config

### "Not all columns visible"
- This is expected if total width > terminal width
- Use `h`/`l` to scroll horizontally
- Consider wider terminal or fewer columns

### "Columns not resizing"
- Ensure you're on latest version with dynamic width feature
- Try manually resizing terminal
- Check for any error messages

## Performance Notes

- Width calculation: ~1-2ms for typical tables
- No impact on query execution time
- Scales well with large column counts (100+ columns tested)
- Memory footprint: negligible (O(n) where n = column count)

## Feedback

This feature is designed to provide the best possible viewing experience across different terminal sizes and data types. If you encounter any issues or have suggestions for improvement, please open an issue on GitHub!

## Next Steps

Try these features next:
- Press `v` for visual selection mode
- Press `y` to copy cell content
- Press `u` to update cell (if table has primary key)
- Press `e` to edit and re-run the query
- Press `f` to sort by column

Happy querying! 🎉