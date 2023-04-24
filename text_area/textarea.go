package textarea

import (
	"fmt"
	bubble_bath "github.com/mieubrisse/bubble-bath"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/runeutil"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	rw "github.com/mattn/go-runewidth"
)

const (
	minHeight        = 1
	minWidth         = 2
	defaultHeight    = 6
	defaultWidth     = 40
	defaultCharLimit = 400
	maxHeight        = 99
	maxWidth         = 500

	lineNumberColorHex = "#5d5d5d"
)

// Paste is a tea.Cmd for pasting from the clipboard into the text input.
func Paste() tea.Msg {
	str, err := clipboard.ReadAll()
	if err != nil {
		return pasteErrMsg{err}
	}
	return pasteMsg(str)
}

// Internal messages for clipboard operations.
type pasteMsg string
type pasteErrMsg struct{ error }

// Model is the Bubble Tea model for this text area element.
type Model struct {
	Err error

	// General settings.

	// Prompt is printed at the beginning of each line.
	//
	// When changing the value of Prompt after the model has been
	// initialized, ensure that SetWidth() gets called afterwards.
	//
	// See also SetPromptFunc().
	Prompt string

	// Placeholder is the text displayed when the user
	// hasn't entered anything yet.
	Placeholder string

	// ShowLineNumbers, if enabled, causes line numbers to be printed
	// after the prompt.
	ShowLineNumbers bool

	// EndOfBufferCharacter is displayed at the end of the input.
	EndOfBufferCharacter rune

	// KeyMap encodes the keybindings recognized by the widget.
	KeyMap KeyMap

	// Styling. FocusedStyle and BlurredStyle are used to style the textarea in
	// focused and blurred states.
	FocusedStyle Style
	BlurredStyle Style
	// style is the current styling to use.
	// It is used to abstract the differences in focus state when styling the
	// model, since we can simply assign the set of styles to this variable
	// when switching focus states.
	style *Style

	// Cursor is the text area cursor.
	Cursor cursor.Model

	// CharLimit is the maximum number of characters this input element will
	// accept. If 0 or less, there's no limit.
	CharLimit int

	// If promptFunc is set, it replaces Prompt as a generator for
	// prompt strings at the beginning of each line.
	promptFunc func(line int) string

	// promptWidth is the width of the prompt.
	promptWidth int

	// width is the maximum number of characters that can be displayed at once.
	// If 0 or less this setting is ignored.
	width int

	// height is the maximum number of lines that can be displayed at once. It
	// essentially treats the text field like a vertically scrolling viewport
	// if there are more lines than the permitted height.
	height int

	// Underlying text value.
	value [][]rune

	// focus indicates whether user input focus should be on this input
	// component. When false, ignore keyboard input and hide the cursor.
	focus bool

	// Cursor column in the 'value' rune grid
	col int

	// Cursor row in the 'value' rune grid
	row int

	// Last character offset, used to maintain state when the cursor is moved
	// vertically such that we can maintain the same navigating position.
	lastCharOffset int

	// lineNumberFormat is the format string used to display line numbers.
	lineNumberFormat string

	// viewport is the vertically-scrollable viewport of the multi-line text
	// input.
	viewport *viewport.Model

	// rune sanitizer for input.
	rsan runeutil.Sanitizer
}

// New creates a new model with default settings.
func New() Model {
	vp := viewport.New(0, 0)
	vp.KeyMap = viewport.KeyMap{}
	cur := cursor.New()

	focusedStyle, blurredStyle := DefaultStyles()

	m := Model{
		CharLimit:            defaultCharLimit,
		Prompt:               lipgloss.ThickBorder().Left + " ",
		style:                &blurredStyle,
		FocusedStyle:         focusedStyle,
		BlurredStyle:         blurredStyle,
		EndOfBufferCharacter: '~',
		ShowLineNumbers:      true,
		Cursor:               cur,
		KeyMap:               DefaultKeyMap,

		value:            make([][]rune, minHeight, maxHeight),
		focus:            false,
		col:              0,
		row:              0,
		lineNumberFormat: "%2v ",

		viewport: &vp,
	}

	m.Resize(defaultWidth, defaultHeight)

	return m
}

// GetValue returns the value of the text input.
func (m *Model) GetValue() string {
	if m.value == nil {
		return ""
	}

	var v strings.Builder
	for _, l := range m.value {
		v.WriteString(string(l))
		v.WriteByte('\n')
	}

	return strings.TrimSuffix(v.String(), "\n")
}

