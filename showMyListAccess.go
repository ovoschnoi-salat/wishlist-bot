package main

import tg "gopkg.in/telebot.v3"

func showMyListAccessHandler(c tg.Context) error {
	sendAlert(c, "not implemented")
	return nil
	//keyboard := getMyListAccessKeyboard()

}

func getMyListAccessKeyboard(ctx *UserCtx) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListSettingsBtn.GetInlineButton(ctx.Language)},
		},
	}
}
