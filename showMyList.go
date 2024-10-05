package main

import (
	tg "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	addWishBtn           = MyLocalizedButton{unique: "add_wish", localKey: "add_wish_btn_text"}
	myWishSelectorBtn    = MySelectorBtn{unique: "my_wish_select"}
	anotherMyListPageBtn = MyPageNavBtn{unique: "another_my_list_page"}
	listSettingsBtn      = MyLocalizedButton{unique: "list_settings", localKey: "list_settings_btn_text"}
	backToMyListBtn      = MyLocalizedButton{unique: "back_to_my_list", localKey: "back_btn_text"}
)

func registerMyListHandlers(b *tg.Bot) {
	b.Handle(&addWishBtn, addWishHandler)
	b.Handle(&myWishSelectorBtn, showMyWishHandler)
	b.Handle(&anotherMyListPageBtn, showAnotherListPageHandler)
	b.Handle(&listSettingsBtn, showListSettingsHandler)
	b.Handle(&backToMyListBtn, backToMyListHandler)
}

func showAnotherListPageHandler(c tg.Context) error {
	if c.Data() == "" {
		sendAlert(c, "error reading page number: empty data")
		return nil
	}
	num, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error parsing page number: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.ListPageNumber = num
	return sendMyList(c, &ctx)
}

func backToMyListHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendMyList(c, &ctx)
}

func showMyListHandler(c tg.Context) error {
	if c.Data() == "" {
		sendAlert(c, "error getting list id: empty data")
		return nil
	}
	listId, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error parsing list id: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.WishId = 0
	ctx.ListId = listId
	return sendMyList(c, &ctx)
}

func sendMyList(c tg.Context, ctx *UserCtx) error {
	list, err := repository.GetListById(db, ctx.ListId)
	if err != nil {
		sendAlert(c, "error getting list: "+err.Error())
		return nil
	}
	listSize, err := repository.GetListSize(db, ctx.ListId)
	if err != nil {
		sendAlert(c, "error getting list size: "+err.Error())
		return nil
	}
	pages := (listSize + 5) / 6
	if ctx.ListPageNumber >= pages {
		ctx.ListPageNumber = max(0, pages-1)
	}
	wishes, err := repository.GetWishes(db, ctx.ListId, ctx.ListPageNumber)
	if err != nil {
		sendAlert(c, "error getting wishes: "+err.Error())
		return nil
	}
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "list_name_msg"))
	sb.WriteString(list.Title)
	sb.WriteString("\n\n")
	writeWishesToBuilder(&sb, ctx, wishes, pages)
	keyboard := getMyListKeyboard(wishes, ctx, ctx.ListPageNumber, pages)
	return myEditOrSend(c, ctx, sb.String(), keyboard, tg.ModeMarkdownV2, tg.NoPreview)
}

func getMyListKeyboard(list []repository.Wish, ctx *UserCtx, page, totalPages int64) *tg.ReplyMarkup {
	keyboard := append(getWishesSelectors(list, myWishSelectorBtn, anotherMyListPageBtn, page, totalPages),
		[]tg.InlineButton{addWishBtn.GetInlineButton(ctx.Language)},
		[]tg.InlineButton{listSettingsBtn.GetInlineButton(ctx.Language)},
		[]tg.InlineButton{backToMyListsBtn.GetInlineButton(ctx.Language)})
	return &tg.ReplyMarkup{
		InlineKeyboard: keyboard,
	}
}
