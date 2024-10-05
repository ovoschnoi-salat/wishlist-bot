package main

import (
	tg "gopkg.in/telebot.v3"
)

var (
	showMyListsBtn       = MyLocalizedButton{unique: "my_lists", localKey: "show_lists_btn_text"}
	showListOfFriendsBtn = MyLocalizedButton{unique: "list_of_friends", localKey: "show_friend_btn_text"}
	shareBtn             = MyLocalizedButton{unique: "share", localKey: "share_btn_text"}
	showSettingsBtn      = MyLocalizedButton{unique: "settings", localKey: "settings_btn_text"}
	backToMainMenuBtn    = MyLocalizedButton{unique: "back_from_friends", localKey: "back_btn_text"}
)

func registerMainMenuHandlers(bot *tg.Bot) {
	bot.Handle(&showMyListsBtn, showMyListsHandler)
	bot.Handle(&showListOfFriendsBtn, showListOfFriendsHandler)
	bot.Handle(&shareBtn, shareLink)
	bot.Handle(&showSettingsBtn, showSettingsHandler)
	bot.Handle(&backToMainMenuBtn, showMainMenuHandler)
}

func showMainMenuHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendsPageNumber = 0
	ctx.ListPageNumber = 0
	return sendMainMenu(c, &ctx)
}

func sendMainMenu(c tg.Context, ctx *UserCtx) error {
	return myEditOrSend(c, ctx, localizer.Get(ctx.Language, "main_menu"), getMainMenuKeyboard(ctx.Language))
}

func getMainMenuKeyboard(lang string) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{ResizeKeyboard: true,
		RemoveKeyboard: true,
		InlineKeyboard: [][]tg.InlineButton{
			{showMyListsBtn.GetInlineButton(lang)},
			{showListOfFriendsBtn.GetInlineButton(lang)},
			{shareBtn.GetInlineButton(lang)},
			{showSettingsBtn.GetInlineButton(lang)},
		},
	}
}
