package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	friendListSelectBtn     = MySelectorBtn{unique: "friend_list_select"}
	removeFriendBtn         = MyLocalizedButton{unique: "remove_friend", localKey: "remove_friend_btn_text"}
	approveFriendRemovalBtn = MyLocalizedButton{unique: "approve_friend_removal", localKey: "remove_friend_btn_text"}
	backToFriendListsBtn    = MyLocalizedButton{unique: "back_to_friend_lists", localKey: "back_btn_text"}
)

func registerFriendListsHandlers(b *tg.Bot) {
	b.Handle(&friendListSelectBtn, showFriendListHandler)
	b.Handle(&removeFriendBtn, removeFriendHandler)
	b.Handle(&approveFriendRemovalBtn, approveFriendRemovalHandler)
	b.Handle(&backToFriendListsBtn, backToFriendListsHandler)
}

func backToFriendListsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.ListId = 0
	return sendFriendLists(c, &ctx)
}

func showFriendListsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	id, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting friend id: %v", err))
		return nil
	}
	ctx.FriendId = id
	return sendFriendLists(c, &ctx)
}

func removeFriendHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	keyboard := &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{
		{approveFriendRemovalBtn.GetInlineButton(ctx.Language)},
		{backToFriendListsBtn.GetInlineButton(ctx.Language)},
	}}
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "remove_friend_msg"))
	sb.WriteByte(' ')
	user, err := repository.GetUserById(db, ctx.FriendId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting user: %v", err))
		return nil
	}
	writeMDV2UserLinkToBuilder(&sb, &user)
	sb.WriteByte('?')
	return myEditOrSend(c, &ctx, sb.String(), keyboard, tg.ModeMarkdownV2)
}

func approveFriendRemovalHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	if err := repository.DeleteFriend(db, ctx.UserId, ctx.FriendId); err != nil {
		sendAlert(c, fmt.Sprintf("error removing friend: %v", err))
		return nil
	}
	ctx.FriendId = 0
	return sendListOfFriendsMessage(c, &ctx)
}

func sendFriendLists(c tg.Context, ctx *UserCtx) error {
	lists, err := repository.GetFriendLists(db, ctx.UserId, ctx.FriendId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting lists: %v", err))
		return nil
	}
	friend, err := repository.GetUserById(db, ctx.FriendId)
	if err != nil {
		sendAlert(c, fmt.Sprintf("error getting user: %v", err))
		return nil
	}
	keyboard := getFriendListsKeyboard(ctx.Language, lists)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(localizer.Get(ctx.Language, "friends_lists_msg"),
		createMDV2Link(friend.Name, "t.me/"+friend.Username)))
	sb.WriteByte('\n')
	buildListsMsg(&sb, lists)
	return myEditOrSend(c, ctx, sb.String(), keyboard, tg.ModeMarkdownV2, tg.NoPreview)
}

func getFriendListsKeyboard(lang string, lists []repository.List) *tg.ReplyMarkup {
	keyboard := append(getListsSelectors(lists, friendListSelectBtn),
		[]tg.InlineButton{removeFriendBtn.GetInlineButton(lang)},
		[]tg.InlineButton{backToListOfFriendsBtn.GetInlineButton(lang)})
	return &tg.ReplyMarkup{
		InlineKeyboard: keyboard,
	}
}
