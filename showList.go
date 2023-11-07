package main

import (
	tg "gopkg.in/telebot.v3"
	"strings"
	"wishlist_bot/storage"
)

const wishBtnUnique = "w"
const anotherWishlistPageBtnUnique = "wp"
const backToFriendsBtnUnique = "bf"

var backBtn = tg.InlineButton{Unique: backToFriendsBtnUnique, Text: "Назад"}

func registerShowListHandlers(bot *tg.Bot) {
	bot.Handle(&showWishlistBtn, showList)
	bot.Handle(getEndpointFromUnique(wishBtnUnique), showWish)
	bot.Handle(getEndpointFromUnique(anotherWishlistPageBtnUnique), showAnotherWishlistPage)
	bot.Handle(&backBtn, showFriendsList)
}

func showList(c tg.Context) error {
	if c.Data() == "" {
		return sendWishlistMessage(c, c.Chat().ID, 0)
	}
	userId, err := getId(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishlistMessage(c, userId, 0)
}

func showAnotherWishlistPage(c tg.Context) error {
	userId, newPage, err := getIdAndData(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendWishlistMessage(c, userId, newPage)
}

func sendWishlistMessage(c tg.Context, userId int64, page int64) error {
	conn, err := pgStorage.Acquire()
	if err != nil {
		return sendError(c, err)
	}
	defer conn.Release()
	wishlistSize, err := conn.GetWishlistSize(userId)
	if err != nil {
		return sendError(c, err)
	}
	pages := (wishlistSize + 5) / 6
	if pages == 0 {
		return c.EditOrSend("пусто", &tg.ReplyMarkup{})
	}
	if page >= pages {
		page = pages - 1
	}
	wishlist, err := conn.GetWishlist(userId, page)
	if err != nil {
		return sendError(c, err)
	}
	keyboard, err := getWishlistKeyboard(wishlist, userId, page, pages, userId != c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}

	b := strings.Builder{}
	if len(wishlist) == 0 {
		b.WriteString("пусто")
	} else {
		for i, wish := range wishlist {
			b.WriteString(emojiNumbers[i])
			b.WriteRune(' ')
			b.WriteString(wish.Title)
			if wish.Url.Valid {
				b.WriteRune('\n')
				b.WriteString(wish.Url.String)
			}
			b.WriteRune('\n')
		}
		addPageNumber(&b, page, pages)
	}
	if userId != c.Chat().ID {
		// todo add back button
	}
	return c.EditOrSend(b.String(), keyboard)
}

func getNextButton(userId, page int64) tg.InlineButton {
	return getNewBtnWithIdAndData(">>", anotherWishlistPageBtnUnique, userId, page+1)
}

func getPrevButton(userId, page int64) tg.InlineButton {
	return getNewBtnWithIdAndData("<<", anotherWishlistPageBtnUnique, userId, page-1)
}

func getWishButton(buttonId int, wishId, page int64) tg.InlineButton {
	return getNewBtnWithIdAndData(emojiNumbers[buttonId], wishBtnUnique, wishId, page)
}

func getWishlistKeyboard(list []storage.Wish, userId, page, pages int64, addBackBtn bool) (*tg.ReplyMarkup, error) {
	keyboard := make([][]tg.InlineButton, 0)
	for i := 0; i < (len(list)+2)/3 && i < 2; i++ {
		row := make([]tg.InlineButton, 0)
		for j := i * 3; j < i*3+3 && j < len(list); j++ {
			row = append(row, getWishButton(j, list[j].Id, page))
		}
		keyboard = append(keyboard, row)
	}
	if pages != 1 {
		if page == 0 {
			keyboard = append(keyboard, []tg.InlineButton{getNextButton(userId, page)})
		} else if page != pages-1 {
			keyboard = append(keyboard, []tg.InlineButton{getPrevButton(userId, page),
				getNextButton(userId, page)})
		} else {
			keyboard = append(keyboard, []tg.InlineButton{getPrevButton(userId, page)})
		}
	}
	if addBackBtn {
		keyboard = append(keyboard, []tg.InlineButton{backBtn})
	}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard}, nil
}