// GetLength returns the number of characters currently in the text input.
func (m *Model) GetLength() int {
	var l int
	for _, row := range m.value {
		l += rw.StringWidth(string(row))
	}
	// We add len(m.value) to include the newline characters.
	return l + len(m.value) - 1
}

// GetNumRows returns the number of lines that are currently in the text input.
func (m *Model) GetNumRows() int {
	return len(m.value)
}

// SetPromptFunc supersedes the Prompt field and sets a dynamic prompt
// instead.
// If the function returns a prompt that is shorter than the
// specified promptWidth, it will be padded to the left.
// If it returns a prompt that is longer, display artifacts
// may occur; the caller is responsible for computing an adequate
// promptWidth.
func (m *Model) SetPromptFunc(promptWidth int, fn func(lineIdx int) string) {
	m.promptFunc = fn
	m.promptWidth = promptWidth
}

// GetCursorColumn gets the column within the rune grid where the cursor is currently at
// Note that the cursor can be beyond the right-hand end of the rune grid!
func (m *Model) GetCursorColumn() int {
	return m.col
}

func (m *Model) GetCursorRow() int {
	return m.row
}

// IsFocused returns the focus state on the model.
func (m *Model) IsFocused() bool {
	return m.focus
}

func (m *Model) SetFocus(isFocused bool) tea.Cmd {
	m.focus = isFocused

	var cmd tea.Cmd
	if isFocused {
		m.style = &m.FocusedStyle
		cmd = m.Cursor.Focus()
	} else {
		m.style = &m.BlurredStyle
	}
	return cmd
}

// GetLineInfo returns the number of characters from the start of the
// (soft-wrapped) line and the (soft-wrapped) line width.
func (m *Model) GetLineInfo() LineInfo {
	grid := wrap(m.value[m.row], m.width)

	// Find out which line we are currently on. This can be determined by the
	// m.col and counting the number of runes that we need to skip.
	var counter int
	for i, line := range grid {
		// We've found the line that we are on
		if counter+len(line) == m.col && i+1 < len(grid) {
			// We wrap around to the next line if we are at the end of the
			// previous line so that we can be at the very beginning of the row
			return LineInfo{
				CharOffset:   0,
				ColumnOffset: 0,
				Height:       len(grid),
				RowOffset:    i + 1,
				StartColumn:  m.col,
				Width:        len(grid[i+1]),
				CharWidth:    rw.StringWidth(string(line)),
			}
		}

		if counter+len(line) >= m.col {
			return LineInfo{
				CharOffset:   rw.StringWidth(string(line[:bubble_bath.GetMaxInt(0, m.col-counter)])),
				ColumnOffset: m.col - counter,
				Height:       len(grid),
				RowOffset:    i,
				StartColumn:  counter,
				Width:        len(line),
				CharWidth:    rw.StringWidth(string(line)),
			}
		}

		counter += len(line)
	}
	return LineInfo{}
}

// GetWidth returns the width of the textarea.
func (m *Model) GetWidth() int {
	return m.width
}

// GetHeight returns the current height of the textarea.
func (m *Model) GetHeight() int {
	return m.height
}

// Resize sets the height & width of the textarea to fit exactly within the given dimensions.
// This means that the textarea will account for the width of the prompt and
// whether or not line numbers are being shown.
//
// Resize should be called after setting the Prompt and ShowLineNumbers,
// It is important that the width of the textarea be exactly the given width
// and no more.
func (m *Model) Resize(width int, height int) {
	m.viewport.Width = bubble_bath.Clamp(width, minWidth, maxWidth)

	// Since the width of the textarea input is dependant on the width of the
	// prompt and line numbers, we need to calculate it by subtracting.
	inputWidth := width
	if m.ShowLineNumbers {
		inputWidth -= rw.StringWidth(fmt.Sprintf(m.lineNumberFormat, 0))
	}

	// Account for base style borders and padding.
	inputWidth -= m.style.Base.GetHorizontalFrameSize()

	if m.promptFunc == nil {
		m.promptWidth = rw.StringWidth(m.Prompt)
	}

	inputWidth -= m.promptWidth
	m.width = bubble_bath.Clamp(inputWidth, minWidth, maxWidth)

	m.height = bubble_bath.Clamp(height, minHeight, maxHeight)

	m.viewport.Height = bubble_bath.Clamp(height, minHeight, maxHeight)
}

