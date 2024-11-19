package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tg "gopkg.in/telebot.v3"
	"gorm.io/gorm"
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
	username = strings.TrimPrefix(username, "@")
	username = strings.TrimPrefix(username, "https://t.me/")
	username = strings.TrimPrefix(username, "http://t.me/")
	friend, err := repository.GetUserByUsername(db, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.State = DefaultState
			return sendSuggestInvitationLink(c, &ctx, username)
		}
		return sendError(c, &ctx, fmt.Sprintf("error getting user %s: %v", username, err))
	}
	chat, err := c.Bot().ChatByID(friend.ID)
	if err != nil {
		return sendError(c, &ctx, fmt.Sprintf("error getting user %s: %v", username, err))
	}
	// TODO localize
	user, err := repository.GetUserById(db, c.Chat().ID)
	requestMsg := strings.Builder{}
	requestMsg.WriteString(localizer.Get(ctx.Language, "new_friends_request_msg"))
	writeMDV2UserLinkToBuilder(&requestMsg, &user)
	_, err = c.Bot().Send(chat, requestMsg.String(), getFriendRequestKeyboard(ctx.Language, ctx.UserId), tg.ModeMarkdownV2)
	if err != nil {
		return err
	}
	ctx.State = DefaultState
	requestSentMsg := strings.Builder{}
	requestSentMsg.WriteString(localizer.Get(ctx.Language, "new_friends_request_sent_msg"))
	writeMDV2UserLinkToBuilder(&requestSentMsg, &friend)
	keyboard := tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListOfFriendsBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, &ctx, requestSentMsg.String(), &keyboard, tg.ModeMarkdownV2)
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
	msg := strings.Builder{}
	writeMDV2LinkToBuilder(&msg, username, "https://t.me/"+username)
	msg.WriteString(EscapeMarkdown(localizer.Get(ctx.Language, "suggest_invitation_link")))
	link := buildInviteLink(c.Chat().Username, username)
	writeMDV2LinkToBuilder(&msg, localizer.Get(ctx.Language, "invitation_link"), link)
	keyboard := &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{backToListOfFriendsBtn.GetInlineButton(ctx.Language)},
		},
	}
	return myEditOrSend(c, ctx, msg.String(), keyboard, tg.ModeMarkdownV2)
}

func buildInviteLink(from, to string) string {
	linkBuilder := strings.Builder{}
	linkBuilder.WriteString("https://t.me/GiftSyncBot?start=")
	usernames := []byte(from + ":" + to)
	usernamesEncoded := make([]byte, b64url.EncodedLen(len(usernames)))
	b64url.Encode(usernamesEncoded, usernames)
	linkBuilder.Write(usernamesEncoded)
	return linkBuilder.String()
}

func parseInviteLink(link string) (from string, to string, err error) {
	//link = strings.TrimPrefix(link, "https://t.me/GiftSyncBot?start=")
	usernamesDecoded := make([]byte, b64url.DecodedLen(len(link)))
	_, err = b64url.Decode(usernamesDecoded, []byte(link))
	if err != nil {
		return
	}
	usernames := strings.Split(string(usernamesDecoded), ":")
	if len(usernames) != 2 {
		return "", "", fmt.Errorf("invalid invite link")
	}
	return usernames[0], usernames[1], nil
}
