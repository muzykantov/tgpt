// Package lang manages language-specific messages for a Telegram chatbot.
// It provides localization support for messages sent by the bot to users.
//
// This package uses the golang.org/x/text/message package to manage and format
// localized messages.
package lang

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	MsgNotAllowed          = "Dear user with ID %d, unfortunately, you are not authorized to use this chatbot. To request access, please contact the administrator %s and provide your user ID."
	MsgUnexpectedError     = "An unexpected error occurred while processing your request. To resolve the issue, please forward this message to the bot administrator %s.\n\nError message: %s."
	MsgNotImplemented      = "Unfortunately, the feature has not been implemented yet."
	MsgNotSupported        = "This type of message is not supported."
	MsgCommandNotSupported = "This command is not supported."
	MsgDone                = "Done."
	MsgStats               = "*Cost statistics*```\nLast message: %s%.2f\nToday       : %s%.2f\nThis month  : %s%.2f\nAll-time    : %s%.2f```"
	MsgGreeting            = "*Welcome to the %s chatbot!*\n\nSend me a message to start a conversation or choose one of the available commands:\n\n"
	MsgSupport             = "For support inquiries, please contact %s."
	MsgCommandHelp         = "Show the help message."
	MsgCommandStats        = "Get usage statistics."
	// MsgCommandResend       = "Resend the last message."
	MsgCommandRestart = "Restart the conversation. Optionally, pass general instructions (for example, /restart you are a helpful assistant)."
)

func init() {
	message.SetString(language.AmericanEnglish, MsgNotAllowed, MsgNotAllowed)
	message.SetString(language.AmericanEnglish, MsgUnexpectedError, MsgUnexpectedError)
	message.SetString(language.AmericanEnglish, MsgNotImplemented, MsgNotImplemented)
	message.SetString(language.AmericanEnglish, MsgNotSupported, MsgNotSupported)
	message.SetString(language.AmericanEnglish, MsgCommandNotSupported, MsgCommandNotSupported)
	message.SetString(language.AmericanEnglish, MsgDone, MsgDone)
	message.SetString(language.AmericanEnglish, MsgStats, MsgStats)
	message.SetString(language.AmericanEnglish, MsgGreeting, MsgGreeting)
	message.SetString(language.AmericanEnglish, MsgSupport, MsgSupport)
	message.SetString(language.AmericanEnglish, MsgCommandHelp, MsgCommandHelp)
	message.SetString(language.AmericanEnglish, MsgCommandStats, MsgCommandStats)
	// message.SetString(language.AmericanEnglish, MsgCommandResend, MsgCommandResend)
	message.SetString(language.AmericanEnglish, MsgCommandRestart, MsgCommandRestart)

	message.SetString(language.Russian, MsgNotAllowed, "Уважаемый пользователь с ID %d, к сожалению, у вас нет доступа к использованию этого чат-бота. Чтобы запросить доступ, пожалуйста, свяжитесь с администратором %s и предоставьте ваш ID пользователя.")
	message.SetString(language.Russian, MsgUnexpectedError, "Произошла неожиданная ошибка при обработке вашего запроса. Для устранения проблемы, пожалуйста, перешлите это сообщение администратору бота %s.\n\nСообщение об ошибке: %s.")
	message.SetString(language.Russian, MsgNotImplemented, "К сожалению, функция еще не реализована.")
	message.SetString(language.Russian, MsgNotSupported, "Этот тип сообщения не поддерживается.")
	message.SetString(language.Russian, MsgCommandNotSupported, "Эта команда не поддерживается.")
	message.SetString(language.Russian, MsgDone, "Готово.")
	message.SetString(language.Russian, MsgStats, "*Статистика расходов*```\nПоследнее сообщ.: %s%.2f\nЗа сегодня      : %s%.2f\nВ этом месяце   : %s%.2f\nЗа все время    : %s%.2f```")
	message.SetString(language.Russian, MsgGreeting, "*Вас приветствует %s чат-бот!*\n\nОтправь мне сообщение для начала беседы или выбери одну из доступных команд:\n\n")
	message.SetString(language.Russian, MsgSupport, "По вопросам поддержки, пожалуйста, обращайтесь к %s.")
	message.SetString(language.Russian, MsgCommandHelp, "Показать справочное сообщение.")
	message.SetString(language.Russian, MsgCommandStats, "Получить статистику использования.")
	// message.SetString(language.Russian, MsgCommandResend, "Повторная отправка последнего сообщения.")
	message.SetString(language.Russian, MsgCommandRestart, "Перезагрузить разговор. По желанию передай общие инструкции (например, /reset ты полезный помощник).")
}
