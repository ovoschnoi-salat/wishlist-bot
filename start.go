package main

import (
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
	if c.Data() != "" {
		log.Println(c.Data())
	}
	if ctx.Language == "" {
		err := repository.AddUser(db, ctx.UserId, c.Chat().Username, "")
		if err != nil {
			log.Println("error adding user to database: " + err.Error())
			return sendError(c, &ctx, "Error: "+err.Error())
		}
		ctx.State = StartState
		return sendLanguageMenu(c, &ctx)
	}
	ctx.State = DefaultState
	return sendMainMenu(c, &ctx)
}
