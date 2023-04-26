package resizable_text_block

import (
	"github.com/mieubrisse/bubble-bath"
)

type Component interface {
	bubble_bath.Component

	// GetContents gets the raw contents of the text block, without truncation
	GetContents() string

	// The minimum width the content can compress down to, doing as much word-wrapping as possible
	// This will be as small as the longest word
	GetMinimumIntrinsicWidth() int

	// The maximum width of the text, if no wrapping is done whatsoever
	GetMaximumIntrinsicWidth() int
}
