package textarea

import (
	rw "github.com/mattn/go-runewidth"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"unicode"
)

// MoveCursorDown moves the cursor down by one line.
// Returns whether or not the cursor blink should be reset.
func (m *Model) MoveCursorDown(bindToLine bool) {
	li := m.GetLineInfo()
	charOffset := bubble_bath.GetMaxInt(m.lastCharOffset, li.CharOffset)
	m.lastCharOffset = charOffset

	if li.RowOffset+1 >= li.Height && m.row < len(m.value)-1 {
		m.row++
		m.col = 0
	} else {
		// Move the cursor to the start of the next line. So that we can get
		// the line information. We need to add 2 columns to account for the
		// trailing space wrapping.
		m.col = bubble_bath.GetMinInt(li.StartColumn+li.Width+2, len(m.value[m.row])-1)
	}

	nli := m.GetLineInfo()
	m.col = nli.StartColumn

	if nli.Width <= 0 {
		return
	}

	stopThreshold := len(m.value[m.row])
	if bindToLine {
		stopThreshold--
	}

	offset := 0
	for offset < charOffset {
		// TODO TESTING
		if m.col >= stopThreshold || offset >= nli.CharWidth-1 {
			break
		}
		offset += rw.RuneWidth(m.value[m.row][m.col])
		m.col++
	}

	m.repositionView()
}

// MoveCursorUp moves the cursor up by one line.
// If bindToLine is set, the cursor will not move past the last character of the line
func (m *Model) MoveCursorUp(bindToLine bool) {
	li := m.GetLineInfo()
	charOffset := bubble_bath.GetMaxInt(m.lastCharOffset, li.CharOffset)
	m.lastCharOffset = charOffset

	if li.RowOffset <= 0 && m.row > 0 {
		m.row--
		m.col = len(m.value[m.row])
	} else {
		// Move the cursor to the end of the previous line.
		// This can be done by moving the cursor to the start of the line and
		// then subtracting 2 to account for the trailing space we keep on
		// soft-wrapped lines.
		m.col = li.StartColumn - 2
	}

	nli := m.GetLineInfo()
	m.col = nli.StartColumn

	if nli.Width <= 0 {
		return
	}

	stopThreshold := len(m.value[m.row])
	if bindToLine {
		stopThreshold--
	}

	offset := 0
	for offset < charOffset {
		if m.col >= stopThreshold || offset >= nli.CharWidth-1 {
			break
		}
		offset += rw.RuneWidth(m.value[m.row][m.col])
		m.col++
	}

	m.repositionView()
}

// SetCursorColumn moves the cursor to the given position. If the position is
// out of bounds the cursor will be moved to the start or end accordingly.
func (m *Model) SetCursorColumn(col int) {
	m.col = bubble_bath.Clamp(col, 0, len(m.value[m.row]))
	// Any time that we move the cursor horizontally we need to reset the last
	// offset so that the horizontal position when navigating is adjusted.
	m.lastCharOffset = 0
}

// MoveCursorToLineStart moves the cursor to the start of the input field.
func (m *Model) MoveCursorToLineStart() {
	m.SetCursorColumn(0)
}

// MoveCursorToLineEnd moves the cursor to the end of the input field.
// If bindToLine is set, only allow going to the last char of the line
func (m *Model) MoveCursorToLineEnd(bindToLine bool) {
	newPosition := len(m.value[m.row])
	if bindToLine {
		newPosition--
	}
	m.SetCursorColumn(newPosition)
}

func (m *Model) SetCursorRow(targetRow int) {
	targetRow = bubble_bath.Clamp(targetRow, 0, len(m.value)-1)

	adjustmentNeeded := targetRow - m.row
	if adjustmentNeeded > 0 {
		for i := 0; i < adjustmentNeeded; i++ {
			m.MoveCursorDown(true)
		}
	} else {
		for i := 0; i > adjustmentNeeded; i-- {
			m.MoveCursorUp(true)
		}
	}

	m.repositionView()
}

func (m *Model) MoveCursorToFirstRow() {
	m.SetCursorRow(0)
}

func (m *Model) MoveCursorToLastRow() {
	m.SetCursorRow(len(m.value) - 1)
}

// MoveCursorRightOneRune moves the cursor one character to the right.
// If bindToLine is set, the cursor will not move psat the last character of the line
func (m *Model) MoveCursorRightOneRune(bindToLine bool) {
	moveRightLimit := len(m.value[m.row])
	if bindToLine {
		moveRightLimit--
	}
	if m.col < moveRightLimit {
		m.SetCursorColumn(m.col + 1)
	}
}

// MoveCursorLeftOneRune moves the cursor one character to the left.
// If insideLine is set, the cursor is moved to the last
// character in the previous line, instead of one past that.
func (m *Model) MoveCursorLeftOneRune() {
	if m.col > 0 {
		m.SetCursorColumn(m.col - 1)
	}
}

func (m *Model) MoveCursorByWord(direction CursorMovementDirection, stopPosition WordwiseMovementStopPosition) {
	m.doWordwiseMovement(direction, stopPosition)
}

