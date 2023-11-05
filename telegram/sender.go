package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Sender defines an interface for sending messages through Telegram's API.
// Implementers of Sender must be able to handle the sending of various types
// of chattable content, such as text messages, photos, audio, etc.
type Sender interface {
	// Send transmits the provided chattable content to a Telegram chat.
	// This method should implement the necessary logic to handle the transmission
	// of different types of messages encapsulated by the Chattable interface.
	//
	// Parameters:
	//   - c: A Chattable instance containing the message or content to be sent.
	//        Chattable is an interface that can represent any sendable content.
	//
	// Returns:
	//   - Message: The Telegram Message object that was sent. This object contains
	//              details about the message, such as its ID, the chat it was sent to,
	//              and the time it was sent.
	//   - error: An error encountered during the sending process. If the message was
	//            sent successfully, this will be nil.
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)

	// Request sends the chattable content to Telegram's API and returns the raw API response.
	// This method is typically used for sending messages and receiving the direct response
	// from the API without any additional processing.
	//
	// Parameters:
	//   - c: A Chattable instance containing the message or content to be sent to the Telegram API.
	//        As with Send, this represents any type of content that conforms to the Chattable interface.
	//
	// Returns:
	//   - APIResponse: A pointer to the APIResponse from Telegram. This response includes
	//                  the raw status, the result, and any error message from the API.
	//   - error: An error encountered while making the request to Telegram's API. If the
	//            request was successful, this will be nil.
	Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
}
