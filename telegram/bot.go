package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/muzykantov/tgpt/chat"
	"github.com/muzykantov/tgpt/lang"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Bot represents a Telegram bot that integrates advanced features such as
// session handling, user permissions, and administrative controls. The bot
// provides a mechanism for sending messages, manages sessions for different chats,
// enforces user access permissions, administers privilege controls, and supports
// localization of messages.
type Bot struct {
	// name of the chat bot.
	name string

	// sender is the mechanism for sending messages.
	sender Sender

	// session handles the session state for different chats.
	// It is a provider that manages chat sessions.
	session chat.SessionProvider

	// model is the name of the model for sessions.
	model string

	// allowedUsers specifies which users are permitted to interact with the bot.
	// It maps user IDs to empty structs, acting as a set to efficiently check user access.
	allowedUsers map[int64]struct{}

	// adminUsers specifies which users have administrative privileges.
	// It maps user IDs to empty structs, similar to allowedUsers,
	// to efficiently check for administrative access.
	adminUsers map[int64]struct{}

	// printer is used for localizing messages based on the provided language tag.
	// It facilitates internationalization by printing messages in the user's language.
	printer *message.Printer

	// adminContact holds the contact information for the bot administrator.
	// This could be used to provide a contact reference for users needing assistance.
	adminContact string

	// currency represents the type of currency used for transactions or
	// monetary values within the bot's functionality.
	currency string

	// rate defines the exchange or conversion rate associated with the
	// currency specified, which may be used in financial calculations.
	rate float64
}

// NewBot creates and initializes a new instance of Bot with the necessary dependencies.
// It configures the bot with mechanisms for message sending, session management,
// and optional user whitelists for access and administrative permissions. The printer
// for localized messages is also initialized based on the specified language.
//
// Parameters:
//   - name: Name of the chat bot.
//   - sender: The message sending mechanism complying with the Sender interface.
//   - sessionProvider: The provider for managing chat session states.
//   - model: The name of the model used for sessions.
//   - allowedUsers: An optional slice of user IDs permitted to interact with the bot.
//   - adminUsers: An optional slice of user IDs granted administrative privileges.
//   - language: The language.Tag used for localized message printing.
//   - adminContact: Contact information for the bot administrator.
//   - currency: A string representing the currency code (e.g., "$", "￥") used for statistics reporting.
//   - rate: A float64 value representing the exchange rate or a conversion factor used for financial statistics.

// Returns:
//   - A pointer to the newly created Bot instance.
func NewBot(
	name string,
	sender Sender,
	sessionProvider chat.SessionProvider,
	model string,
	allowedUsers, adminUsers []int64,
	language language.Tag,
	adminContact string,
	currency string,
	rate float64,
) *Bot {
	bot := &Bot{
		name:         name,
		sender:       sender,
		session:      sessionProvider,
		model:        model,
		allowedUsers: make(map[int64]struct{}),
		adminUsers:   make(map[int64]struct{}),
		printer:      message.NewPrinter(language),
		adminContact: adminContact,
		currency:     currency,
		rate:         rate,
	}

	// Populate the allowedUsers map
	for _, userID := range allowedUsers {
		bot.allowedUsers[userID] = struct{}{}
	}

	// Populate the adminUsers map
	for _, userID := range adminUsers {
		bot.allowedUsers[userID] = struct{}{}
		bot.adminUsers[userID] = struct{}{}
	}

	return bot
}

// ProcessUpdates listens for incoming updates from the Telegram bot API
// and processes each message update asynchronously.
//
// ctx: The context to control the lifecycle of the update processing. If the context
// is canceled, the method will stop processing updates and return.
//
// updates: A channel through which the Telegram bot API sends updates.
//
// Returns:
// - An error if the context is canceled, otherwise runs indefinitely without returning.
func (b *Bot) ProcessUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case update := <-updates:
			if update.Message == nil { // Ignore any non-Message updates.
				continue
			}

			go b.handleMessage(ctx, update.Message)
		}
	}
}

// Reply sends a textual reply to a specific message within a Telegram chat.
// It constructs a message configuration targeting the original message
// and uses the Bot's sender to dispatch the reply. If an error occurs
// during the sending process, it is logged with the context of the chat and
// the message that was attempted to be replied to.
//
// Parameters:
//
//	to   - A pointer to the Message instance to which the reply is directed.
//	       This must include the necessary identifiers to target the correct message.
//	with - The text content to be sent as the reply message.
//
// No return values, but errors during message sending are logged.
func (b *Bot) Reply(to *tgbotapi.Message, with string) {
	// Creating a message configuration for replying to the specific message.
	msg := tgbotapi.NewMessage(to.Chat.ID, with)
	msg.ReplyToMessageID = to.MessageID
	msg.ParseMode = "markdown"

	// Using the Send method of the sender to dispatch the message.
	_, err := b.sender.Send(msg)
	if err != nil {
		slog.Error(
			"reply error",
			slog.Int64("chatID", to.Chat.ID),
			slog.Int("messageID", to.MessageID),
			slog.String("error", err.Error()),
		)
	}
}

