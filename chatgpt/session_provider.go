package chatgpt

import (
	"context"
	"sync"
	"time"

	"github.com/muzykantov/tgpt/chat"
	"github.com/sashabaranov/go-openai"
)

// ensure that the concrete type Session implements the chat.Session interface
var _ chat.SessionProvider = (*SessionProvider)(nil)

// SessionProvider maintains a cache of chat sessions with a defined time-to-live (TTL) and
// manages their lifecycle. It uses a cleanup interval to periodically remove expired sessions.
// The SessionProvider initializes sessions with OpenAI Client and a storage interface,
// and controls concurrent access using a mutex.
type SessionProvider struct {
	// sessions keep track of all active chat sessions by their unique IDs.
	sessions map[chat.ID]*sessionInfo

	// ttl defines the time-to-live for sessions in the cache before they expire.
	ttl time.Duration

	// client is the OpenAI client used to interface with the GPT API for chat interactions.
	client *openai.Client

	// storage is the abstract storage layer for session data persistence.
	storage chat.Storage

	// params hold the parameters used to customize the OpenAI request.
	params RequestParams

	// mu provides concurrency control for accessing the sessions map.
	mu sync.RWMutex

	// cleanupInterval specifies the frequency at which expired sessions are cleaned up.
	cleanupInterval time.Duration
}

// NewSessionProvider creates a new chat.Provider with the specified OpenAI client, storage interface,
// request parameters, TTL for sessions, and cleanup interval.
// It starts a background goroutine to periodically clean up expired sessions.
//
// client: Instance of the OpenAI Client for API interactions.
// storage: Storage backend for session data persistence.
// params: Default request parameters for GPT API interactions.
// ttl: Time-to-live for sessions to determine their expiration.
// cleanupInterval: Interval to execute cleanup of expired sessions.
//
// Returns a pointer to a newly created Provider.
func NewSessionProvider(
	client *openai.Client,
	storage chat.Storage,
	params RequestParams,
	ttl, cleanupInterval time.Duration,
) *SessionProvider {
	sm := &SessionProvider{
		sessions:        make(map[chat.ID]*sessionInfo),
		ttl:             ttl,
		client:          client,
		storage:         storage,
		params:          params,
		cleanupInterval: cleanupInterval,
	}

	go sm.cleanupScheduler()

	return sm
}

// GetOrCreateSession retrieves an existing session associated with the given ID from the session manager,
// or creates a new one if it does not exist. It ensures that only one session is created or retrieved
// at a time through mutual exclusion.
//
// id: The unique identifier for the chat session.
//
// Returns:
// chat.Session: The imlementation of retrieved or newly created session.
// error: An error if encountered during the session creation process.
func (m *SessionProvider) ProvideSession(_ context.Context, id chat.ID) (chat.Session, error) {
	return m.GetOrCreateSession(id)
}

// GetOrCreateSession retrieves an existing session associated with the given ID from the session manager,
// or creates a new one if it does not exist. It ensures that only one session is created or retrieved
// at a time through mutual exclusion.
//
// id: The unique identifier for the chat session.
//
// Returns:
// *Session: A pointer to the retrieved or newly created session.
// error: An error if encountered during the session creation process.
func (m *SessionProvider) GetOrCreateSession(id chat.ID) (*Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sInfo, exists := m.sessions[id] // Check if the session already exists.
	if !exists {
		// If the session does not exist, create a new session.
		newSession := NewSession(id, m.client, m.storage)
		newSession.SetRequestParams(m.params) // Set request parameters for the new session.
		sInfo = &sessionInfo{
			session:    newSession, // Assign the new session.
			lastAccess: chat.Now(), // Set the current time as the last access time.
		}
		m.sessions[id] = sInfo // Store the new session information in the map.
	} else {
		// If the session exists, update the last access time.
		sInfo.lastAccess = chat.Now()
	}

	return sInfo.session, nil // Return the session.
}

// Clear terminates all managed sessions within the SessionManager. This is done by clearing
// the session map which holds all active sessions.
//
// This function locks the SessionManager to ensure exclusive access to the sessions map
// and defers the unlocking until the operation is complete, to prevent any race conditions.
func (m *SessionProvider) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for id := range m.sessions {
		delete(m.sessions, id) // Remove each session from the map, effectively ending it.
	}
}

// cleanupScheduler triggers cleanupExpiredSessions at every cleanup interval to remove expired sessions.
func (m *SessionProvider) cleanupScheduler() {
	for {
		time.Sleep(m.cleanupInterval)
		m.cleanupExpiredSessions()
	}
}

// cleanupExpiredSessions iterates through the sessions and deletes any that have not been accessed within the TTL.
func (m *SessionProvider) cleanupExpiredSessions() {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	for id, sInfo := range m.sessions {
		if now.Sub(sInfo.lastAccess) > m.ttl {
			delete(m.sessions, id)
		}
	}
}

// sessionInfo holds the data for a session along with the last access timestamp.
// It is used internally by the SessionManager to manage session state and lifecycle.
type sessionInfo struct {
	// session is a pointer to the session instance.
	session *Session

	// lastAccess records the last time the session was accessed, used to determine session expiration.
	lastAccess time.Time
}
