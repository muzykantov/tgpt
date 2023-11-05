package chat

import "context"

// Storage defines an interface for managing the persistence of chat history and statistics.
// It provides an abstraction over the actual storage mechanism, which could be implemented
// using various systems such as files, databases, or other storage backends.
type Storage interface {
	// SaveHistory persists a given chat history into the storage.
	// The method ensures that the provided History object is stored and retrievable
	// by an identifier. In the event of a failure during the save operation,
	// an error will be returned.
	//
	// ctx: A context.Context to allow for cancellation and timeout control during the save process.
	// history: The History object containing chat history data to be saved.
	//
	// Returns an error if the save operation encounters issues.
	SaveHistory(ctx context.Context, history *History) error

	// LoadHistory retrieves the chat history associated with the given ID from storage.
	// If no history is associated with the ID, a new, empty History object is returned.
	// This method does not consider the absence of a history record as an error condition.
	//
	// ctx: A context.Context to allow for cancellation and timeout control during the load process.
	// id: The unique identifier used to retrieve the chat history.
	//
	// Returns the retrieved or new History object, and an error if the load operation fails
	// for reasons other than the history not being found.
	LoadHistory(ctx context.Context, id ID) (*History, error)

	// SaveStatistics persists the given chat statistics into the storage.
	// This method ensures that the provided Statistics object is stored and retrievable
	// by an identifier. Should the save operation encounter a failure, an error is returned.
	//
	// ctx: A context.Context to allow for cancellation and timeout control during the save process.
	// statistics: The Statistics object containing chat statistics data to be saved.
	//
	// Returns an error if the save operation encounters issues.
	SaveStatistics(ctx context.Context, statistics *Statistics) error

	// LoadStatistics retrieves the chat statistics associated with the given ID from storage.
	// If no statistics are associated with the ID, a new, empty Statistics object is returned.
	// Similar to LoadHistory, this method does not consider the absence of statistics as an error condition.
	//
	// ctx: A context.Context to allow for cancellation and timeout control during the load process.
	// id: The unique identifier used to retrieve the chat statistics.
	//
	// Returns the retrieved or new Statistics object, and an error if the load operation fails
	// for reasons other than the statistics not being found.
	LoadStatistics(ctx context.Context, id ID) (*Statistics, error)
}