// Update is the Bubble Tea update loop.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.focus {
		m.Cursor.Blur()
		return nil
	}

	// Used to determine if the cursor should blink.
	oldRow, oldCol := m.cursorLineNumber(), m.col

	var cmds []tea.Cmd

	if m.value[m.row] == nil {
		m.value[m.row] = make([]rune, 0)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.DeleteAfterCursor):
			m.col = bubble_bath.Clamp(m.col, 0, len(m.value[m.row]))
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
				break
			}
			m.DeleteAfterCursor()
		case key.Matches(msg, m.KeyMap.DeleteBeforeCursor):
			m.col = bubble_bath.Clamp(m.col, 0, len(m.value[m.row]))
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				break
			}
			m.DeleteBeforeCursor()
		case key.Matches(msg, m.KeyMap.DeleteCharacterBackward):
			m.col = bubble_bath.Clamp(m.col, 0, len(m.value[m.row]))
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				break
			}
			if len(m.value[m.row]) > 0 {
				m.value[m.row] = append(m.value[m.row][:bubble_bath.GetMaxInt(0, m.col-1)], m.value[m.row][m.col:]...)
				if m.col > 0 {
					m.SetCursorColumn(m.col - 1)
				}
			}
		case key.Matches(msg, m.KeyMap.DeleteCharacterForward):
			if len(m.value[m.row]) > 0 && m.col < len(m.value[m.row]) {
				m.value[m.row] = append(m.value[m.row][:m.col], m.value[m.row][m.col+1:]...)
			}
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
				break
			}
		case key.Matches(msg, m.KeyMap.DeleteWordBackward):
			if m.col <= 0 {
				m.mergeLineAbove(m.row)
				break
			}
			m.deleteWordLeft()
		case key.Matches(msg, m.KeyMap.DeleteWordForward):
			m.col = bubble_bath.Clamp(m.col, 0, len(m.value[m.row]))
			if m.col >= len(m.value[m.row]) {
				m.mergeLineBelow(m.row)
				break
			}
			m.deleteWordRight()
		case key.Matches(msg, m.KeyMap.InsertNewline):
			if len(m.value) >= maxHeight {
				return nil
			}
			m.col = bubble_bath.Clamp(m.col, 0, len(m.value[m.row]))
			m.splitLine(m.row, m.col)
		case key.Matches(msg, m.KeyMap.LineEnd):
			// If the user is going to the end of the line, do allow them to go off the end of the line
			m.MoveCursorToLineEnd(false)
		case key.Matches(msg, m.KeyMap.LineStart):
			m.MoveCursorToLineStart()
		case key.Matches(msg, m.KeyMap.CharacterForward):
			// When we're editing normally, we DO want to allow moving off the end of the line
			m.MoveCursorRightOneRune(false)
		case key.Matches(msg, m.KeyMap.LineNext):
			// If the user is using the arrow keys, we don't want to be binding to the end of the line because they're
			// almost definitely in insert mode
			m.MoveCursorDown(false)
		case key.Matches(msg, m.KeyMap.WordForward):
			m.MoveCursorByWord(Right, Incidence)
		case key.Matches(msg, m.KeyMap.Paste):
			return Paste
		case key.Matches(msg, m.KeyMap.CharacterBackward):
			m.MoveCursorLeftOneRune()
		case key.Matches(msg, m.KeyMap.LinePrevious):
			// If the user is using the arrow keys, we don't want to be binding to the end of the line because they're
			// almost definitely in insert mode
			m.MoveCursorUp(false)
		case key.Matches(msg, m.KeyMap.WordBackward):
			// Note that "End" here is actually the start of the word
			m.MoveCursorByWord(Left, Terminus)
		case key.Matches(msg, m.KeyMap.InputBegin):
			m.MoveCursorToFirstRow()
			m.MoveCursorToLineStart()
		case key.Matches(msg, m.KeyMap.InputEnd):
			m.MoveCursorToLastRow()
			m.MoveCursorToLineEnd(false)

		default:
			m.insertRunesFromUserInput(msg.Runes)
		}

	case pasteMsg:
		m.insertRunesFromUserInput([]rune(msg))

	case pasteErrMsg:
		m.Err = msg
	}

	vp, cmd := m.viewport.Update(msg)
	m.viewport = &vp
	cmds = append(cmds, cmd)

	newRow, newCol := m.cursorLineNumber(), m.col
	m.Cursor, cmd = m.Cursor.Update(msg)
	if newRow != oldRow || newCol != oldCol {
		m.Cursor.Blink = false
		cmd = m.Cursor.BlinkCmd()
	}
	cmds = append(cmds, cmd)

	m.repositionView()

	return tea.Batch(cmds...)
}

