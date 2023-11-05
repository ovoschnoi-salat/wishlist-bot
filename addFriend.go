package main

import tg "gopkg.in/telebot.v3"

func registerAddFriendHandlers(bot *tg.Bot) {
	bot.Handle(&showFriendsListBtn, showFriendsList)
}

func addFriend(c tg.Context) error {
	return c.Send("...")
}
