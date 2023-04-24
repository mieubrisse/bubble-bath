package textarea

// WordwiseMovementStopPosition is the position in the word that the cursor will stop at *IN THE DIRECTION OF CURSOR TRAVEL*
type WordwiseMovementStopPosition int

const (
	// These integers actually represent the index offset (relative to the direction of cursor travel)
	// where we'll look for a word boundary
	// If we're going right, "incident" = beginning of word and "terminus" = end of word
	// If we're going left, "incident" = *end* of word, and "terminus" = beginning of word
	// We used "incident" and "terminus" to avoid confusion with "start" and "end" when going left

	// Incidence tells the cursor to stop as soon as hit a word (in direction of cursor travel)
	Incidence WordwiseMovementStopPosition = -1

	// Terminus tells the cursor to stop as soon as it would leave a word (in direction of cursor travel)
	Terminus WordwiseMovementStopPosition = 1
)
