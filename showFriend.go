package main

import (
	tg "gopkg.in/telebot.v3"
)

func showFriendsWishlist(c tg.Context) error {
	userId, _, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishlistMessage(c, userId, 0)
}
