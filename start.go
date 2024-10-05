package main

import (
	tg "gopkg.in/telebot.v3"
	"log"
	"time"
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
		ctx.State = StartState
		return sendLanguageMenu(c, &ctx)
	}
	ctx.State = DefaultState
	return sendMainMenu(c, &ctx)
}
