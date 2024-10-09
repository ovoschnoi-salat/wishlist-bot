package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"time"
	"wishlist_bot/repository"
)

func startHandler(c tg.Context) error {
	defer func() {
		time.Sleep(1 * time.Second)
		err := c.Delete()
		if err != nil {
			log.Println("error deleting start message:", err)
		}
	}()
	ctx := GetUserState(c.Chat().ID)

	if ctx.Language == "" {
		err := repository.AddUser(db, ctx.UserId, c.Chat().Username, "")
		if err != nil {
			log.Println("error adding user to database: " + err.Error())
			return sendError(c, &ctx, "Error: "+err.Error())
		}
		checkForInvitation(c, &ctx)
		ctx.State = StartState
		return sendLanguageMenu(c, &ctx)
	}
	checkForInvitation(c, &ctx)
	if c.Data() == "" {
		ctx.ListPageNumber = 0
		ctx.ListId = 0
		ctx.WishId = 0
		ctx.FriendId = 0
		ctx.FriendsPageNumber = 0
	}
	ctx.State = DefaultState
	return sendMainMenu(c, &ctx)
}

func checkForInvitation(c tg.Context, ctx *UserCtx) {
	if c.Data() != "" {
		if err := addFriendsFromStartPayload(c); err != nil {
			if err := sendError(c, ctx, "Error adding new friend: "+err.Error()); err != nil {
				log.Println("error sending error message:", err)
			}
		}
	}
}

func addFriendsFromStartPayload(c tg.Context) error {
	from, to, err := parseInviteLink(c.Data())
	if err != nil {
		return err
	}
	fromUser, err := repository.GetUserByUsername(db, from)
	if err != nil {
		return err
	}
	toUser, err := repository.GetUserByUsername(db, to)
	if err != nil {
		return err
	}
	if toUser.ID != c.Chat().ID {
		return fmt.Errorf("this link ment for another user")
	}
	return repository.AddFriend(db, fromUser.ID, toUser.ID)
}
