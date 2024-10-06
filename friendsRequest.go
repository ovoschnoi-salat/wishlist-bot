package main

import (
	"errors"
	"fmt"
	tg "gopkg.in/telebot.v3"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	acceptFriendsRequestBtn = MyLocalizedDataButton{unique: "accept_friends_request", localKey: "accept_friends_request_btn_text"}
	rejectFriendsRequestBtn = MyLocalizedButton{unique: "reject_friends_request", localKey: "reject_friends_request_btn_text"}
)

func registerFriendRequestHandlers(b *tg.Bot) {
	b.Handle(&acceptFriendsRequestBtn, acceptFriendsRequest)
	b.Handle(&rejectFriendsRequestBtn, rejectFriendsRequest)

	textStateHandlers[ReadNewFriendUsernameState] = sendFriendRequest

}

func addFriendHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = ReadNewFriendUsernameState
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListOfFriendsBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, &ctx, localizer.Get(ctx.Language, "friend_username_input_msg"), &keyboard)
}

func sendFriendRequest(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	username := c.Text()
	if username == "" {
		return sendError(c, &ctx, localizer.Get(ctx.Language, "empty_input_msg"))
	}
	username = strings.ToLower(username)
	username = strings.TrimPrefix(username, "@")
	username = strings.TrimPrefix(username, "https://t.me/")
	username = strings.TrimPrefix(username, "http://t.me/")
	user, err := repository.GetUserByUsername(db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.State = DefaultState
			return sendSuggestInvitationLink(c, &ctx, username)
		}
		return sendError(c, &ctx, fmt.Sprintf("error getting user %s: %v", username, err))
	}
	chat, err := c.Bot().ChatByID(user.ID)
	if err != nil {
		return sendError(c, &ctx, fmt.Sprintf("error getting user %s: %v", username, err))
	}
	// TODO localize
	user, err = repository.GetUserById(db, c.Chat().ID)
	sb := strings.Builder{}
	sb.WriteString(localizer.Get(ctx.Language, "new_friends_request_msg"))
	writeMDV2UserLinkToBuilder(&sb, &user)
	_, err = c.Bot().Send(chat, sb.String(), getFriendRequestKeyboard(ctx.Language, ctx.UserId), tg.ModeMarkdownV2)
	if err != nil {
		return err
	}
	ctx.State = DefaultState
	return sendListOfFriendsMessage(c, &ctx)
}

func getFriendRequestKeyboard(lang string, userId int64) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{acceptFriendsRequestBtn.GetInlineButton(lang, strconv.FormatInt(userId, 10))},
			{rejectFriendsRequestBtn.GetInlineButton(lang)},
		},
	}
}

func acceptFriendsRequest(c tg.Context) error {
	err := c.Delete()
	if err != nil {
		sendAlert(c, "error deleting message")
	}
	friendIdMsg := c.Data()
	if friendIdMsg == "" {
		sendAlert(c, "no friend id")
		return nil
	}
	friendId, err := strconv.ParseInt(friendIdMsg, 10, 64)
	if err != nil {
		sendAlert(c, "invalid friend id: "+friendIdMsg)
		return nil
	}
	err = repository.AddFriend(db, c.Chat().ID, friendId)
	if err != nil {
		sendAlert(c, "error adding friend: "+err.Error())
	}
	return nil
}

func rejectFriendsRequest(c tg.Context) error {
	return c.Delete()
}

func sendSuggestInvitationLink(c tg.Context, ctx *UserCtx, username string) error {
	return sendNotImplemented(c)
}
