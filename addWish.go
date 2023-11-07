package main

import (
	"encoding/base64"
	"errors"
	tg "gopkg.in/telebot.v3"
	"strings"
)

var (
	acceptNewWishBtn  = tg.InlineButton{Unique: "acceptWish", Text: "✅"}
	declineNewWishBtn = tg.InlineButton{Unique: "declineWish", Text: "❌"}
	noUrlBtn          = tg.InlineButton{Unique: "noUrl", Text: "Не добавлять ссылку"}
	addUrlKeyboard    = &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{{noUrlBtn}}}
)

func registerAddWishHandlers(bot *tg.Bot) {
	bot.Handle(&addWishBtn, addWish)
	bot.Handle(&noUrlBtn, setNoUrl)
	bot.Handle(&acceptNewWishBtn, acceptUnexpectedWish)
	bot.Handle(&declineNewWishBtn, cancel)
}

func addWish(c tg.Context) error {
	err := userStates.SetUserState(c.Chat().ID, NewWishState)
	if err != nil {
		return sendError(c, err)
	}
	return c.Send("Введите название:", cancelKeyboard)
}

func readWish(c tg.Context) error {
	id, err := pgStorage.AddWish(c.Chat().ID, c.Text(), "")
	if err != nil {
		return sendError(c, err)
	}
	_ = userStates.SetUserWholeState(c.Chat().ID, UserState{ChosenWish: id, State: ReadUrlState})
	if err != nil {
		return sendError(c, err)
	}
	return c.Send("Введите url:", addUrlKeyboard)
}

func readUrl(c tg.Context) error {
	state, err := userStates.GetUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	err = pgStorage.EditWishUrl(state.ChosenWish, c.Text())
	if err != nil {
		return sendError(c, err)
	}
	err = userStates.DeleteUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, state.ChosenWish, 0, false)
}

func setNoUrl(c tg.Context) error {
	state, err := userStates.GetUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	err = userStates.DeleteUserState(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, state.ChosenWish, 0, false)
}

func readUnexpectedWish(c tg.Context) error {
	title, url := parseUnexpectedWish(c.Data())
	buf := make([]byte, b64.EncodedLen(len(title)+len(url))+1)
	b64.Encode(buf, []byte(title))
	buf[b64.EncodedLen(len(title))] = ':'
	b64.Encode(buf[b64.EncodedLen(len(title))+1:], []byte(url))
	msg := strings.Builder{}
	msg.WriteString("Вы хотели добавить делание?\nНазвание желания: ")
	msg.WriteString(title)
	if url != "" {
		msg.WriteString("\nСсылка желания: ")
		msg.WriteString(url)
	}
	msg.WriteString("\nДобавить новое желание?")
	return c.Send(msg.String(), getUnexpectedWishKeyboard(string(buf)))
}

func getUnexpectedWishKeyboard(data string) *tg.ReplyMarkup {
	keyboard := [][]tg.InlineButton{{getBtnWithData(acceptNewWishBtn, data), declineNewWishBtn}}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard, RemoveKeyboard: true}
}

func parseUnexpectedWish(msg string) (string, string) {
	array := strings.Split(msg, "\n")
	if len(array) > 1 && rxURL.MatchString(array[len(array)-1]) {
		url := array[len(array)-1]
		return msg[:len(msg)-len(url)-1], url
	}
	return msg, ""
}

func acceptUnexpectedWish(c tg.Context) error {
	array := strings.Split(c.Data(), ":")
	if len(array) != 2 {
		return sendError(c, errors.New("wrong data passed:"+c.Data()))
	}
	b64 := base64.StdEncoding
	title, err := b64.DecodeString(array[0])
	if err != nil {
		return sendError(c, err)
	}
	url, err := b64.DecodeString(array[0])
	if err != nil {
		return sendError(c, err)
	}
	conn, err := pgStorage.Acquire()
	if err != nil {
		return sendError(c, err)
	}
	wishId, err := conn.AddWish(c.Chat().ID, string(title), string(url))
	if err != nil {
		return sendError(c, err)
	}
	return sendWishMessage(c, wishId, 0, false)
}
