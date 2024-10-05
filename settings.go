package main

import tg "gopkg.in/telebot.v3"

var (
	changeLanguageBtn = MyLocalizedButton{unique: "change_language", localKey: "change_language_btn_text"}
	changeNameBtn     = MyLocalizedButton{unique: "change_name", localKey: "change_name_btn_text"}
	backToSettingsBtn = MyLocalizedButton{unique: "back_to_settings", localKey: "cancel_btn_text"}
)

func registerSettingsHandlers(b *tg.Bot) {
	b.Handle(&changeLanguageBtn, changeLangHandler)
	b.Handle(&changeNameBtn, changeNameHandler)
	b.Handle(&backToSettingsBtn, showSettingsHandler)
}

func showSettingsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendSettings(c, &ctx)
}

func sendSettings(c tg.Context, ctx *UserCtx) error {
	keyboard := getSettingsKeyboard(ctx.Language)
	return myEditOrSend(c, ctx, localizer.Get(ctx.Language, "settings_msg"), keyboard)
}

func getSettingsKeyboard(lang string) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{changeLanguageBtn.GetInlineButton(lang)},
			{changeNameBtn.GetInlineButton(lang)},
			{backToMainMenuBtn.GetInlineButton(lang)},
		},
	}
}
