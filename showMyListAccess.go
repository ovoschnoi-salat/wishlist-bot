package main

import tg "gopkg.in/telebot.v3"

func showMyListAccessHandler(c tg.Context) error {
	return sendNotImplemented(c)
}

func getMyListAccessKeyboard(ctx *UserCtx) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListSettingsBtn.GetInlineButton(ctx.Language)},
		},
	}
}
