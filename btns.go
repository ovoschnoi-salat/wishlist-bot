package main

import (
	"strconv"

	tg "gopkg.in/telebot.v3"
)

func registerAllHandlers(b *tg.Bot) {
	registerMainMenuHandlers(b)
	registerFriendRequestHandlers(b)
	registerMyListsHandlers(b)
	registerMyListHandlers(b)
	registerSettingsHandlers(b)
	registerListOfFriendsHandlers(b)
	registerLangHandlers(b)
	registerNameHandlers(b)
	registerFriendListsHandlers(b)
	registerFriendListHandlers(b)
	registerNewListHandlers(b)
	registerListSettingsHandlers(b)
	registerFriendWishHandlers(b)
	registerMyWishHandlers(b)
	registerAddWishHandlers(b)
	registerMiscHandlers(b)
	registerMyListAccessHandlers(b)
}

type MyBasicButton struct {
	unique string
}

func (b *MyBasicButton) GetInlineButton(text, data string) tg.InlineButton {
	return tg.InlineButton{
		Unique: b.unique,
		Text:   text,
		Data:   data,
	}
}

func (b *MyBasicButton) CallbackUnique() string {
	return "\f" + b.unique
}

type MyLocalizedButton struct {
	unique   string
	localKey string
}

func (b *MyLocalizedButton) GetInlineButton(lang string) tg.InlineButton {
	return tg.InlineButton{
		Unique: b.unique,
		Text:   localizer.Get(lang, b.localKey),
	}
}

func (b *MyLocalizedButton) CallbackUnique() string {
	return "\f" + b.unique
}

type MyLocalizedDataButton struct {
	unique   string
	localKey string
}

func (b *MyLocalizedDataButton) GetInlineButton(lang string, data string) tg.InlineButton {
	return tg.InlineButton{
		Unique: b.unique,
		Text:   localizer.Get(lang, b.localKey),
		Data:   data,
	}
}

func (b *MyLocalizedDataButton) CallbackUnique() string {
	return "\f" + b.unique
}

type MySelectorBtn struct {
	unique string
}

func (b *MySelectorBtn) GetInlineButton(num int, selectedId int64) tg.InlineButton {
	return tg.InlineButton{
		Unique: b.unique,
		Text:   emojiNumbers[num],
		Data:   strconv.FormatInt(selectedId, 10),
	}
}

func (b *MySelectorBtn) CallbackUnique() string {
	return "\f" + b.unique
}

type MyPageNavBtn struct {
	unique string
}

func (b *MyPageNavBtn) GetInlineButton(text string, selectedId int64) tg.InlineButton {
	return tg.InlineButton{
		Unique: b.unique,
		Text:   text,
		Data:   strconv.FormatInt(selectedId, 10),
	}
}

func (b *MyPageNavBtn) CallbackUnique() string {
	return "\f" + b.unique
}
