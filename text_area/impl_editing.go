package textarea

import (
	"github.com/charmbracelet/bubbles/runeutil"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"unicode"
)

// SetValue sets the value of the text input.
func (m *implementation) SetValue(s string) {
	m.Reset()
	m.InsertString(s)
	m.MoveCursorLeftOneRune()
}

// InsertString inserts a string at the cursor position.
func (m *implementation) InsertString(s string) {
	m.insertRunesFromUserInput([]rune(s))
}

// InsertRune inserts a rune at the cursor position.
func (m *implementation) InsertRune(r rune) {
	m.insertRunesFromUserInput([]rune{r})
}

// DeleteBeforeCursor deletes all text before the cursor. Returns whether or
// not the cursor blink should be reset.
func (m *implementation) DeleteBeforeCursor() {
	m.value[m.row] = m.value[m.row][m.col:]
	m.SetCursorColumn(0)
}

// DeleteAfterCursor deletes all text after the cursor. Returns whether or not
// the cursor blink should be reset. If input is masked delete everything after
// the cursor so as not to reveal word breaks in the masked input.
func (m *implementation) DeleteAfterCursor() {
	m.value[m.row] = m.value[m.row][:m.col]
	m.SetCursorColumn(len(m.value[m.row]) - 1)
}

// DeleteOnCursor deletes the single character on the cursor
// Returns the rune that was deleted (if any)
func (m *implementation) DeleteOnCursor() []rune {
	currentRow := m.value[m.row]
	if len(currentRow) == 0 {
		return make([]rune, 0)
	}

	deletedChar := []rune{currentRow[m.col]}
	newRow := currentRow[:m.col]
	if m.col < len(currentRow)-1 {
		// i.e., there are more characters after the cursor
		newRow = append(newRow, currentRow[m.col+1:]...)
	}
	m.value[m.row] = newRow

	newCol := bubble_bath.GetMinInt(m.col, len(newRow)-1)
	m.SetCursorColumn(newCol)

	return deletedChar
}

func (m *implementation) InsertLineAbove() {
	newValue := make([][]rune, 0, maxHeight)

	preCursorLines := m.value[0:m.row]
	newValue = append(newValue, preCursorLines...)

	newValue = append(newValue, make([]rune, 0))

	cursorLineAndAfter := m.value[m.row:]
	newValue = append(newValue, cursorLineAndAfter...)

	m.row++
	m.value = newValue
}

func (m *implementation) InsertLineBelow() {
	newValue := make([][]rune, 0, maxHeight)

	cursorLineAndPrevious := m.value[0 : m.row+1]
	newValue = append(newValue, cursorLineAndPrevious...)

	newValue = append(newValue, make([]rune, 0))

	postCursorLines := m.value[m.row+1:]
	newValue = append(newValue, postCursorLines...)

	m.value = newValue
}

func (m *implementation) DeleteLine() {
	if len(m.value) <= 1 {
		m.value = make([][]rune, minHeight, maxHeight)
		m.SetCursorColumn(0)
		return
	}

	preCursorLines := m.value[:m.row]
	postCursorLines := make([][]rune, 0)
	if m.row < len(m.value)-1 {
		postCursorLines = m.value[m.row+1:]
	}

	newValue := make([][]rune, 0, maxHeight)
	newValue = append(newValue, preCursorLines...)
	newValue = append(newValue, postCursorLines...)

	m.value = newValue

	m.row = bubble_bath.Clamp(m.row, 0, len(m.value)-1)
}

func (m *implementation) ClearLine() {
	m.value[m.row] = make([]rune, 0, maxWidth)
	m.SetCursorColumn(0)
}

// Reset sets the input to its default state with no input.
func (m *implementation) Reset() {
	m.value = make([][]rune, minHeight, maxHeight)
	m.col = 0
	m.row = 0
	m.viewport.GotoTop()
	m.SetCursorColumn(0)
}

// ====================================================================================================
//	Private Helper Functions
// ====================================================================================================

// rsan initializes or retrieves the rune sanitizer.
func (m *implementation) san() runeutil.Sanitizer {
	if m.rsan == nil {
		// Textinput has all its input on a single line so collapse
		// newlines/tabs to single spaces.
		m.rsan = runeutil.NewSanitizer()
	}
	return m.rsan
}

