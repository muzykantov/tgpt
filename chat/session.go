package chat

import "context"

// Session is an interface that abstracts the operations of a chat session.
// It defines the contract for a session that can send messages, set prompts,
// and retrieve history and statistics.
type Session interface {
	// Ask takes a context, a message string, and a reset flag as inputs. It sends
	// the message to an underlying chat service and returns the service's reply.
	// The reset flag indicates whether to start a new conversation (clearing history) or
	// continue with the existing one.
	//
	// ctx: The context for the API call, which allows for deadline control and cancelation.
	// message: The message string to send to the chat service.
	// reset: A boolean flag indicating whether to clear the conversation history.
	//
	// Returns the reply from the chat service as a string and an error if the operation fails.
	Ask(ctx context.Context, message string, reset bool) (reply string, err error)

	// Reset terminates the current session and clears any saved state or history associated with it.
	// This function is intended to restart the session as if it were new, without any memory of previous interactions.
	//
	// ctx: The context for the operation, which allows for deadline control and cancelation.
	//
	// Returns an error if the reset operation fails, for instance, due to a failure in underlying storage
	// or network systems.
	Reset(ctx context.Context) error

	// SetPrompt updates the prompt for the session to the given string. It affects the
	// conversation flow and can be used to provide context or instructions that persist across
	// exchanges in the session.
	//
	// ctx: The context for the operation, which allows for deadline control and cancelation.
	// prompt: The new prompt string to be used for subsequent interactions.
	//
	// Returns an error if the operation fails.
	SetPrompt(ctx context.Context, prompt string) error

	// History retrieves a copy of the chat history associated with the session.
	// The history reflects all messages sent and received during the session's lifecycle.
	//
	// ctx: The context for the operation, which allows for deadline control and cancelation.
	//
	// Returns a pointer to a History object containing the session's chat history and an error if the operation fails.
	History(ctx context.Context) (*History, error)

	// Statistics retrieves a copy of the chat statistics associated with the session.
	// Statistics may include metrics like the number of messages exchanged, word counts,
	// and other relevant data points that characterize the session's usage.
	//
	// ctx: The context for the operation, which allows for deadline control and cancelation.
	//
	// Returns a pointer to a Statistics object containing the session's chat statistics and an error if the operation fails.
	Statistics(ctx context.Context) (*Statistics, error)
}
