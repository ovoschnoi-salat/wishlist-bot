package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"wishlist_bot/repository"
)

var (
	useUsernameBtn = MyLocalizedButton{unique: "use_username", localKey: "use_username_btn_text"}
)

func registerNameHandlers(b *tg.Bot) {
	b.Handle(&useUsernameBtn, useUsernameHandler)

	textStateHandlers[ReadNameState] = readNameHandler
}

func changeNameHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = ReadNameState
	return sendNameInputMsg(c, &ctx)
}

func sendNameInputMsg(c tg.Context, ctx *UserCtx) error {
	keyboard := &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{
		{useUsernameBtn.GetInlineButton(ctx.Language)},
		{backToSettingsBtn.GetInlineButton(ctx.Language)},
	}}
	msg := fmt.Sprintf("%s\n%s%s",
		localizer.Get(ctx.Language, "name_input_msg"),
		localizer.Get(ctx.Language, "name_input_current_name_msg"),
		ctx.Name)
	return myEditOrSend(c, ctx, msg, keyboard)
}

func useUsernameHandler(c tg.Context) error {
	newName := c.Chat().Username
	if newName == "" {
		sendAlert(c, fmt.Sprintf("error getting username"))
		return nil
	}
	return saveName(c, "@"+newName)
}

func readNameHandler(c tg.Context) error {
	newName := c.Text()
	// TODO: validate name
	return saveName(c, newName)
}

func saveName(c tg.Context, name string) error {
	log.Printf("saving name: %v", name)
	ctx := GetUserState(c.Chat().ID)
	err := db.Save(&repository.User{ID: ctx.UserId, Language: ctx.Language, Username: c.Chat().Username, Name: name}).Error
	if err != nil {
		sendAlert(c, fmt.Sprintf("error saving new name: %v", err))
		return nil
	}
	ctx.Name = name
	ctx.State = DefaultState
	return sendSettings(c, &ctx)
}