// insertRunesFromUserInput inserts runes at the current cursor position.
func (m *implementation) insertRunesFromUserInput(runes []rune) {
	// Clean up any special characters in the input provided by the
	// clipboard. This avoids bugs due to e.g. tab characters and
	// whatnot.
	runes = m.san().Sanitize(runes)

	var availSpace int
	if m.CharLimit > 0 {
		availSpace = m.CharLimit - m.GetLength()
		// If the char limit's been reached, cancel.
		if availSpace <= 0 {
			return
		}
		// If there's not enough space to paste the whole thing cut the pasted
		// runes down so they'll fit.
		if availSpace < len(runes) {
			runes = runes[:len(runes)-availSpace]
		}
	}

	// Split the input into lines.
	var lines [][]rune
	lstart := 0
	for i := 0; i < len(runes); i++ {
		if runes[i] == '\n' {
			lines = append(lines, runes[lstart:i])
			lstart = i + 1
		}
	}
	if lstart < len(runes) {
		// The last line did not end with a newline character.
		// Take it now.
		lines = append(lines, runes[lstart:])
	}

	// Obey the maximum height limit.
	if len(m.value)+len(lines)-1 > maxHeight {
		allowedHeight := bubble_bath.GetMaxInt(0, maxHeight-len(m.value)+1)
		lines = lines[:allowedHeight]
	}

	if len(lines) == 0 {
		// Nothing left to insert.
		return
	}

	// Save the reminder of the original line at the current
	// cursor position.
	tail := make([]rune, len(m.value[m.row][m.col:]))
	copy(tail, m.value[m.row][m.col:])

	// Paste the first line at the current cursor position.
	m.value[m.row] = append(m.value[m.row][:m.col], lines[0]...)
	m.col += len(lines[0])

	if numExtraLines := len(lines) - 1; numExtraLines > 0 {
		// Add the new lines.
		// We try to reuse the slice if there's already space.
		var newGrid [][]rune
		if cap(m.value) >= len(m.value)+numExtraLines {
			// Can reuse the extra space.
			newGrid = m.value[:len(m.value)+numExtraLines]
		} else {
			// No space left; need a new slice.
			newGrid = make([][]rune, len(m.value)+numExtraLines)
			copy(newGrid, m.value[:m.row+1])
		}
		// Add all the rows that were after the cursor in the original
		// grid at the end of the new grid.
		copy(newGrid[m.row+1+numExtraLines:], m.value[m.row+1:])
		m.value = newGrid
		// Insert all the new lines in the middle.
		for _, l := range lines[1:] {
			m.row++
			m.value[m.row] = l
			m.col = len(l)
		}
	}

	// Finally add the tail at the end of the last line inserted.
	m.value[m.row] = append(m.value[m.row], tail...)

	m.SetCursorColumn(m.col)
}

// deleteWordLeft deletes the word left to the cursor. Returns whether or not
// the cursor blink should be reset.
func (m *implementation) deleteWordLeft() {
	if m.col == 0 || len(m.value[m.row]) == 0 {
		return
	}

	// Linter note: it's critical that we acquire the initial cursor position
	// here prior to altering it via SetCursorColumn() below. As such, moving this
	// call into the corresponding if clause does not apply here.
	oldCol := m.col //nolint:ifshort

	m.SetCursorColumn(m.col - 1)
	for unicode.IsSpace(m.value[m.row][m.col]) {
		if m.col <= 0 {
			break
		}
		// ignore series of whitespace before cursor
		m.SetCursorColumn(m.col - 1)
	}

	for m.col > 0 {
		if !unicode.IsSpace(m.value[m.row][m.col]) {
			m.SetCursorColumn(m.col - 1)
		} else {
			if m.col > 0 {
				// keep the previous space
				m.SetCursorColumn(m.col + 1)
			}
			break
		}
	}

	if oldCol > len(m.value[m.row]) {
		m.value[m.row] = m.value[m.row][:m.col]
	} else {
		m.value[m.row] = append(m.value[m.row][:m.col], m.value[m.row][oldCol:]...)
	}
}

