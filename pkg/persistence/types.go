package persistence

import "time"

// SaveInfo contains metadata about a save file
type SaveInfo struct {
	Filename  string    // File name without extension
	Version   string    // Save format version
	Timestamp time.Time // When the game was saved
	Size      int64     // File size in bytes
	ModTime   time.Time // Last modification time
}
