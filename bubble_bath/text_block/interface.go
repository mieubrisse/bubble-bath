package text_block

import (
	"github.com/mieubrisse/bubble-bath/components"
)

type Component interface {
	components.Component

	// TODO remove this??
	GetContents() string
}