// deleteWordRight deletes the word right to the cursor.
func (m *implementation) deleteWordRight() {
	if m.col >= len(m.value[m.row]) || len(m.value[m.row]) == 0 {
		return
	}

	oldCol := m.col
	m.SetCursorColumn(m.col + 1)
	for unicode.IsSpace(m.value[m.row][m.col]) {
		// ignore series of whitespace after cursor
		m.SetCursorColumn(m.col + 1)

		if m.col >= len(m.value[m.row]) {
			break
		}
	}

	for m.col < len(m.value[m.row]) {
		if !unicode.IsSpace(m.value[m.row][m.col]) {
			m.SetCursorColumn(m.col + 1)
		} else {
			break
		}
	}

	if m.col > len(m.value[m.row]) {
		m.value[m.row] = m.value[m.row][:oldCol]
	} else {
		m.value[m.row] = append(m.value[m.row][:oldCol], m.value[m.row][m.col:]...)
	}

	m.SetCursorColumn(oldCol)
}

func (m *implementation) doWordRight(fn func(charIdx int, pos int)) {
	haveEncounteredWhitespace := false
	for {
		// If we're at (or beyond) the last char of the line (which may be empty)...
		if m.col >= len(m.value[m.row])-1 {
			// ...and there are no more lines, we're done
			if m.row == len(m.value)-1 {
				return
			}

			// ...and the next line is empty, then move to the next line and we're done
			// This is a bit odd, but it's Vim behaviour
			if len(m.value[m.row+1]) == 0 {
				// Copied from MoveCursorRightOneRune
				m.MoveCursorToLineStart()
				m.row++
				return
			}

			// ...and the next line is not empty, prep to stop on the next word
			m.row++
			m.MoveCursorToLineStart()
			haveEncounteredWhitespace = true
			continue
		}

		charUnderCursor := m.value[m.row][m.col]

		// We've already left a word and found another; we're done
		if haveEncounteredWhitespace && !unicode.IsSpace(charUnderCursor) {
			return
		}

		// We've seen a space, so we're ready to stop
		if unicode.IsSpace(charUnderCursor) {
			haveEncounteredWhitespace = true
		}

		// We don't want the cursor to move off the end of everything
		m.MoveCursorRightOneRune(true)
	}

	/*
		charIdx := 0
		for m.col < len(m.value[m.row]) {
			if unicode.IsSpace(m.value[m.row][m.col]) {
				break
			}
			fn(charIdx, m.col)
			m.SetCursorColumn(m.col + 1)
			charIdx++
		}

	*/
}

// mergeLineBelow merges the current line with the line below.
func (m *implementation) mergeLineBelow(row int) {
	if row >= len(m.value)-1 {
		return
	}

	// To perform a merge, we will need to combine the two lines and then
	m.value[row] = append(m.value[row], m.value[row+1]...)

	// Shift all lines up by one
	for i := row + 1; i < len(m.value)-1; i++ {
		m.value[i] = m.value[i+1]
	}

	// And, remove the last line
	if len(m.value) > 0 {
		m.value = m.value[:len(m.value)-1]
	}
}

// mergeLineAbove merges the current line the cursor is on with the line above.
func (m *implementation) mergeLineAbove(row int) {
	if row <= 0 {
		return
	}

	m.col = len(m.value[row-1])
	m.row = m.row - 1

	// To perform a merge, we will need to combine the two lines and then
	m.value[row-1] = append(m.value[row-1], m.value[row]...)

	// Shift all lines up by one
	for i := row; i < len(m.value)-1; i++ {
		m.value[i] = m.value[i+1]
	}

	// And, remove the last line
	if len(m.value) > 0 {
		m.value = m.value[:len(m.value)-1]
	}
}

func (m *implementation) splitLine(row, col int) {
	// To perform a split, take the current line and keep the content before
	// the cursor, take the content after the cursor and make it the content of
	// the line underneath, and shift the remaining lines down by one
	head, tailSrc := m.value[row][:col], m.value[row][col:]
	tail := make([]rune, len(tailSrc))
	copy(tail, tailSrc)

	m.value = append(m.value[:row+1], m.value[row:]...)

	m.value[row] = head
	m.value[row+1] = tail

	m.col = 0
	m.row++
}
