package engine

// GameError represents a game rule violation or error
type GameError struct {
	Code    string
	Message string
}

func (e *GameError) Error() string {
	return e.Message
}
