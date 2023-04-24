package text_block

import (
	"github.com/mieubrisse/bubble-bath/bubble_bath"
)

type Component interface {
	bubble_bath.Component

	// GetContents gets the raw contents of the text block, without truncation
	GetContents() string
}
