package textarea

// CharacterwiseMovementStopPosition is the position that the character will stop when doing characterwise movement, as
// determined by the direction of travel
type CharacterwiseMovementStopPosition int

const (
	// On stops on the character (corresponds to 'f' in Vim)
	On = 0

	// Before stops one character before the target character (corresponds to 't' in Vim)
	Before = 1
)
