package chat

import "context"

// SessionProvider is an interface for managing and providing access to user chat sessions.
// It defines a method for obtaining a Session based on a unique identifier.
type SessionProvider interface {
	// ProvideSession retrieves or creates a chat session for the given user ID.
	// It accepts a context for deadline control and cancellation and an ID representing
	// the unique identifier for the user's chat session.
	//
	// ctx: The context for the operation, which allows for deadline control and cancellation.
	// id: The unique identifier for the desired chat session.
	//
	// Returns a Session corresponding to the user ID and an error if the retrieval or creation fails.
	ProvideSession(ctx context.Context, id ID) (Session, error)
}
