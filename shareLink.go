package main

import (
	tg "gopkg.in/telebot.v3"
	"strings"
)

func registerShareLinkHandlers(bot *tg.Bot) {
	bot.Handle(&shareBtn, shareLink)
}

func shareLink(c tg.Context) error {
	sb := strings.Builder{}
	sb.WriteString("[Ссылка для твоих друзей](https://t.me/AddPresentBot?start=")
	sb.WriteString(b64url.EncodeToString(calcHash([]byte(c.Chat().Username+"+"), c.Chat().ID, c.Chat().Username)))
	sb.WriteString(")")
	return c.Send(sb.String(), &tg.SendOptions{ParseMode: tg.ModeMarkdownV2})
}
