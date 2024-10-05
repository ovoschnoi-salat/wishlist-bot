package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"strconv"
	"strings"

	"wishlist_bot/repository"
)

var (
	bookBtn          = MyLocalizedButton{unique: "book", localKey: "book_wish_btn_text"}
	cancelBookingBtn = MyLocalizedButton{unique: "cancelBooking", localKey: "cancel_booking_btn_text"}
)

func registerFriendWishHandlers(b *tg.Bot) {
	b.Handle(&bookBtn, bookFriendWishHandler)
	b.Handle(&cancelBookingBtn, cancelFriendWishBookingHandler)
}

func showFriendWishHandler(c tg.Context) error {
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
	return sendFriendWish(c, &ctx)
}

func sendFriendWish(c tg.Context, ctx *UserCtx) error {
	wish, err := repository.GetWish(db, ctx.WishId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting wish: %v", err))
		return nil
	}
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "my_wish_title_msg"))
	sb.WriteString(wish.Title)
	sb.WriteByte('\n')
	if wish.Url != "" {
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_url_msg"))
		sb.WriteString(wish.Url)
		sb.WriteByte('\n')
	}
	if wish.Description != "" {
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_description_msg"))
		sb.WriteString(wish.Description)
		sb.WriteByte('\n')
	}
	if wish.Price != "" {
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_price_msg"))
		sb.WriteString(wish.Price)
		sb.WriteByte('\n')
	}
	if wish.ReservationFree {
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_reservation_free_msg"))
	} else if wish.ReservedBy == 0 {
		sb.WriteString(localizer.Get(ctx.Language, "my_wish_reservable_msg"))
	} else {
		sb.WriteString(localizer.Get(ctx.Language, "friend_wish_reserved_msg"))
	}
	keyboard := getFriendWishKeyboard(ctx, wish.ReservationFree, wish.ReservedBy)
	return myEditOrSend(c, ctx, sb.String(), keyboard, tg.NoPreview)

}

func getFriendWishKeyboard(ctx *UserCtx, reservationFree bool, reservedBy int64) *tg.ReplyMarkup {
	keyboard := make([][]tg.InlineButton, 0, 2)
	if !reservationFree {
		if reservedBy == ctx.UserId {
			keyboard = append(keyboard, []tg.InlineButton{cancelBookingBtn.GetInlineButton(ctx.Language)})
		} else if reservedBy == 0 {
			keyboard = append(keyboard, []tg.InlineButton{bookBtn.GetInlineButton(ctx.Language)})
		}
	}
	keyboard = append(keyboard, []tg.InlineButton{backToFriendListBtn.GetInlineButton(ctx.Language)})
	return &tg.ReplyMarkup{InlineKeyboard: keyboard}
}

func bookFriendWishHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	err := repository.ReserveWish(db, ctx.WishId, ctx.UserId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("Error reserving wish: %v", err))
		return nil
	}
	return sendFriendWish(c, &ctx)
}

func cancelFriendWishBookingHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	err := repository.UndoReservation(db, ctx.WishId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("Error reserving wish: %v", err))
		return nil
	}
	return sendFriendWish(c, &ctx)
}
