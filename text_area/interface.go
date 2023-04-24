package textarea

import bubble_bath "github.com/mieubrisse/bubble-bath"

type Component interface {
	bubble_bath.InteractiveComponent

	/* ---- getters ----- */

	GetValue() string
	GetLength() int
	GetNumRows() int
	GetLineInfo() LineInfo
	GetCursorColumn() int
	GetCursorRow() int

	/* ---- prompt func ----- */

	SetPromptFunc(promptWidth int, fn func(lineIdx int) string)

	/* ---- editing ----- */

	SetValue(s string)
	Reset()
	InsertString(s string)
	InsertRune(r rune)
	DeleteBeforeCursor()
	DeleteAfterCursor()
	DeleteOnCursor() []rune
	InsertLineAbove()
	InsertLineBelow()
	DeleteLine()
	ClearLine()

	/* ---- cursor movement ----- */

	SetCursorColumn(col int)
	MoveCursorDown(bindToLine bool)
	MoveCursorUp(bindToLine bool)
	MoveCursorToLineStart()
	MoveCursorToLineEnd(bindToLine bool)
	MoveCursorRightOneRune(bindToLine bool)
	MoveCursorLeftOneRune()
	MoveCursorByWord(direction CursorMovementDirection, stopPosition WordwiseMovementStopPosition)

	SetCursorRow(targetRow int)
	MoveCursorToFirstRow()
	MoveCursorToLastRow()
}
