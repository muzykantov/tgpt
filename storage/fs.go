package storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/muzykantov/tgpt/chat"
)

// FS represents a file-based storage system that provides methods to persist and retrieve
// chat-related data structures like History and Statistics to and from the file system.
type FS struct {
	BaseDir string // BaseDir is the base directory for storing and retrieving data files.
}

// SaveHistory persists a given History object to the file system.
// It creates a uniquely named JSON file based on the History ID within the BaseDir.
// If a file with the same name exists, it will be overwritten.
//
// history: The History object to be saved.
//
// Returns:
// error: An error if encountered during file operations or serialization.
func (fs *FS) SaveHistory(_ context.Context, history *chat.History) error {
	// Generate the path to save the history using the ID.
	filename := fmt.Sprintf(
		"history-%d-%d-%s.json",
		history.ID.User,
		history.ID.Chat,
		history.ID.Model,
	)
	path := filepath.Join(fs.BaseDir, filename)

	// Open or create the file.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open or create the file: %w", err)
	}
	defer file.Close()

	// Write the history to the file in JSON format.
	err = history.Write(file)
	if err != nil {
		return fmt.Errorf("error writing the history to the file: %w", err)
	}

	return nil
}

// LoadHistory retrieves a chat history from the file system using the provided ID.
// If the file does not exist, a new History instance is returned.
// It constructs a file name from the ID and tries to open the corresponding JSON file
// within the BaseDir, then deserializes the file's content into a History object.
//
// id: The ID of the chat history to be loaded.
//
// Returns:
// *History: A pointer to the retrieved or newly created History object.
// error: An error if encountered during file operations or deserialization, except for file not found error.
func (fs *FS) LoadHistory(_ context.Context, id chat.ID) (*chat.History, error) {
	// Generate the path to load the history using the ID.
	filename := fmt.Sprintf("history-%d-%d-%s.json", id.User, id.Chat, id.Model)
	path := filepath.Join(fs.BaseDir, filename)

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, return a new instance of chat.History.
			return &chat.History{
				ID:     id,
				Prompt: "",
				Log:    []chat.Message{},
			}, nil
		}
		// For other errors, return an error.
		return nil, fmt.Errorf("could not open the file: %w", err)
	}
	defer file.Close()

	// Decode the history from the file.
	history := new(chat.History)
	err = history.Read(file)
	if err != nil {
		return nil, fmt.Errorf("error reading the history from the file: %w", err)
	}

	return history, nil
}

// SaveStatistics persists the given chat statistics to the file system.
// This method generates a unique filename based on the ID of the chat statistics
// and writes the statistics to a JSON file within the BaseDir.
// If a file with the same name already exists, it will be overwritten.
//
// statistics: The chat statistics to be saved.
//
// Returns:
// error: An error if encountered during file operations or serialization.
func (fs *FS) SaveStatistics(_ context.Context, statistics *chat.Statistics) error {
	// Generate the path to save the statistics using the ID.
	filename := fmt.Sprintf(
		"statistics-%d-%d-%s.json",
		statistics.ID.User,
		statistics.ID.Chat,
		statistics.ID.Model,
	)
	path := filepath.Join(fs.BaseDir, filename)

	// Open or create the file.
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open or create the file: %w", err)
	}
	defer file.Close()

	// Write the statistics to the file in JSON format.
	err = statistics.Write(file)
	if err != nil {
		return fmt.Errorf("error writing the statistics to the file: %w", err)
	}

	return nil
}

// LoadStatistics retrieves a chat statistics from the file system using the provided ID.
// If the file does not exist, a new Statistics instance is returned.
// It constructs a file name from the ID and tries to open the corresponding JSON file
// within the BaseDir, then deserializes the file's content into a Statistics object.
//
// id: The ID of the chat statistics to be loaded.
//
// Returns:
// *Statistics: A pointer to the retrieved or newly created Statistics object.
// error: An error if encountered during file operations or deserialization, except for file not found error.
func (fs *FS) LoadStatistics(_ context.Context, id chat.ID) (*chat.Statistics, error) {
	// Generate the path to load the statistics using the ID.
	filename := fmt.Sprintf("statistics-%d-%d-%s.json", id.User, id.Chat, id.Model)
	path := filepath.Join(fs.BaseDir, filename)

	// Open the file.
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// If the file does not exist, return a new instance of chat.History.
			return &chat.Statistics{
				ID:          id,
				LastMessage: 0,
				Daily:       0,
				Monthly:     map[time.Month]chat.Cost{},
				Total:       0,
			}, nil
		}
		// For other errors, return an error.
		return nil, fmt.Errorf("could not open the file: %w", err)
	}
	defer file.Close()

	// Decode the statistics from the file.
	statistics := new(chat.Statistics)
	err = statistics.Read(file)
	if err != nil {
		return nil, fmt.Errorf("error reading the statistics from the file: %w", err)
	}

	return statistics, nil
}
