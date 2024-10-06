package main

import (
	tg "gopkg.in/telebot.v3"
	"log"
)

var (
	closeErrorBtn = MyLocalizedButton{unique: "close_error", localKey: "close_btn_text"}
)

func registerMiscHandlers(b *tg.Bot) {
	b.Handle(&closeErrorBtn, closeErrorHandler)
}

func sendNotImplemented(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendError(c, &ctx, localizer.Get(ctx.Language, "not_implemented_msg"))
}

func cancel(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.State = DefaultState
	return sendMainMenu(c, &ctx)
}

func unknownCommandHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendError(c, &ctx, "Простите, бот немного запутался и не смог понять, что вы сейчас сделали.\nПопробуйте заново.")
}

func myEditOrSend(c tg.Context, ctx *UserCtx, what interface{}, opts ...interface{}) error {
	if ctx.ErrorMessageId != 0 {
		if err := c.Bot().Delete((*ErrorCtx)(ctx)); err != nil {
			log.Println("error deleting error message:", err)
		} else {
			ctx.ErrorMessageId = 0
		}
	}
	if c.Message().Sender.ID == 6729731825 && ctx.LastMessageId == 0 {
		ctx.LastMessageId = c.Message().ID
	}
	if ctx.LastMessageId != 0 {
		if _, err := c.Bot().Edit(ctx, what, opts...); err == nil {
			ctxStorage.Set(c.Chat().ID, *ctx)
			//log.Println("save", ctx)
			return nil
		} else {
			log.Println("error editing message:", err)
		}
		if err := c.Bot().Delete(ctx); err != nil {
			log.Println("error deleting message:", err)
		}
	}
	message, err := c.Bot().Send(c.Recipient(), what, opts...)
	if err != nil {
		ctxStorage.Set(c.Chat().ID, *ctx)
		//log.Println("save", ctx)
		return err
	}
	ctx.LastMessageId = message.ID
	log.Println("resent msg")
	ctxStorage.Set(c.Chat().ID, *ctx)
	return nil
}

func sendError(c tg.Context, ctx *UserCtx, msg string) error {
	keyboard := getErrorKeyboard(ctx)
	if ctx.ErrorMessageId != 0 {
		_, err := c.Bot().Edit((*ErrorCtx)(ctx), msg, keyboard)
		if err == nil {
			ctxStorage.Set(c.Chat().ID, *ctx)
			return nil
		}
		log.Println("error editing error message:", err)
		if err := c.Bot().Delete((*ErrorCtx)(ctx)); err != nil {
			log.Println("error deleting error message:", err)
		}
	}
	message, err := c.Bot().Send(c.Recipient(), msg, keyboard)
	if err != nil {
		ctxStorage.Set(c.Chat().ID, *ctx)
		return err
	}
	ctx.ErrorMessageId = message.ID
	log.Println("resent error")
	ctxStorage.Set(c.Chat().ID, *ctx)
	return nil
}

func getErrorKeyboard(ctx *UserCtx) *tg.ReplyMarkup {
	return &tg.ReplyMarkup{
		InlineKeyboard: [][]tg.InlineButton{
			{closeErrorBtn.GetInlineButton(ctx.Language)},
		},
	}

}

func sendAlert(c tg.Context, msg string) {
	r := tg.CallbackResponse{
		Text:      msg,
		ShowAlert: true,
	}
	log.Println(msg)
	err := c.Respond(&r)
	if err != nil {
		log.Println("error sending alert:", err)
	}
}

func closeErrorHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.ErrorMessageId = 0
	ctxStorage.Set(c.Chat().ID, ctx)
	return c.Delete()
}
