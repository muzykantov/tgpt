package chatgpt

import (
	"context"
	"fmt"
	"sync"

	"github.com/muzykantov/tgpt/chat"
	"github.com/sashabaranov/go-openai"
)

// ensure that the concrete type Session implements the chat.Session interface
var _ chat.Session = (*Session)(nil)

// Session encapsulates the state and management of a ChatGPT session.
// It includes a unique session identifier, an OpenAI client for interactions,
// a storage mechanism for persisting session data, and a session-specific cache.
type Session struct {
	chat.ID // Embedding chat.ID provides the unique identifiers for the user and the session.

	client  *openai.Client // client is the OpenAI client used to interface with the GPT API.
	storage chat.Storage   // storage is the abstract storage layer for saving and loading history and statistics.
	params  RequestParams  // params holds the parameters used to customize the OpenAI request.

	cache *sessionCache // cache holds the session's history and statistics to minimize storage access.
	mu    *sync.RWMutex // cacheMu is a read/write mutex for thread-safe access to the fields.
}

// NewSession creates a new chat Session with default request parameters.
//
// id: A composite identifier that includes the User, Chat, and Model information.
// client: An instance of the OpenAI Client.
// storage: An abstraction for the storage backend where session data is saved.
//
// Returns:
// A pointer to a new Session instance.
func NewSession(id chat.ID, client *openai.Client, storage chat.Storage) *Session {
	return &Session{
		ID:      id,
		client:  client,
		storage: storage,
		params:  DefaultRequestParams,
		cache:   &sessionCache{},
		mu:      &sync.RWMutex{},
	}
}

// SetRequestParams updates the session's request parameters with the provided values.
//
// params: The new request parameters to be used for subsequent OpenAI requests.
func (s *Session) SetRequestParams(params RequestParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.params = params
}

// SetPrompt updates the session's prompt with the provided string and persists the updated history.
// It locks the session for exclusive write access to prevent concurrent read/write issues.
// The method first ensures that the session's cache is loaded and then proceeds to update
// the prompt in the session's history cache before saving it to the storage.
//
// ctx: The context in which the network operations will be made, it allows for
// deadlines and cancellation signals to be carried through the function calls.
//
// prompt: The new prompt string to be set for the session.
//
// Returns:
// error: An error if encountered during the cache loading or while saving the history.
func (s *Session) SetPrompt(ctx context.Context, prompt string) error {
	s.mu.Lock()         // Lock the session for exclusive access.
	defer s.mu.Unlock() // Ensure the session is unlocked after this function returns.

	// Load the session cache if it's not already loaded.
	if err := s.loadCacheIfNeeded(ctx); err != nil {
		return err
	}

	// Don't need to update the same prompt.
	if s.cache.History.Prompt == prompt {
		return nil
	}

	s.cache.History.Prompt = prompt // Update the prompt in the history cache.

	// Persist the updated history to the storage.
	if err := s.storage.SaveHistory(ctx, s.cache.History); err != nil {
		return fmt.Errorf("error saving the history to the storage: %w", err)
	}

	return nil
}

// Ask sends a message to the OpenAI API and updates the session's history and statistics.
// The session's cache is loaded before making the request to ensure the latest data is used.
// If 'reset' is true, the history is cleared before sending the message; otherwise, the message
// is appended to the existing history. After receiving a response, the method calculates the cost,
// updates the history and statistics, and persists them to storage.
//
// ctx: The context in which the API call will be made. It may carry deadlines, cancellation signals,
// and other request-scoped values.
// message: The user message to send to the OpenAI API.
// reset: A flag indicating whether to reset the conversation history before sending the message.
//
// Returns:
// reply: The AI-generated response to the message.
// err: Any error encountered during the process. Errors may arise from loading cache, communicating
// with the OpenAI API, calculating costs, or persisting data to storage.
func (s *Session) Ask(ctx context.Context, message string, reset bool) (reply string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load session cache if necessary.
	if err := s.loadCacheIfNeeded(ctx); err != nil {
		return "", err
	}

	// Prepare the message history for the API request.
	msgs := make(
		[]openai.ChatCompletionMessage,
		0,
		len(s.cache.History.Log)+2,
	) // Prompt + Message + Log

	// Append the system prompt if available.
	if s.cache.History.Prompt != "" {
		msgs = append(msgs, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: s.cache.History.Prompt,
		})
	}

	// Append existing conversation if not resetting.
	if !reset {
		for _, msg := range s.cache.History.Log {
			msgs = append(msgs,
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleUser,
					Content: msg.User,
				},
				openai.ChatCompletionMessage{
					Role:    openai.ChatMessageRoleAssistant,
					Content: msg.Assistant,
				},
			)
		}
	}

	// Add the new user message to the history.
	msgs = append(msgs, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: message,
	})

	// Send the message to the OpenAI API.
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:            s.ID.Model,
		Messages:         msgs,
		MaxTokens:        s.params.MaxTokens,
		Temperature:      s.params.Temperature,
		TopP:             s.params.TopP,
		N:                1,
		Stream:           false,
		PresencePenalty:  s.params.PresencePenalty,
		FrequencyPenalty: s.params.FrequencyPenalty,
	})
	if err != nil {
		return "", fmt.Errorf("error creating chat completion: %w", err)
	}

	// Extract the AI's reply from the response.
	reply = resp.Choices[0].Message.Content

	// Calculate the cost of the interaction.
	usage := &Usage{
		Input:  resp.Usage.PromptTokens,
		Output: resp.Usage.CompletionTokens,
	}

	cost, err := usage.CalculateCostByModel(s.ID.Model)
	if err != nil {
		return "", fmt.Errorf("error calculating the cost: %w", err)
	}

	// Update the history and statistics unless we're resetting the history.
	if !reset {
		s.cache.History.Add(chat.Message{
			User:      message,
			Assistant: reply,
		})
	} else {
		s.cache.History.Clear()
	}

	s.cache.Statistics.AddCost(cost)

	// Persist the updated history and statistics.
	if err := s.storage.SaveHistory(ctx, s.cache.History); err != nil {
		return "", fmt.Errorf("error saving history to storage: %w", err)
	}

	if err := s.storage.SaveStatistics(ctx, s.cache.Statistics); err != nil {
		return "", fmt.Errorf("error saving statistics to storage: %w", err)
	}

	return reply, nil
}

