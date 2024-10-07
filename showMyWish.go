package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	editWishTitleBtn         = MyLocalizedButton{unique: "edit_wish_title", localKey: "edit_wish_title_btn_text"}
	editWishUrlBtn           = MyLocalizedButton{unique: "edit_wish_url", localKey: "edit_wish_url_btn_text"}
	removeWishUrlBtn         = MyLocalizedButton{unique: "remove_wish_url", localKey: "remove_wish_url_btn_text"}
	editWishDescriptionBtn   = MyLocalizedButton{unique: "edit_wish_description", localKey: "edit_wish_description_btn_text"}
	removeWishDescriptionBtn = MyLocalizedButton{unique: "remove_wish_description", localKey: "remove_wish_description_btn_text"}
	editWishPriceBtn         = MyLocalizedButton{unique: "edit_wish_price", localKey: "edit_wish_price_btn_text"}
	removeWishPriceBtn       = MyLocalizedButton{unique: "remove_wish_price", localKey: "remove_wish_price_btn_text"}
	makeReservationFreeBtn   = MyLocalizedButton{unique: "make_reservation_free", localKey: "make_reservation_free_btn_text"}
	makeReservableBtn        = MyLocalizedButton{unique: "make_reservable", localKey: "make_reservable_btn_text"}
	deleteWishBtn            = MyLocalizedButton{unique: "delete_wish", localKey: "delete_wish_btn_text"}
	backToMyWishBtn          = MyLocalizedButton{unique: "back_to_my_wish", localKey: "cancel_btn_text"}
)

func registerMyWishHandlers(b *tg.Bot) {
	b.Handle(&editWishTitleBtn, editWishTitleHandler)
	b.Handle(&editWishUrlBtn, editWishUrlHandler)
	b.Handle(&editWishDescriptionBtn, editWishDescriptionHandler)
	b.Handle(&editWishPriceBtn, editWishPriceHandler)
	b.Handle(&deleteWishBtn, deleteWishHandler)
	b.Handle(&backToMyWishBtn, backToMyWishHandler)
	b.Handle(&makeReservationFreeBtn, makeWishReservationFreeHandler)
	b.Handle(&makeReservableBtn, makeWishReservableHandler)

	textStateHandlers[ReadWishNewTitleState] = readWishNewTitleHandler
	textStateHandlers[ReadWishNewUrlState] = readWishNewUrlHandler
	textStateHandlers[ReadWishNewDescriptionState] = readWishNewDescriptionHandler
	textStateHandlers[ReadWishNewPriceState] = readWishNewPriceHandler
}

func backToMyWishHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendMyWish(c, &ctx)
}

func showMyWishHandler(c tg.Context) error {
	if c.Data() == "" {
		sendAlert(c, "error getting wish id")
		return nil
	}
	wishId, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error parsing wish id")
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.WishId = wishId
	return sendMyWish(c, &ctx)
}

func sendMyWish(c tg.Context, ctx *UserCtx) error {
	wish, err := repository.GetWish(db, ctx.WishId)
	if err != nil {
		return sendError(c, ctx, fmt.Sprintf("error getting wish: %v", err))
	}
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "my_wish_title_msg"))
	sb.WriteString(wish.Title)
	sb.WriteByte('\n')
	if wish.Url != "" {
		sb.WriteByte('\n')
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_url_msg"))
		sb.WriteString(wish.Url)
		sb.WriteByte('\n')
	}
	if wish.Description != "" {
		sb.WriteByte('\n')
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_description_msg"))
		sb.WriteString(wish.Description)
		sb.WriteByte('\n')
	}
	if wish.Price != "" {
		sb.WriteByte('\n')
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_price_msg"))
		sb.WriteString(wish.Price)
		sb.WriteByte('\n')
	}
	if wish.ReservationFree {
		sb.WriteByte('\n')
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_reservation_free_msg"))
	} else {
		sb.WriteByte('\n')
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_reservable_msg"))
	}
	keyboard := getMyWishKeyboard(ctx, wish.ReservationFree)
	return myEditOrSend(c, ctx, sb.String(), keyboard, tg.NoPreview)
}

