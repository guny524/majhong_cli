package persistence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/guny524/majhong_cli/pkg/mahjong/engine"
)

const (
	// DefaultFileExtension is the default file extension for save files
	DefaultFileExtension = ".mjson"
)

// Storage interface defines methods for saving and loading game states
type Storage interface {
	Save(state *engine.GameState, filename string) error
	Load(filename string) (*engine.GameState, error)
}

// FileStorage implements file-based storage
type FileStorage struct {
	baseDir string
}

// NewFileStorage creates a new file storage with the given base directory
func NewFileStorage(baseDir string) *FileStorage {
	return &FileStorage{
		baseDir: baseDir,
	}
}

// Save saves a game state to a file
func (fs *FileStorage) Save(state *engine.GameState, filename string) error {
	if state == nil {
		return fmt.Errorf("cannot save nil game state")
	}

	// Serialize the game state
	data, err := Serialize(state)
	if err != nil {
		return fmt.Errorf("failed to serialize game state: %w", err)
	}

	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != DefaultFileExtension {
		filename = filename + DefaultFileExtension
	}

	// Create full path
	fullPath := filepath.Join(fs.baseDir, filename)

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Load loads a game state from a file
func (fs *FileStorage) Load(filename string) (*engine.GameState, error) {
	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != DefaultFileExtension {
		filename = filename + DefaultFileExtension
	}

	// Create full path
	fullPath := filepath.Join(fs.baseDir, filename)

	// Read file
	data, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save file not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Deserialize the game state
	state, err := Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize game state: %w", err)
	}

	// Validate the loaded state
	// First convert to serializable format for validation
	saveFormat := convertGameState(state)
	if err := Validate(&saveFormat); err != nil {
		return nil, fmt.Errorf("loaded game state is invalid: %w", err)
	}

	return state, nil
}

// List returns all save files in the base directory
func (fs *FileStorage) List() ([]string, error) {
	files, err := os.ReadDir(fs.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var saveFiles []string
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == DefaultFileExtension {
			// Remove extension for display
			name := file.Name()[:len(file.Name())-len(DefaultFileExtension)]
			saveFiles = append(saveFiles, name)
		}
	}

	return saveFiles, nil
}

// Delete removes a save file
func (fs *FileStorage) Delete(filename string) error {
	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != DefaultFileExtension {
		filename = filename + DefaultFileExtension
	}

	// Create full path
	fullPath := filepath.Join(fs.baseDir, filename)

	// Delete file
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("save file not found: %s", filename)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a save file exists
func (fs *FileStorage) Exists(filename string) bool {
	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != DefaultFileExtension {
		filename = filename + DefaultFileExtension
	}

	// Create full path
	fullPath := filepath.Join(fs.baseDir, filename)

	_, err := os.Stat(fullPath)
	return err == nil
}

// GetInfo returns information about a save file
func (fs *FileStorage) GetInfo(filename string) (*SaveInfo, error) {
	// Ensure the filename has the correct extension
	if filepath.Ext(filename) != DefaultFileExtension {
		filename = filename + DefaultFileExtension
	}

	// Create full path
	fullPath := filepath.Join(fs.baseDir, filename)

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save file not found: %s", filename)
		}
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Read file to get metadata
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON to get version and timestamp
	var saveFormat SaveFormat
	if err := json.Unmarshal(data, &saveFormat); err != nil {
		return nil, fmt.Errorf("failed to parse save file: %w", err)
	}

	return &SaveInfo{
		Filename:    filename[:len(filename)-len(DefaultFileExtension)],
		Version:     saveFormat.Version,
		Timestamp:   saveFormat.Timestamp,
		Size:        info.Size(),
		ModTime:     info.ModTime(),
	}, nil
}