// View renders the text area in its current state.
func (m *Model) View() string {
	if m.GetValue() == "" && m.row == 0 && m.col == 0 && m.Placeholder != "" {
		return m.placeholderView()
	}
	m.Cursor.TextStyle = m.style.CursorLine

	var s strings.Builder
	var style lipgloss.Style
	lineInfo := m.GetLineInfo()

	var newLines int

	displayLine := 0
	for l, line := range m.value {
		wrappedLines := wrap(line, m.width)

		if m.row == l {
			style = m.style.CursorLine
		} else {
			style = m.style.Text
		}

		for wl, wrappedLine := range wrappedLines {
			prompt := m.getPromptString(displayLine)
			prompt = m.style.Prompt.Render(prompt)
			s.WriteString(style.Render(prompt))
			displayLine++

			if m.ShowLineNumbers {
				if wl == 0 {
					if m.row == l {
						s.WriteString(style.Render(m.style.CursorLineNumber.Render(fmt.Sprintf(m.lineNumberFormat, l+1))))
					} else {
						s.WriteString(style.Render(m.style.LineNumber.Render(fmt.Sprintf(m.lineNumberFormat, l+1))))
					}
				} else {
					s.WriteString(m.style.LineNumber.Render(style.Render("   ")))
				}
			}

			strwidth := rw.StringWidth(string(wrappedLine))
			padding := m.width - strwidth
			// If the trailing space causes the line to be wider than the
			// width, we should not draw it to the screen since it will result
			// in an extra space at the end of the line which can look off when
			// the cursor line is showing.
			if strwidth > m.width {
				// The character causing the line to be wider than the width is
				// guaranteed to be a space since any other character would
				// have been wrapped.
				wrappedLine = []rune(strings.TrimSuffix(string(wrappedLine), " "))
				padding -= m.width - strwidth
			}
			if m.row == l && lineInfo.RowOffset == wl {
				s.WriteString(style.Render(string(wrappedLine[:lineInfo.ColumnOffset])))
				if m.col >= len(line) && lineInfo.CharOffset >= m.width {
					m.Cursor.SetChar(" ")
					s.WriteString(m.Cursor.View())
				} else {
					m.Cursor.SetChar(string(wrappedLine[lineInfo.ColumnOffset]))
					s.WriteString(style.Render(m.Cursor.View()))
					s.WriteString(style.Render(string(wrappedLine[lineInfo.ColumnOffset+1:])))
				}
			} else {
				s.WriteString(style.Render(string(wrappedLine)))
			}
			s.WriteString(style.Render(strings.Repeat(" ", bubble_bath.GetMaxInt(0, padding))))
			s.WriteRune('\n')
			newLines++
		}
	}

	// Always show at least `m.GetHeight` lines at all times.
	// To do this we can simply pad out a few extra new lines in the view.
	for i := 0; i < m.height; i++ {
		prompt := m.getPromptString(displayLine)
		prompt = m.style.Prompt.Render(prompt)
		s.WriteString(prompt)
		displayLine++

		if m.ShowLineNumbers {
			lineNumber := m.style.EndOfBuffer.Render(fmt.Sprintf(m.lineNumberFormat, string(m.EndOfBufferCharacter)))
			s.WriteString(lineNumber)
		}
		s.WriteRune('\n')
	}

	m.viewport.SetContent(s.String())
	return m.style.Base.Render(m.viewport.View())
}

// ====================================================================================================
//	Private Helper Functions
// ====================================================================================================

// repositionView repositions the view of the viewport based on the defined
// scrolling behavior.
func (m *Model) repositionView() {
	min := m.viewport.YOffset
	max := min + m.viewport.Height - 1

	if row := m.cursorLineNumber(); row < min {
		m.viewport.LineUp(min - row)
	} else if row > max {
		m.viewport.LineDown(row - max)
	}
}

