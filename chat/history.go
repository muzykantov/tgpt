// Package chat provides structures and functionalities to manage and persist chat
// histories. This includes the representation of individual chat messages, usage
// details, and methods to serialize and deserialize chat histories.
package chat

import (
	"encoding/json"
	"io"
)

// ID uniquely identifies a chat session. It consists of a user ID, chat session ID,
// and the model name used for that chat session.
type ID struct {
	User  int64  // User is the unique identifier for the user.
	Chat  int64  // Chat is the unique identifier for the chat session.
	Model string // Model represents the name of the model used for this chat.
}

// Message encapsulates a single chat interaction, including the message from the
// user and the corresponding response from the assistant.
type Message struct {
	User      string // User is the message provided by the user.
	Assistant string // Assistant is the response from the assistant.
}

// History captures the details of a chat session, including its unique ID,
// initial prompt, conversation log, and total token usage.
type History struct {
	ID               // ID is the unique identifier for this chat session.
	Prompt string    // Prompt is the initial statement or question that started the chat.
	Log    []Message // Log maintains a sequential record of the chat interactions.
}

// Add includes a new chat interaction to the history and updates the total token usage.
//
// msg: The new chat message to be appended to the log.
func (h *History) Add(msg Message) {
	h.Log = append(h.Log, msg)
}

// Clear removes all entries from the conversation log in the chat session history.
// This method resets the log to an empty state without modifying the initial prompt
// or the unique session ID.
func (h *History) Clear() {
	h.Log = []Message{}
}

// Write serializes the chat history and writes it to the provided io.Writer in JSON format.
//
// w: The writer to which the serialized history should be written.
//
// Returns:
// error: An error if encountered during the serialization or writing process.
func (h *History) Write(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(h)
}

// Read deserializes the chat history from the provided io.Reader which should contain
// the chat history in JSON format.
//
// r: The reader from which the serialized history should be read.
//
// Returns:
// error: An error if encountered during the deserialization process.
func (h *History) Read(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(h)
}

// Clone creates a deep copy of the History object. This method ensures that the
// returned History instance is a complete replica of the original, including all
// messages in the chat log. This is useful when you need to work with a copy of the
// History without affecting the original instance, maintaining thread safety.
//
// Returns:
// *History: A new instance of History that is a deep copy of the original.
func (h *History) Clone() *History {
	// Create a new History object with the ID and Prompt copied from the original.
	clone := &History{
		ID:     h.ID,     // ID can be copied directly as it is composed of primitive types.
		Prompt: h.Prompt, // String is immutable in Go, safe to directly assign.
	}

	// Make a deep copy of the Log slice to ensure independent manipulation.
	clone.Log = make([]Message, len(h.Log))
	copy(clone.Log, h.Log)

	return clone
}