// ====================================================================================================
//
//	Private Helper Functions
//
// ====================================================================================================
func (m *Model) doWordwiseMovement(direction CursorMovementDirection, stopPosition WordwiseMovementStopPosition) {
	// This function utilizes the insight that the textarea string can be thought of as a "tape" of words, joined by whitespace
	// With this insight, we can handle both (left,right) and (word_start,word_end) by simply sliding along the tape in
	// the appropriate direction looking for the sequence we want

	// If no lines, abort immediately
	if len(m.value) == 0 {
		return
	}

	// Factor applied to index calculations to account for the desired direction of cursor travel
	directionMultiplier := int(direction)

	// Factor applied to index calculations to account for stopPosition, where "incidence" means "look at the
	// character BEHIND the cursor to figure out if you should stop" and "terminus" means "look at the character
	// AHEAD of the cursor to figure out if you should stop" (where "ahead" and "behind" are measured by direction of
	// cursor travel)
	stopPositionMultiplier := int(stopPosition)

	// Figure out what the row index of each end of the tape is
	limitRowIndex := len(m.value) - 1
	if direction == Left {
		limitRowIndex = 0
	}

	// Our column might be off the right edge of the line; ensure we account for that
	sanitizedColIdx := bubble_bath.GetMinInt(m.col, len(m.value[m.row])-1)

	// At some point the proposed new column might be off either end of the line
	// Therefore, let's first calculate the boundary beyond which we know
	// that the proposed column is off the edge of the line
	limitColIndex := len(m.value[m.row]) - 1
	if direction == Left {
		limitColIndex = 0
	}

	// To start the algo off, shove the cursor column one character in the direction of travel
	// This prevents this being a noop if you're already at a word start/end
	nextColIdx := sanitizedColIdx + directionMultiplier
	for {
		// Now, check if the proposed column would indeed put us off the line
		remainingColsBeforeLimit := directionMultiplier * (limitColIndex - nextColIdx)
		if remainingColsBeforeLimit < 0 {
			// We're off the end; can we go to another line to keep going?
			nextRowIdx := m.row + directionMultiplier
			remainingRowsBeforeLimit := directionMultiplier * (limitRowIndex - nextRowIdx)
			if remainingRowsBeforeLimit < 0 {
				// We're at the end of the "tape"; nothing to do
				return
			}

			// We still have at least one more line, so let's use it (which means we're crossing a newline char which
			// means we're crossing a word boundary)
			m.row += directionMultiplier
			if direction == Right {
				nextColIdx = 0
			} else {
				nextColIdx = len(m.value[m.row]) - 1
			}
		}

		// Now that we know we're moving the column index to a valid line, move it
		m.SetCursorColumn(nextColIdx)
		nextColIdx = m.col + directionMultiplier // Prep for next iteration

		// We still might have moved the cursor to an empty line, making the cursor location invalid!
		// Vim will stop on these empty lines, so we try to as well
		if len(m.value[m.row]) == 0 {
			return
		}

		cursorChar := m.value[m.row][m.col]

		// Grab a comparison column which must be whitespace to stop the algorithm
		// The stopPosition multiplier means that "incident" will require whitespace to be *behind* the cursor in the direction
		// of algorithm travel, whereas "terminus" will require whitespace to be *ahead* of the cursor
		// in the direction of algorithm travel
		cursorAdjacentColIdx := m.col + stopPositionMultiplier*directionMultiplier

		// Depending on stopPosition, the adjacent col might actually be *behind* the cursor, meaning we need to use the opposite boundary
		// as the direction of algo travel along the "tape"
		// This happens if we're moving right but we're stopping at word incidence, OR if we're moving left and stopping at word terminus
		// NOTE: This is actually an XOR
		var adjacentColLimitIndex int
		if ((direction == Right) && (stopPosition == Terminus)) || ((direction == Left) && (stopPosition == Incidence)) {
			// The adjacent col limit index will be the list's right limit
			adjacentColLimitIndex = len(m.value[m.row]) - 1
		} else {
			adjacentColLimitIndex = 0
		}

		// Now grab a character that may or may not be whitespace
		// If the col we're grabbing is off the end of the line, it's automatically a newline
		remainingColsBeforeCursorAdjacentColIsOff := stopPositionMultiplier * directionMultiplier * (adjacentColLimitIndex - cursorAdjacentColIdx)
		var candidateWhitespaceChar rune
		if remainingColsBeforeCursorAdjacentColIsOff < 0 {
			candidateWhitespaceChar = '\n'
		} else {
			candidateWhitespaceChar = m.value[m.row][cursorAdjacentColIdx]
		}

		// Evaluate if we reached our target
		if !unicode.IsSpace(cursorChar) && unicode.IsSpace(candidateWhitespaceChar) {
			return
		}

		// We're not done, so prep for next iteration
	}
}

func (m *Model) doCharacterwiseMovement(targetChar rune, direction CursorMovementDirection, stopPosition CharacterwiseMovementStopPosition) {
	directionMultiplier := int(direction)

	stopPositionOffset := int(stopPosition)

	// To start the algo off, shove the proposed column one character in the direction of travel
	// This prevents this being a noop if you're already on the character you're looking for
	newColIdx := m.col + directionMultiplier

	// Based on the stop position, we need to examine different cells (either on the cursor, or ahead of it)
	examinationColIdx := newColIdx + stopPositionOffset

	for {
		// If our examintion column is out-of-bounds, abort; we haven't found anything
		if examinationColIdx < 0 || examinationColIdx >= len(m.value[m.row]) {
			return
		}

		if m.value[m.row][examinationColIdx] == targetChar {
			m.SetCursorColumn(newColIdx)
			return
		}

		newColIdx += directionMultiplier
		examinationColIdx += directionMultiplier
	}
}