func getMyWishKeyboard(ctx *UserCtx, reservationFree bool) *tg.ReplyMarkup {
	if reservationFree {
		return &tg.ReplyMarkup{
			InlineKeyboard: [][]tg.InlineButton{
				{editWishTitleBtn.GetInlineButton(ctx.Language)},
				{editWishDescriptionBtn.GetInlineButton(ctx.Language)},
				{editWishUrlBtn.GetInlineButton(ctx.Language)},
				{editWishPriceBtn.GetInlineButton(ctx.Language)},
				{makeReservableBtn.GetInlineButton(ctx.Language)},
				{deleteWishBtn.GetInlineButton(ctx.Language)},
				{backToMyListBtn.GetInlineButton(ctx.Language)},
			},
		}
	}
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{editWishTitleBtn.GetInlineButton(ctx.Language)},
			{editWishDescriptionBtn.GetInlineButton(ctx.Language)},
			{editWishUrlBtn.GetInlineButton(ctx.Language)},
			{editWishPriceBtn.GetInlineButton(ctx.Language)},
			{makeReservationFreeBtn.GetInlineButton(ctx.Language)},
			{deleteWishBtn.GetInlineButton(ctx.Language)},
			{backToMyListBtn.GetInlineButton(ctx.Language)},
		},
	}
}

func editWishTitleHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	msg := localizer.Get(ctx.Language, "edit_wish_title_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToMyWishBtn.GetInlineButton(ctx.Language)},
		},
	}
	ctx.State = ReadWishNewTitleState
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readWishNewTitleHandler(c tg.Context) error {
	newTitle := c.Text()
	if newTitle == "" {
		sendAlert(c, "error getting new wish title: empty title") // todo alert doesn't work
	}
	// TODO validate title
	ctx := GetUserState(c.Chat().ID)
	err := repository.UpdateWishField(db, ctx.WishId, "title", newTitle)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	ctx.State = DefaultState
	return sendMyWish(c, &ctx)
}

func editWishDescriptionHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	msg := localizer.Get(ctx.Language, "edit_wish_description_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{removeWishDescriptionBtn.GetInlineButton(ctx.Language)},
			{backToMyWishBtn.GetInlineButton(ctx.Language)},
		},
	}
	ctx.State = ReadWishNewDescriptionState
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readWishNewDescriptionHandler(c tg.Context) error {
	newDescription := c.Text()
	if newDescription == "" {
		sendAlert(c, "error getting new wish title: empty description")
		return nil
	}
	if len(newDescription) > 200 {
		sendAlert(c, "error: new wish description is too long (maximum is 200 characters)") // todo localize
		return nil
	}
	// TODO validate description
	ctx := GetUserState(c.Chat().ID)
	err := repository.UpdateWishField(db, ctx.WishId, "description", newDescription)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	ctx.State = DefaultState
	return sendMyWish(c, &ctx)
}

func editWishUrlHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	msg := localizer.Get(ctx.Language, "edit_wish_url_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{removeWishUrlBtn.GetInlineButton(ctx.Language)},
			{backToMyWishBtn.GetInlineButton(ctx.Language)},
		},
	}
	ctx.State = ReadWishNewUrlState
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readWishNewUrlHandler(c tg.Context) error {
	newUrl := c.Text()
	if newUrl == "" {
		sendAlert(c, "error getting new wish title: empty url")
	}
	// TODO validate url
	ctx := GetUserState(c.Chat().ID)
	err := repository.UpdateWishField(db, ctx.WishId, "url", newUrl)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	ctx.State = DefaultState
	return sendMyWish(c, &ctx)
}

func editWishPriceHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	msg := localizer.Get(ctx.Language, "edit_wish_price_msg")
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{removeWishPriceBtn.GetInlineButton(ctx.Language)},
			{backToMyWishBtn.GetInlineButton(ctx.Language)},
		},
	}
	ctx.State = ReadWishNewPriceState
	return myEditOrSend(c, &ctx, msg, &keyboard)
}

func readWishNewPriceHandler(c tg.Context) error {
	newPrice := c.Text()
	if newPrice == "" {
		sendAlert(c, "error getting new wish title: empty price")
	}
	// TODO validate price
	ctx := GetUserState(c.Chat().ID)
	err := repository.UpdateWishField(db, ctx.WishId, "price", newPrice)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	ctx.State = DefaultState
	return sendMyWish(c, &ctx)
}

func deleteWishHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	err := repository.DeleteWish(db, ctx.WishId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error deleting wish: %v", err))
		return nil
	}
	ctx.WishId = 0
	return sendMyList(c, &ctx)
}

func makeWishReservationFreeHandler(c tg.Context) error {
	log.Println("makeWishReservationFreeHandler")
	ctx := GetUserState(c.Chat().ID)
	err := repository.MakeWishReservationFree(db, ctx.WishId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	log.Println("makeWishReservationFreeHandler finished")
	return sendMyWish(c, &ctx)
}

func makeWishReservableHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	err := repository.MakeWishReservable(db, ctx.WishId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error updating wish: %v", err))
		return nil
	}
	return sendMyWish(c, &ctx)
}