func (m *Model) getPromptString(displayLine int) (prompt string) {
	prompt = m.Prompt
	if m.promptFunc == nil {
		return prompt
	}
	prompt = m.promptFunc(displayLine)
	pl := rw.StringWidth(prompt)
	if pl < m.promptWidth {
		prompt = fmt.Sprintf("%*s%s", m.promptWidth-pl, "", prompt)
	}
	return prompt
}

// placeholderView returns the prompt and placeholder view, if any.
func (m *Model) placeholderView() string {
	var (
		s     strings.Builder
		p     = rw.Truncate(m.Placeholder, m.width, "...")
		style = m.style.Placeholder.Inline(true)
	)

	prompt := m.getPromptString(0)
	prompt = m.style.Prompt.Render(prompt)
	s.WriteString(m.style.CursorLine.Render(prompt))

	if m.ShowLineNumbers {
		s.WriteString(m.style.CursorLine.Render(m.style.CursorLineNumber.Render((fmt.Sprintf(m.lineNumberFormat, 1)))))
	}

	m.Cursor.TextStyle = m.style.Placeholder
	m.Cursor.SetChar(string(p[0]))
	s.WriteString(m.style.CursorLine.Render(m.Cursor.View()))

	// The rest of the placeholder text
	s.WriteString(m.style.CursorLine.Render(style.Render(p[1:] + strings.Repeat(" ", bubble_bath.GetMaxInt(0, m.width-rw.StringWidth(p))))))

	// The rest of the new lines
	for i := 1; i < m.height; i++ {
		s.WriteRune('\n')
		prompt := m.getPromptString(i)
		prompt = m.style.Prompt.Render(prompt)
		s.WriteString(prompt)

		if m.ShowLineNumbers {
			eob := m.style.EndOfBuffer.Render((fmt.Sprintf(m.lineNumberFormat, string(m.EndOfBufferCharacter))))
			s.WriteString(eob)
		}
	}

	m.viewport.SetContent(s.String())
	return m.style.Base.Render(m.viewport.View())
}

// cursorLineNumber returns the line number that the cursor is on.
// This accounts for soft wrapped lines.
func (m *Model) cursorLineNumber() int {
	line := 0
	for i := 0; i < m.row; i++ {
		// Calculate the number of lines that the current line will be split
		// into.
		line += len(wrap(m.value[i], m.width))
	}
	line += m.GetLineInfo().RowOffset
	return line
}

// Wrap a rune string into an array of rune strings
func wrap(runes []rune, width int) [][]rune {
	var (
		lines  = [][]rune{{}}
		word   = []rune{}
		row    int
		spaces int
	)

	// Word wrap the runes
	for _, r := range runes {
		if unicode.IsSpace(r) {
			spaces++
		} else {
			word = append(word, r)
		}

		if spaces > 0 {
			if rw.StringWidth(string(lines[row]))+rw.StringWidth(string(word))+spaces > width {
				row++
				lines = append(lines, []rune{})
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			} else {
				lines[row] = append(lines[row], word...)
				lines[row] = append(lines[row], repeatSpaces(spaces)...)
				spaces = 0
				word = nil
			}
		} else {
			// If the last character is a double-width rune, then we may not be able to add it to this line
			// as it might cause us to go past the width.
			lastCharLen := rw.RuneWidth(word[len(word)-1])
			if rw.StringWidth(string(word))+lastCharLen > width {
				// If the current line has any content, let's move to the next
				// line because the current word fills up the entire line.
				if len(lines[row]) > 0 {
					row++
					lines = append(lines, []rune{})
				}
				lines[row] = append(lines[row], word...)
				word = nil
			}
		}
	}

	if rw.StringWidth(string(lines[row]))+rw.StringWidth(string(word))+spaces >= width {
		lines = append(lines, []rune{})
		lines[row+1] = append(lines[row+1], word...)
		// We add an extra space at the end of the line to account for the
		// trailing space at the end of the previous soft-wrapped lines so that
		// behaviour when navigating is consistent and so that we don't need to
		// continually add edges to handle the last line of the wrapped input.
		spaces++
		lines[row+1] = append(lines[row+1], repeatSpaces(spaces)...)
	} else {
		lines[row] = append(lines[row], word...)
		spaces++
		lines[row] = append(lines[row], repeatSpaces(spaces)...)
	}

	return lines
}

func repeatSpaces(n int) []rune {
	return []rune(strings.Repeat(string(' '), n))
}