// Send dispatches a non-reply message to a specified chat in Telegram.
// It creates a new message configuration with the designated chat ID and
// message content, then sends it using the Bot's sender. If the message
// fails to Send, an error is logged including the chat ID and the error message.
//
// Parameters:
//
//	chat    - The chat ID to which the message should be sent.
//	message - The text content of the message to be sent.
//
// No return values, but errors during message sending are logged.
func (b *Bot) Send(chat int64, message string) {
	// Creating a message configuration.
	msg := tgbotapi.NewMessage(chat, message)
	msg.ParseMode = "markdown"

	// Using the Send method of the sender to dispatch the message.
	_, err := b.sender.Send(msg)
	if err != nil {
		slog.Error(
			"send error",
			slog.Int64("chatID", chat),
			slog.String("error", err.Error()),
		)
	}
}

// Typing simulates typing activity in a chat until the provided context is cancelled.
//
// Parameters:
//
//	ctx    - The context that controls the cancellation of the typing action.
//	chatID - The ID of the chat where the typing action will be shown.
//
// The function does not return a value. Errors encountered during sending the typing action
// are logged but not returned.
func (b *Bot) Typing(ctx context.Context, chatID int64) {
	typingCfg := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
	b.sender.Send(typingCfg)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// The context has been cancelled, stop the typing action.
			return

		case <-ticker.C:
			// Send a typing action to the chat.
			b.sender.Send(typingCfg)
		}
	}
}

// IsUserAllowed checks if the user with the given ID is allowed to interact with the bot.
//
// userID: The Telegram user ID to check for permission.
//
// Returns:
// - true if the user is allowed, false otherwise.
func (b *Bot) IsUserAllowed(userID int64) bool {
	_, allowed := b.allowedUsers[userID]
	return allowed
}

// IsUserAdmin checks if the user with the given ID has administrative privileges.
//
// userID: The Telegram user ID to check for administrative status.
//
// Returns:
// - true if the user is an admin, false otherwise.
func (b *Bot) IsUserAdmin(userID int64) bool {
	_, admin := b.adminUsers[userID]
	return admin
}

// handleMessage processes a received Telegram message.
// Depending on the content and the sender of the message, it will perform different actions.
// This method is designed to be called as a goroutine to handle each message concurrently.
//
// ctx: The context to control the lifecycle of the message processing. If the context
// is canceled, the method should cease processing and return.
//
// msg: The Telegram message to process.
func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	start := time.Now()

	slog.Info(
		"handleMessage started",
		slog.Int64("chatID", msg.Chat.ID),
		slog.Int("messageID", msg.MessageID),
		slog.String("messageText", msg.Text),
		slog.Time("started", start),
	)

	defer func() {
		end := time.Now()
		slog.Info(
			"handleMessage finished",
			slog.Int64("chatID", msg.Chat.ID),
			slog.Int("messageID", msg.MessageID),
			slog.String("messageText", msg.Text),
			slog.Time("finished", end),
			slog.Duration("elapsed", end.Sub(start)),
		)
	}()

	// First, check if the user or admin is allowed to interact with the bot.
	if !b.IsUserAllowed(msg.From.ID) {
		b.Reply(msg, b.printer.Sprintf(lang.MsgNotAllowed, msg.From.ID, b.adminContact))
		return
	}

	if msg.IsCommand() {
		// Handle the command.
		b.handleCommand(ctx, msg)
	} else {
		// Handle a regular message.
		b.handleRegularMessage(ctx, msg)
	}
}