// Reset clears the current session's chat history and updates the storage to reflect these changes.
// This method is protected by a mutex to ensure thread safety during the reset operation.
// It first loads the session cache if it is not already loaded, then clears the history,
// and finally persists the cleared history back to the storage.
//
// ctx: The context for the operation, which allows for deadline control and cancelation.
//
// Returns an error if the reset operation fails. This could happen if there is an issue loading
// the session cache, or if there is a problem saving the cleared history to storage. Errors from
// underlying operations are wrapped to provide more context.
func (s *Session) Reset(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Load session cache if necessary.
	if err := s.loadCacheIfNeeded(ctx); err != nil {
		return err
	}

	s.cache.History.Clear()

	// Persist the updated history and statistics.
	if err := s.storage.SaveHistory(ctx, s.cache.History); err != nil {
		return fmt.Errorf("error saving history to storage: %w", err)
	}

	return nil
}

// History returns a copy of the chat history from the session's cache.
// If the cache is not loaded, it attempts to load it before returning the history.
// This function ensures that any modifications to the returned History object
// do not affect the original data in the cache.
//
// ctx: The context for handling cancellation and timeouts.
//
// Returns:
// *chat.History: A copy of the chat history for the session.
// error: An error if encountered during cache loading.
func (s *Session) History(ctx context.Context) (*chat.History, error) {
	s.mu.RLock() // Use read lock to allow concurrent reads.
	defer s.mu.RUnlock()

	// Load session cache if necessary.
	if err := s.loadCacheIfNeeded(ctx); err != nil {
		return nil, err
	}

	// Return a clone of the history to prevent modifications to the cached version.
	return s.cache.History.Clone(), nil
}

// Statistics returns a copy of the chat statistics from the session's cache.
// If the cache is not loaded, it attempts to load it before returning the statistics.
// This function ensures that any modifications to the returned Statistics object
// do not affect the original data in the cache.
//
// ctx: The context for handling cancellation and timeouts.
//
// Returns:
// *chat.Statistics: A copy of the chat statistics for the session.
// error: An error if encountered during cache loading.
func (s *Session) Statistics(ctx context.Context) (*chat.Statistics, error) {
	s.mu.RLock() // Use read lock to allow concurrent reads.
	defer s.mu.RUnlock()

	// Load session cache if necessary.
	if err := s.loadCacheIfNeeded(ctx); err != nil {
		return nil, err
	}

	// Return a clone of the statistics to prevent modifications to the cached version.
	return s.cache.Statistics.Clone(), nil
}

// loadCacheIfNeeded checks if the session cache has been loaded and if not,
// loads the history and statistics from the storage.
//
// This method is unexported and meant to be used internally within Session methods.
//
// ctx: The context for controlling cancellations and timeouts.
//
// Returns:
// error: An error if any occurs during the loading process from the storage.
func (s *Session) loadCacheIfNeeded(ctx context.Context) (err error) {
	if s.cache.History != nil && s.cache.Statistics != nil {
		// Cache is already loaded, no need to load again
		return nil
	}

	// Load history and statistics from the storage backend.
	s.cache.History, err = s.storage.LoadHistory(ctx, s.ID)
	if err != nil {
		return err
	}

	s.cache.Statistics, err = s.storage.LoadStatistics(ctx, s.ID)
	if err != nil {
		return err
	}

	return nil
}

// sessionCache holds the temporary session data including history and statistics,
// along with metadata about the last update time and the duration it should be considered valid.
type sessionCache struct {
	*chat.History    // History holds the conversation history of the session.
	*chat.Statistics // Statistics holds the usage and performance statistics of the session.
}
