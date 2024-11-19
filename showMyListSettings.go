package main

import (
	"fmt"
	"strings"

	tg "gopkg.in/telebot.v3"
	"wishlist_bot/repository"
)

var (
	changeListTitleBtn      = MyLocalizedButton{unique: "change_list_title", localKey: "change_list_title_btn_text"}
	makeListPrivateBtn      = MyLocalizedButton{unique: "make_list_private", localKey: "make_list_private_btn_text"}
	makeListPublicBtn       = MyLocalizedButton{unique: "make_list_public", localKey: "make_list_public_btn_text"}
	showMyListAccessListBtn = MyLocalizedButton{unique: "show_my_list_access_list", localKey: "show_my_list_access_list_btn_text"}
	removeListBtn           = MyLocalizedButton{unique: "remove_list", localKey: "remove_list_btn_text"}
	approveRemovalBtn       = MyLocalizedButton{unique: "remove_list_removal", localKey: "approve_list_removal_btn_text"}
	backToListSettingsBtn   = MyLocalizedButton{unique: "back_to_list_settings", localKey: "back_btn_text"}
)

func registerListSettingsHandlers(b *tg.Bot) {
	b.Handle(&changeListTitleBtn, showChangeListTitleHandler)
	b.Handle(&makeListPrivateBtn, makeListPrivateHandler)
	b.Handle(&makeListPublicBtn, makeListOpenHandler)
	b.Handle(&showMyListAccessListBtn, showMyListAccessHandler)
	b.Handle(&backToListSettingsBtn, showListSettingsHandler)
	b.Handle(&removeListBtn, askToRemoveListHandler)
	b.Handle(&approveRemovalBtn, removeListHandler)

	textStateHandlers[ReadListNewTitleState] = readListNewTitleHandler
}

func showListSettingsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendsPageNumber = 0
	return sendListSettings(c, &ctx)
}

func sendListSettings(c tg.Context, ctx *UserCtx) error {
	list, err := repository.GetListById(db, ctx.ListId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("Failed to get list: %v", err))
	}
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "list_settings_msg"))
	sb.WriteByte('\n')
	sb.WriteString(localizer.Get(ctx.Language, "list_name_msg"))
	sb.WriteString(list.Title)
	sb.WriteByte('\n')
	sb.WriteString(localizer.Get(ctx.Language, "list_access_msg"))
	if list.Open {
		sb.WriteString(localizer.Get(ctx.Language, "open_list_access_msg"))
	} else {
		sb.WriteString(localizer.Get(ctx.Language, "private_list_access_msg"))
	}
	keyboard := getListSettingsKeyboard(ctx, list.Open)
	return myEditOrSend(c, ctx, sb.String(), keyboard)
}

func getListSettingsKeyboard(ctx *UserCtx, open bool) *tg.ReplyMarkup {
	if open {
		return &tg.ReplyMarkup{
			InlineKeyboard: [][]tg.InlineButton{
				{changeListTitleBtn.GetInlineButton(ctx.Language)},
				{makeListPrivateBtn.GetInlineButton(ctx.Language)},
				{removeListBtn.GetInlineButton(ctx.Language)},
				{backToMyListBtn.GetInlineButton(ctx.Language)},
			},
		}
	}
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{changeListTitleBtn.GetInlineButton(ctx.Language)},
			{makeListPublicBtn.GetInlineButton(ctx.Language)},
			{showMyListAccessListBtn.GetInlineButton(ctx.Language)},
			{removeListBtn.GetInlineButton(ctx.Language)},
			{backToMyListBtn.GetInlineButton(ctx.Language)},
		},
	}
}

func showChangeListTitleHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = ReadListNewTitleState
	msg := localizer.Get(ctx.Language, "change_list_title_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListSettingsBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readListNewTitleHandler(c tg.Context) error {
	newTitle := c.Message().Text
	if newTitle == "" {
		sendAlert(c, "received empty title")
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	if err := repository.UpdateListField(db, ctx.ListId, "Title", newTitle); err != nil {
		sendAlert(c, fmt.Sprintf("Failed to update list title: %v", err))
		return nil
	}
	return sendListSettings(c, &ctx)
}

func makeListPrivateHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	if err := repository.UpdateListField(db, ctx.ListId, "Open", false); err != nil {
		sendAlert(c, fmt.Sprintf("error updating list access: %v", err))
		return nil
	}
	return sendListSettings(c, &ctx)
}

func makeListOpenHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	if err := repository.ClearListAccess(db, ctx.ListId); err != nil {
		sendAlert(c, fmt.Sprintf("error opening list: %v", err))
	}
	if err := repository.UpdateListField(db, ctx.ListId, "Open", true); err != nil {
		sendAlert(c, fmt.Sprintf("error updating list access: %v", err))
	}
	return sendListSettings(c, &ctx)
}

func askToRemoveListHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	msg := localizer.Get(ctx.Language, "ask_to_remove_list_msg")
	keyboard := &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{approveRemovalBtn.GetInlineButton(ctx.Language)},
			{backToListSettingsBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, &ctx, msg, keyboard)
}

func removeListHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	err := repository.DeleteList(db, ctx.ListId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error removing list: %v", err))
		return nil
	}
	ctx.ListId = 0
	return sendMyLists(c, &ctx)
}
