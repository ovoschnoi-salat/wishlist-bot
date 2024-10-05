package main

import (
	tg "gopkg.in/telebot.v3"
	"wishlist_bot/repository"
)

func registerAddWishHandlers(b *tg.Bot) {
	textStateHandlers[ReadNewWishTitleState] = readWishTitleHandler
}

func addWishHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = ReadNewWishTitleState
	msg := localizer.Get(ctx.Language, "add_wish_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToMyListBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readWishTitleHandler(c tg.Context) error {
	newTitle := c.Text()
	if newTitle == "" {
		sendAlert(c, "error: empty title")
		return nil
	}
	if len(newTitle) > 100 {
		sendAlert(c, "error: title too long (max 100 characters)")
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	wishId, err := repository.AddWish(db, ctx.UserId, ctx.ListId, newTitle)
	if err != nil {
		sendAlert(c, "error creating new wish: "+err.Error())
		return nil
	}
	ctx.State = DefaultState
	ctx.WishId = wishId
	return sendMyWish(c, &ctx)
}
