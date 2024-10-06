package main

import (
	tg "gopkg.in/telebot.v3"
	"log"
	"wishlist_bot/repository"
)

const languageMenu = "Choose language:\nВыберите язык:"

var (
	languagesKeyboard [][]tg.InlineButton

	langSelectBtn = MyBasicButton{unique: "language_select"}
)

func registerLangHandlers(b *tg.Bot) {
	b.Handle(&langSelectBtn, langHandler)
}

func changeLangHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendLanguageMenu(c, &ctx)
}

func sendLanguageMenu(c tg.Context, ctx *UserCtx) error {
	keyboard := &tg.ReplyMarkup{
		InlineKeyboard: getChangeLangKeyboard(ctx.Language),
	}
	return myEditOrSend(c, ctx, languageMenu, keyboard)
}

func getChangeLangKeyboard(lang string) [][]tg.InlineButton {
	if lang != "" {
		return append(languagesKeyboard, []tg.InlineButton{backToSettingsBtn.GetInlineButton(lang)})
	}
	return languagesKeyboard
}

func langHandler(c tg.Context) error {
	lang := c.Data()
	ctx := GetUserState(c.Chat().ID)
	if !localizer.CheckLangAvailable(lang) {
		log.Println("wrong language received:", lang)
		sendAlert(c, "Error: language not available")
		return sendLanguageMenu(c, &ctx)
	}
	if ctx.Language == lang {
		sendAlert(c, "same lang")
	}
	ctx.Language = lang

	err := repository.UpdateUserLanguage(db, ctx.UserId, lang)
	if err != nil {
		log.Println(err)
		return sendError(c, &ctx, "Error: "+err.Error())
	}
	if ctx.State == StartState {
		ctx.State = DefaultState
		return sendMainMenu(c, &ctx)
	}
	return sendSettings(c, &ctx)
}
