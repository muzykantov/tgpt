package chat

import (
	"encoding/json"
	"io"
	"time"
)

// Now is a variable function that returns the current time. By default, it is
// initialized to use UTC time. This can be overridden for testing or other
// purposes where control over time is required.
var Now func() time.Time

// init initializes the package variables. It sets the Now function to return
// the current UTC time.
func init() {
	Now = time.Now().UTC
}

// Cost represents the cost associated with a chat operation. It is used to
// track costs for messages, daily operations, monthly aggregates, and the total
// cost across the chat session's lifetime.
type Cost float64

// Statistics contains data related to the cost and usage of chat sessions.
// It embeds the ID type to associate these statistics with a particular chat session.
type Statistics struct {
	ID                              // Embedded ID to uniquely identify the chat session.
	LastMessage Cost                // LastMessage is the cost of the last message in the chat session.
	Daily       Cost                // Daily is the total cost of the chat session for the current day.
	Monthly     map[time.Month]Cost // Monthly is a map tracking the cost per month.
	Total       Cost                // Total is the cumulative cost of the chat session.
	LastUpdate  time.Time           // LastUpdate records the timestamp of the last time the Statistics were modified.
}

// AddCost updates the Statistics instance with a new cost from a chat interaction.
// It checks if the current day is different from the last update day and resets
// the daily cost to zero if a new day has started. Then, it adds the new cost to
// the last message, daily, monthly, and total costs. For monthly tracking, if there
// is no entry for the current month, it creates one. It also updates the last update
// time to the current time.
//
// The method assumes there is a LastUpdate field of type time.Time in the Statistics
// structure to keep track of when the statistics were last updated.
//
// newCost: The cost from the new chat interaction to add to the statistics.
func (s *Statistics) AddCost(newCost Cost) {
	now := Now()

	// Check if a new day has started.
	if now.Day() != s.LastUpdate.Day() || now.Month() != s.LastUpdate.Month() || now.Year() != s.LastUpdate.Year() {
		s.Daily = 0
	}

	s.LastMessage = newCost
	s.Daily += newCost
	s.Total += newCost

	currentMonth := now.Month()

	if s.Monthly == nil {
		s.Monthly = make(map[time.Month]Cost)
	}

	s.Monthly[currentMonth] += newCost
	s.LastUpdate = now
}

// Write serializes the Statistics instance and writes it to the provided io.Writer in JSON format.
//
// w: The writer to which the serialized statistics should be written.
//
// Returns:
// error: An error if encountered during the serialization or writing process.
func (s *Statistics) Write(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	return encoder.Encode(s)
}

// Read deserializes the Statistics instance from the provided io.Reader which should contain
// the statistics in JSON format.
//
// r: The reader from which the serialized statistics should be read.
//
// Returns:
// error: An error if encountered during the deserialization process.
func (s *Statistics) Read(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(s)
}

// Clone creates a deep copy of the Statistics object. This is particularly useful
// when you want to duplicate a Statistics object to make thread-safe operations
// without affecting the original object.
//
// Returns:
// *Statistics: A new instance of Statistics which is a deep copy of the original.
func (s *Statistics) Clone() *Statistics {
	// Create a new Statistics object with shallow-copied fields.
	clone := &Statistics{
		ID:          s.ID, // ID can be shallow copied as it contains only primitive types.
		LastMessage: s.LastMessage,
		Daily:       s.Daily,
		Total:       s.Total,
	}

	// Make a deep copy of the Monthly map to ensure independent manipulation.
	clone.Monthly = make(map[time.Month]Cost, len(s.Monthly))
	for k, v := range s.Monthly {
		clone.Monthly[k] = v
	}

	return clone
}
