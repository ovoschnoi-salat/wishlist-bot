package main

import (
	tg "gopkg.in/telebot.v3"
)

func shareLink(c tg.Context) error {
	sendAlert(c, "not implemented")
	return nil
	//sb := strings.Builder{}
	//sb.WriteString("[Ссылка для твоих друзей](https://t.me/AddPresentBot?start=")
	//sb.WriteString(b64url.EncodeToString(calcHash([]byte(c.Chat().Username+"+"), c.Chat().ID, c.Chat().Username)))
	//sb.WriteString(")")
	//return c.Send(sb.String(), markdownV2, tg.NoPreview)
}