// handleCommand processes a command received in a message.
// Actual implementation will depend on the specific commands your bot supports.
//
// ctx: The context for controlling the processing lifecycle.
// msg: The message containing the command to process.
func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	session, err := b.session.ProvideSession(ctx, chat.ID{
		User:  msg.From.ID,
		Chat:  msg.Chat.ID,
		Model: b.model,
	})
	if err != nil {
		b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
		slog.Error(
			"handleCommand ProvideSession error",
			slog.Int64("chatID", msg.Chat.ID),
			slog.Int("messageID", msg.MessageID),
			slog.String("messageText", msg.Text),
			slog.String("error", err.Error()),
		)
		return
	}

	commands := []tgbotapi.BotCommand{
		{Command: "help", Description: b.printer.Sprintf(lang.MsgCommandHelp)},
		{Command: "stats", Description: b.printer.Sprintf(lang.MsgCommandStats)},
		{Command: "restart", Description: b.printer.Sprintf(lang.MsgCommandRestart)},
	}

	switch msg.Command() {
	case "start":
		if _, err := b.sender.Request(tgbotapi.NewSetMyCommands(commands...)); err != nil {
			b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
			slog.Error(
				"handleCommand Request error",
				slog.Int64("chatID", msg.Chat.ID),
				slog.Int("messageID", msg.MessageID),
				slog.String("messageText", msg.Text),
				slog.String("error", err.Error()),
			)
		}

		fallthrough

	case "help":
		sb := &strings.Builder{}
		sb.WriteString(b.printer.Sprintf(lang.MsgGreeting, b.name))
		for _, cmd := range commands {
			sb.WriteString(
				fmt.Sprintf("/%s — %s\n\n", cmd.Command, cmd.Description),
			)
		}
		sb.WriteString(b.printer.Sprintf(lang.MsgSupport, b.adminContact))
		b.Send(msg.Chat.ID, sb.String())

	case "restart":
		if err := session.Reset(ctx); err != nil {
			b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
			slog.Error(
				"handleCommand Reset error",
				slog.Int64("chatID", msg.Chat.ID),
				slog.Int("messageID", msg.MessageID),
				slog.String("messageText", msg.Text),
				slog.String("error", err.Error()),
			)
		}

		args := msg.CommandArguments()
		if args != "" {
			if err := session.SetPrompt(ctx, args); err != nil {
				b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
				slog.Error(
					"handleCommand SetPrompt error",
					slog.Int64("chatID", msg.Chat.ID),
					slog.Int("messageID", msg.MessageID),
					slog.String("messageText", msg.Text),
					slog.String("error", err.Error()),
				)
			}
		}

		b.Reply(msg, b.printer.Sprintf(lang.MsgDone))

	case "stats":
		stats, err := session.Statistics(ctx)
		if err != nil {
			b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
			slog.Error(
				"handleCommand ProvideSession error",
				slog.Int64("chatID", msg.Chat.ID),
				slog.Int("messageID", msg.MessageID),
				slog.String("messageText", msg.Text),
				slog.String("error", err.Error()),
			)
			return
		}

		now := chat.Now()
		b.Send(msg.Chat.ID, b.printer.Sprintf(
			lang.MsgStats,
			b.currency, b.rate*float64(stats.LastMessage),
			b.currency, b.rate*float64(stats.Daily),
			b.currency, b.rate*float64(stats.Monthly[now.Month()]),
			b.currency, b.rate*float64(stats.Total),
		))

	default:
		b.Reply(msg, b.printer.Sprintf(lang.MsgCommandNotSupported))
	}
}

// handleRegularMessage processes a standard message that is not a command.
// The processing logic can include sending replies, performing actions, etc.
//
// ctx: The context for controlling the processing lifecycle.
// msg: The message to process.
func (b *Bot) handleRegularMessage(ctx context.Context, msg *tgbotapi.Message) {
	start := time.Now()

	slog.Info(
		"handleRegularMessage started",
		slog.Int64("chatID", msg.Chat.ID),
		slog.Int("messageID", msg.MessageID),
		slog.String("messageText", msg.Text),
		slog.Time("started", start),
	)

	var replyText string
	defer func() {
		end := time.Now()
		slog.Info(
			"handleRegularMessage finished",
			slog.Int64("chatID", msg.Chat.ID),
			slog.Int("messageID", msg.MessageID),
			slog.String("replyText", replyText),
			slog.Time("finished", end),
			slog.Duration("elapsed", end.Sub(start)),
		)
	}()

	if msg.Text == "" {
		b.Reply(msg, b.printer.Sprintf(lang.MsgNotSupported))
		return
	}

	session, err := b.session.ProvideSession(ctx, chat.ID{
		User:  msg.From.ID,
		Chat:  msg.Chat.ID,
		Model: b.model,
	})
	if err != nil {
		b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
		slog.Error(
			"handleRegularMessage ProvideSession error",
			slog.Int64("chatID", msg.Chat.ID),
			slog.Int("messageID", msg.MessageID),
			slog.String("messageText", msg.Text),
			slog.String("error", err.Error()),
		)
		return
	}

	typingCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go b.Typing(typingCtx, msg.Chat.ID)

	reply, err := session.Ask(ctx, msg.Text, false)
	if err != nil {
		b.Reply(msg, b.printer.Sprintf(lang.MsgUnexpectedError, b.adminContact, err.Error()))
		slog.Error(
			"handleRegularMessage Ask error",
			slog.Int64("chatID", msg.Chat.ID),
			slog.Int("messageID", msg.MessageID),
			slog.String("messageText", msg.Text),
			slog.String("error", err.Error()),
		)
		return
	}

	b.Reply(msg, reply)
	replyText = reply
}
