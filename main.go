package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/muzykantov/tgpt/chatgpt"
	"github.com/muzykantov/tgpt/storage"
	"github.com/muzykantov/tgpt/telegram"
	openai "github.com/sashabaranov/go-openai"
	lang "golang.org/x/text/language"

	// Autoload package which will read in .env on import.
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var (
		telegramBotToken = getEnv("TGPT_TELEGRAM_BOT_TOKEN", "")
		openaiApiKey     = getEnv("TGPT_OPENAI_API_KEY", "")

		name         = getEnv("TGPT_NAME", "TGPT")
		model        = getEnv("MODEL", "gpt-4")
		allowedUsers = getEnvAsSlice("TGPT_ALLOWED_USERS", []int64{}, ",")
		adminUsers   = getEnvAsSlice("TGPT_ADMIN_USERS", []int64{}, ",")
		language     = getEnv("TGPT_LANGUAGE", "us")
		adminContact = getEnv("TGPT_ADMIN_CONTACT", "https://github.com/muzykantov/tgpt")
		currency     = getEnv("TGPT_CURRENCY", "TGPT")
		rate         = getEnvAsFloat("TGPT_RATE", 1.0)

		cacheTTL         = time.Duration(getEnvAsInt("TGPT_CACHE_TTL_SEC", 3600)) * time.Second
		dbDir            = getEnv("TGPT_DB_DIR", ".db")
		maxTokens        = getEnvAsInt("TGPT_MAX_TOKENS", chatgpt.DefaultRequestParams.MaxTokens)
		temperature      = getEnvAsFloat32("TGPT_TEMPERATURE", chatgpt.DefaultRequestParams.Temperature)
		topP             = getEnvAsFloat32("TGPT_TOP_P", chatgpt.DefaultRequestParams.TopP)
		presencePenalty  = getEnvAsFloat32("TGPT_PRESENCE_PENALTY", chatgpt.DefaultRequestParams.PresencePenalty)
		frequencyPenalty = getEnvAsFloat32("TGPT_FREQUENCY_PENALTY", chatgpt.DefaultRequestParams.FrequencyPenalty)
	)

	fmt.Printf("Bot '%s' is starting...\n", name)

	var (
		tgClient     = must(tgbotapi.NewBotAPI(telegramBotToken))
		openaiClient = openai.NewClient(openaiApiKey)
	)

	// Parse the language tag
	langTag, err := lang.Parse(language)
	if err != nil {
		fmt.Printf("Error parsing language tag: %v\n", err)
		langTag = lang.English
	}

	tgpt := telegram.NewBot(
		name,
		tgClient,
		chatgpt.NewSessionProvider(
			openaiClient,
			&storage.FS{
				BaseDir: dbDir,
			},
			chatgpt.RequestParams{
				MaxTokens:        maxTokens,
				Temperature:      temperature,
				TopP:             topP,
				PresencePenalty:  presencePenalty,
				FrequencyPenalty: frequencyPenalty,
			},
			cacheTTL,
			cacheTTL/2,
		),
		model,
		allowedUsers,
		adminUsers,
		langTag,
		adminContact,
		currency,
		rate,
	)

	// Setup a channel to listen for interrupt signal (Ctrl+C) and SIGTERM.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start processing updates in a separate goroutine.
	go func() {
		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60

		updates := tgClient.GetUpdatesChan(u)

		if err := tgpt.ProcessUpdates(ctx, updates); err != nil {
			// Handle the error according to your application's needs.
			fmt.Println("Error processing updates:", err)
			cancel() // Signal the context to cancel.
		}
	}()

	// Print a message indicating that the bot has started and is ready to receive updates.
	fmt.Println("Bot is now running. Press Ctrl+C to exit.")

	// Block until a signal is received.
	sig := <-sigChan
	fmt.Printf("\nReceived %v, initiating graceful shutdown.\n", sig)

	// Cancel the context to signal any ongoing processes to finish.
	cancel()

	// Wait for a moment. We need to give some operations a chance to finish.
	time.Sleep(time.Second * 5)

	fmt.Println("Shutdown complete.")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultValue []int64, separator string) []int64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	var slice []int64
	for _, str := range strings.Split(valStr, separator) {
		if i, err := strconv.ParseInt(str, 10, 64); err == nil {
			slice = append(slice, i)
		} else {
			fmt.Printf("Error parsing slice from env var '%s': %v\n", key, err)
		}
	}
	return slice
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	if value, err := strconv.ParseFloat(valStr, 64); err == nil {
		return value
	} else {
		fmt.Printf("Error parsing float from env var '%s': %v\n", key, err)
		return defaultValue
	}
}

func getEnvAsInt(key string, defaultValue int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	if value, err := strconv.Atoi(valStr); err == nil {
		return value
	} else {
		fmt.Printf("Error parsing int from env var '%s': %v\n", key, err)
		return defaultValue
	}
}

func getEnvAsFloat32(key string, defaultValue float32) float32 {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultValue
	}

	if value, err := strconv.ParseFloat(valStr, 32); err == nil {
		return float32(value)
	} else {
		fmt.Printf("Error parsing float32 from env var '%s': %v\n", key, err)
		return defaultValue
	}
}

func must[T any](result T, err error) T {
	if err != nil {
		panic(err)
	}

	return result
}
