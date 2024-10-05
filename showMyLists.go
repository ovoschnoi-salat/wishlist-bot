package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"strings"
	"wishlist_bot/repository"
)

var (
	myListSelectBtn  = MySelectorBtn{unique: "my_list_select"}
	addListBtn       = MyLocalizedButton{unique: "add_list", localKey: "add_list_btn_text"}
	backToMyListsBtn = MyLocalizedButton{unique: "back_to_my_lists", localKey: "back_btn_text"}
)

func registerMyListsHandlers(b *tg.Bot) {
	b.Handle(&myListSelectBtn, showMyListHandler)
	b.Handle(&addListBtn, sendNewListMessage)
	b.Handle(&backToMyListsBtn, showMyListsHandler)
}

func showMyListsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendId = 0
	return sendMyLists(c, &ctx)
}

func sendMyLists(c tg.Context, ctx *UserCtx) error {
	lists, err := repository.GetUserLists(db, ctx.UserId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting lists: %v", err))
		return nil
	}
	keyboard := getMyListsKeyboard(ctx.Language, lists)
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "my_lists_msg"))
	sb.WriteByte('\n')
	buildListsMsg(&sb, lists)
	return myEditOrSend(c, ctx, sb.String(), keyboard)
}

func getMyListsKeyboard(lang string, lists []repository.List) *tg.ReplyMarkup {
	keyboard := append(getListsSelectors(lists, myListSelectBtn),
		[]tg.InlineButton{addListBtn.GetInlineButton(lang)},
		[]tg.InlineButton{backToMainMenuBtn.GetInlineButton(lang)})
	return &tg.ReplyMarkup{
		InlineKeyboard: keyboard,
	}
}
