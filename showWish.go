package main

import (
	tg "gopkg.in/telebot.v3"
)

var (
	wishKeyboard = &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{
		{editWishTitleBtn},
		{editWishUrlBtn},
		{deleteWishBtn}}}
	bookBtn          = tg.InlineButton{Unique: "book", Text: "Забронировать"}
	cancelBookingBtn = tg.InlineButton{Unique: "cancelBooking", Text: "Снять бронь"}
	editWishTitleBtn = tg.InlineButton{Unique: "editWishTitle", Text: "Изменить название"}
	editWishUrlBtn   = tg.InlineButton{Unique: "editWishUrl", Text: "Изменить ссылку"}
	deleteWishBtn    = tg.InlineButton{Unique: "deleteWish", Text: "Удалить желание"}
	backToListBtn    = tg.InlineButton{Unique: "backToList", Text: "Назад"}
	backToWishBtn    = tg.InlineButton{Unique: "backToWish", Text: "Отмена"}
)

func registerShowWishHandlers(bot *tg.Bot) {
	bot.Handle(&bookBtn, bookWish)
	bot.Handle(&cancelBookingBtn, cancelBooking)
	bot.Handle(&editWishTitleBtn, editWishTitle)
	bot.Handle(&editWishUrlBtn, editWishUrl)
	bot.Handle(&deleteWishBtn, deleteWish)
	bot.Handle(&backToListBtn, backToList)
	bot.Handle(&backToWishBtn, backToWish)
}

func showWish(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, wishId, returnPage, true)
}

func sendWishMessage(c tg.Context, wishId, returnPage int64, addBackBtn bool) error {
	wish, err := pgStorage.GetWish(wishId)
	if err != nil {
		return sendError(c, err)
	}
	msg := wish.Title
	if wish.Url.Valid {
		msg += "\n" + wish.Url.String
	}
	if c.Chat().ID == wish.Owner {
		return c.EditOrSend(msg, getEditWishKeyboard(wish.Owner, returnPage, wishId, addBackBtn))
	}
	if wish.ReservedBy.Valid {
		if wish.ReservedBy.Int64 == c.Chat().ID {
			return c.EditOrSend(msg, getCancelBookingKeyboard(wish.Owner, returnPage, wish.ReservedBy.Int64))
		}
		return c.EditOrSend(msg, getBackToListKeyboard(wish.Owner, returnPage))
	}
	return c.EditOrSend(msg, getBookWishKeyboard(wish.Owner, returnPage, wishId))
}

func bookWish(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	err = pgStorage.ReserveWish(wishId, c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, wishId, returnPage, true)
}

func cancelBooking(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	err = pgStorage.UndoReservation(wishId)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, wishId, returnPage, true)
}

func editWishUrl(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	err = userStates.SetUserWholeState(c.Chat().ID, UserState{State: ReadUrlState, ChosenWish: wishId})
	if err != nil {
		return sendError(c, err)
	}
	wish, err := pgStorage.GetWish(wishId)
	return c.Edit("Отправьте новую ссылку для \""+wish.Title+"\":", getBackToWishKeyboard(wishId, returnPage))
}

func editWishTitle(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	err = userStates.SetUserWholeState(c.Chat().ID, UserState{State: ReadNewTitleState, ChosenWish: wishId})
	if err != nil {
		return sendError(c, err)
	}
	wish, err := pgStorage.GetWish(wishId)
	if err != nil {
		return sendError(c, err)
	}
	return c.Edit("Отправьте новый заголовок для \""+wish.Title+"\":", getBackToWishKeyboard(wishId, returnPage))
}

func deleteWish(c tg.Context) error {
	wishId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	conn, err := pgStorage.Acquire()
	if err != nil {
		return sendError(c, err)
	}
	wish, err := conn.GetWish(wishId)
	if err != nil {
		return sendError(c, err)
	}
	err = conn.DeleteWish(wishId)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishlistMessage(c, wish.Owner, returnPage)
}

func getEditWishKeyboard(userId, returnPage, wishId int64, addBackBtn bool) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{
		{getBtnWithIdAndData(editWishTitleBtn, wishId, returnPage)},
		{getBtnWithIdAndData(editWishUrlBtn, wishId, returnPage)},
		{getBtnWithIdAndData(deleteWishBtn, wishId, returnPage)}}
	if addBackBtn {
		keyboard = append(keyboard, []tg.InlineButton{getBtnWithIdAndData(backToListBtn, userId, returnPage)})
	}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func getBookWishKeyboard(userId, returnPage, wishId int64) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{{getBtnWithIdAndData(bookBtn, wishId, returnPage)},
		{getBtnWithIdAndData(backToListBtn, userId, returnPage)}}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func getCancelBookingKeyboard(userId, returnPage, wishId int64) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{{getBtnWithIdAndData(cancelBookingBtn, wishId, returnPage)},
		{getBtnWithIdAndData(backToListBtn, userId, returnPage)}}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func getBackToListKeyboard(userId, returnPage int64) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{{getBtnWithIdAndData(backToListBtn, userId, returnPage)}}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func getBackToWishKeyboard(wishId, returnPage int64) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{{getBtnWithIdAndData(backToWishBtn, wishId, returnPage)}}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func readNewTitle(c tg.Context) error {
	state, err := userStates.GetUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	err = pgStorage.EditWishTitle(state.ChosenWish, c.Text())
	if err != nil {
		return sendError(c, err)
	}
	_ = userStates.DeleteUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, state.ChosenWish, 0, true)
}

func backToList(c tg.Context) error {
	userId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishlistMessage(c, userId, returnPage)
}

func backToWish(c tg.Context) error {
	userId, returnPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, userId, returnPage, true)
}
