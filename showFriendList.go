package main

import (
	tg "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	friendWishSelectorBtn    = MySelectorBtn{unique: "friend_wish_select"}
	anotherFriendListPageBtn = MyPageNavBtn{unique: "another_friend_list_page"}
	backToFriendListBtn      = MyLocalizedButton{unique: "back_to_friend_list", localKey: "back_btn_text"}
)

func registerFriendListHandlers(b *tg.Bot) {
	b.Handle(&friendWishSelectorBtn, showFriendWishHandler)
	b.Handle(&anotherFriendListPageBtn, showAnotherFriendListPageHandler)
	b.Handle(&backToFriendListBtn, backToFriendListHandler)
}

func showAnotherFriendListPageHandler(c tg.Context) error {
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
	return sendFriendList(c, &ctx)
}

func backToFriendListHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.WishId = 0
	return sendFriendList(c, &ctx)
}

func showFriendListHandler(c tg.Context) error {
	if c.Data() == "" {
		sendAlert(c, "error getting friends list id")
		return nil
	}
	friendListId, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error parsing list id: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.ListId = friendListId
	return sendFriendList(c, &ctx)
}

func sendFriendList(c tg.Context, ctx *UserCtx) error {
	friend, err := repository.GetUserById(db, ctx.FriendId)
	if err != nil {
		sendAlert(c, "error getting friend: "+err.Error())
		return nil
	}
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
	sb.WriteByte('\n')
	sb.WriteString(localizer.Get(ctx.Language, "list_owner_msg"))
	writeMDV2UserLinkToBuilder(&sb, &friend)
	sb.WriteString("\n\n")
	writeWishesToBuilder(&sb, ctx, wishes, pages)
	keyboard := getFriendListKeyboard(wishes, ctx, ctx.ListPageNumber, pages)
	return myEditOrSend(c, ctx, sb.String(), &keyboard, tg.ModeMarkdownV2, tg.NoPreview)
}

func getFriendListKeyboard(list []repository.Wish, ctx *UserCtx, page, totalPages int64) tg.ReplyMarkup {
	keyboard := append(getWishesSelectors(list, friendWishSelectorBtn, anotherFriendListPageBtn, page, totalPages),
		[]tg.InlineButton{backToFriendListsBtn.GetInlineButton(ctx.Language)})
	return tg.ReplyMarkup{
		InlineKeyboard: keyboard,
	}
}
