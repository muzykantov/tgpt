# TGPT - Lightweight GPT Telegram Bot

TGPT is a minimalist, yet powerful GPT Telegram bot built in Go with a pure focus on functionality, speed, and efficiency. The bot doesn't require any database or extra dependencies, is lightweight enough that it can even run on lower-end hardware such as a Raspberry Pi Zero W.

## Main Features
- Fast and Efficient: Designed and built to work swiftly, ensuring optimal user experience.
- Minimal Requirements: Runs smoothly even on low-end hardware. No database or additional dependencies required.
- Accurate Cost Calculation: Precisely calculates cost, ensuring correct estimates.
- Versions Supported: Fully supports the GPT-3.5 Turbo and future-proof with GPT-4 support. Different context length can be handled.
- Chat History: Allows to maintain chat history, enabling continuity in user interactions.
- Light on Hardware: Among the unique advantages of TGPT is its low hardware requirements, making it easier to host and maintain than some other options.

## Installation
Clone the project into your Go directory.
 
git clone https://github.com/muzykantov/tgpt.git

Go to the project directory and build it. The binary file will be named tgpt.
cd tgpt
go build

## Usage
To run the bot, simply start the executable.
./tgpt

## Configuration

Before you can run the bot, you need to configure it by setting environment variables. These variables can either be set in your environment directly or by using a `.env` file in the root directory of the project.

Below is a list of all the environment variables used by the bot, along with a description for each:

### API Parameters (Required)

- `TGPT_TELEGRAM_BOT_TOKEN`: Your Telegram bot token obtained from BotFather.
- `TGPT_OPENAI_API_KEY`: Your OpenAI API key for accessing GPT models.

### Bot Parameters (Optional)

- `TGPT_NAME`: The name you want to give to your Telegram bot (default is "TGPT").
- `TGPT_MODEL`: The language model to use, default is "gpt-4".
- `TGPT_ALLOWED_USERS`: Comma-separated list of user IDs allowed to interact with the bot.
- `TGPT_ADMIN_USERS`: Comma-separated list of admin user IDs with extended permissions.
- `TGPT_LANGUAGE`: The language code for bot responses (default is "en").
- `TGPT_ADMIN_CONTACT`: The username or channel name of the admin for contact purposes.
- `TGPT_CURRENCY`: The currency symbol to use in financial interactions, e.g., for donations (default is "$").
- `TGPT_RATE`: The exchange rate used for converting currencies, if applicable (default is "1.0").

### OPENAI Client Parameters (Optional)

- `TGPT_CACHE_TTL_SEC`: Time-to-live for the cache, in seconds (default is "3600").
- `TGPT_DB_DIR`: The directory where the database files will be stored (default is ".db").
- `TGPT_MAX_TOKENS`: The maximum number of tokens the model should generate in each response.
- `TGPT_TEMPERATURE`: Controls the randomness in the model's output, with lower values leading to more deterministic responses.
- `TGPT_TOP_P`: Influences the range of token probabilities considered for generating each token in a response.
- `TGPT_PRESENCE_PENALTY`: Adjusts the model to prefer tokens from the input, which can encourage the model to talk about new topics.
- `TGPT_FREQUENCY_PENALTY`: Adjusts the model to avoid using tokens from the input, which can discourage the model from repeating itself.

### Setting Up the `.env` File

To use a `.env` file for your configuration:

1. Create a file named `.env` in the root directory of the project.
2. Copy the contents of the `.env.example` file (if provided) into the `.env` file.
3. Fill in the values for the environment variables.

**Important**: Never commit sensitive keys and tokens to version control. Always keep `TGPT_TELEGRAM_BOT_TOKEN` and `TGPT_OPENAI_API_KEY` confidential.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
TGPT is under the GPL license. See LICENSE for more information.

## Keywords
Go, Chatbot, GPT, Telegram Bot, Lightweight Chatbot, GPT-3.5 Turbo, GPT-4, Raspberry Pi, Minimalistic Chatbot, Accurate Cost Calculation, Fast and Efficient, Chat History Support

Crafted with care and efficiency. Always aimed at surpassing user's expectations. TGPT is the next step your Telegram Bot experience deserves.
