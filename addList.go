package main

import (
	tg "gopkg.in/telebot.v3"
	"wishlist_bot/repository"
)

func registerNewListHandlers(b *tg.Bot) {
	textStateHandlers[ReadNewListTitleState] = readNewListTitleHandler
}

func sendNewListMessage(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	count, err := repository.CountUserLists(db, ctx.UserId)
	if err != nil {
		sendAlert(c, "error counting lists: "+err.Error())
	}
	if count > 5 {
		sendAlert(c, localizer.Get(ctx.Language, "too_many_lists"))
		return nil
	}
	ctx.State = ReadNewListTitleState
	keyboard := &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{{backToMyListsBtn.GetInlineButton(ctx.Language)}}}
	return myEditOrSend(c, &ctx, localizer.Get(ctx.Language, "new_list_msg"), keyboard)
}

func readNewListTitleHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = DefaultState
	list := repository.List{Title: c.Text(), OwnerID: ctx.UserId}
	err := repository.AddList(db, &list)
	if err != nil {
		return sendError(c, &ctx, "error creating list: "+err.Error())
	}
	ctx.ListId = list.ID
	return sendMyList(c, &ctx)
}
