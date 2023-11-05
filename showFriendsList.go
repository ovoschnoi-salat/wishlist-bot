package main

import (
	tg "gopkg.in/telebot.v3"
	"strings"
	"wishlist_bot/storage"
)

const friendBtnUnique = "f"
const friendsPageBtnUnique = "fp"

func registerShowFriendsListHandlers(bot *tg.Bot) {
	bot.Handle(&showFriendsListBtn, showFriendsList)
	bot.Handle(getEndpointFromUnique(friendBtnUnique), showFriendsWishlist)
	bot.Handle(getEndpointFromUnique(friendsPageBtnUnique), showAnotherFriendsPage)
}

func showFriendsList(c tg.Context) error {
	return sendFriendsListMessage(c, 0)
}

func showAnotherFriendsPage(c tg.Context) error {
	id, err := getId(c)
	if err != nil {
		return sendError(c, err)
	}
	return sendFriendsListMessage(c, id)
}

func sendFriendsListMessage(c tg.Context, page int64) error {
	conn, err := pgStorage.Acquire()
	if err != nil {
		return sendError(c, err)
	}
	defer conn.Release()
	wishlistSize, err := conn.GetFriendsListSize(c.Chat().ID)
	if err != nil {
		return sendError(c, err)
	}
	pages := (wishlistSize + 5) / 6
	if pages == 0 {
		return c.EditOrSend("пусто", &tg.ReplyMarkup{})
	}
	if page >= (wishlistSize+5)/6 {
		page = pages - 1
	}
	friendsList, err := conn.GetFriendsList(c.Chat().ID, page)
	if err != nil {
		return sendError(c, err)
	}
	keyboard, err := getFriendsKeyboard(friendsList, page, pages)
	if err != nil {
		return sendError(c, err)
	}
	b := strings.Builder{}
	if len(friendsList) == 0 {
		b.WriteString("пусто")
	} else {
		for i, friend := range friendsList {
			b.WriteString(emojiNumbers[i])
			b.WriteRune(' ')
			b.WriteString(friend.Username)
			b.WriteRune('\n')
		}
		addPageNumber(&b, page, pages)
	}
	return c.EditOrSend(b.String(), keyboard)
}

func getPrevFriendsButton(page int64) tg.InlineButton {
	return getNewBtnWithId(prev, anotherWishlistPageBtnUnique, page-1)
}

func getNextFriendsButton(page int64) tg.InlineButton {
	return getNewBtnWithId(next, anotherWishlistPageBtnUnique, page+1)
}

func getFriendButton(buttonId int, friendId, page int64) tg.InlineButton {
	return getNewBtnWithIdAndData(emojiNumbers[buttonId], friendBtnUnique, friendId, page)
}

func getFriendsKeyboard(list []storage.User, page, pages int64) (*tg.ReplyMarkup, error) {
	keyboard := make([][]tg.InlineButton, 0)
	for i := 0; i < (len(list)+2)/3 && i < 2; i++ {
		row := make([]tg.InlineButton, 0)
		for j := i * 3; j < i*3+3 && j < len(list); j++ {
			row = append(row, getFriendButton(j, list[j].Id, page))
		}
		keyboard = append(keyboard, row)
	}
	if pages != 1 {
		if page == 0 {
			keyboard = append(keyboard, []tg.InlineButton{getNextFriendsButton(page)})
		} else if page != pages-1 {
			keyboard = append(keyboard, []tg.InlineButton{getPrevFriendsButton(page),
				getNextFriendsButton(page)})
		} else {
			keyboard = append(keyboard, []tg.InlineButton{getPrevFriendsButton(page)})
		}
	}
	return &tg.ReplyMarkup{InlineKeyboard: keyboard}, nil
}
