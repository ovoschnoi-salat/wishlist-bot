package main

import tg "gopkg.in/telebot.v3"

var (
	friendsRequestSelector = &tg.ReplyMarkup{}
	btnAccept              = friendsRequestSelector.Data("да", "accept")
	btnDeny                = friendsRequestSelector.Data("нет", "deny")
)

func registerFriendsRequestHandlers(bot *tg.Bot) {
	friendsRequestSelector.Inline(
		friendsRequestSelector.Row(btnAccept, btnDeny),
	)
	bot.Handle(&btnAccept, acceptFriend)
	bot.Handle(&btnDeny, denyFriend)
}

func sendFriendRequest(b *tg.Bot, username string) error {
	chat, err := b.ChatByUsername(username)
	if err != nil {

	}
	_, err = b.Send(chat, "new friends request from "+username, friendsRequestSelector)
	return err
}

func acceptFriend(c tg.Context) error {
	return c.Send("...", friendsRequestSelector)
}

func denyFriend(c tg.Context) error {
	return c.Send("...", friendsRequestSelector)
}
