package main

import (
	"fmt"
	tg "gopkg.in/telebot.v3"
	"log"
	"strconv"
	"strings"
	"wishlist_bot/repository"
)

var (
	friendSelectBtn             = MySelectorBtn{unique: "friend_select"}
	anotherListOfFriendsPageBtn = MyPageNavBtn{unique: "another_list_of_friends_page"}
	addFriendBtn                = MyLocalizedButton{unique: "add_friend", localKey: "add_friend_btn_text"}
	backToListOfFriendsBtn      = MyLocalizedButton{unique: "back_to_list_of_friends", localKey: "back_btn_text"}
)

func registerListOfFriendsHandlers(b *tg.Bot) {
	b.Handle(&friendSelectBtn, showFriendListsHandler)
	b.Handle(&anotherListOfFriendsPageBtn, showAnotherListOfFriendsPage)
	b.Handle(&addFriendBtn, addFriendHandler)
	b.Handle(&backToListOfFriendsBtn, showListOfFriendsHandler)

}

func showListOfFriendsHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendId = 0
	ctx.FriendsPageNumber = 0
	return sendListOfFriendsMessage(c, &ctx)
}

func showAnotherListOfFriendsPage(c tg.Context) error {
	num, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error reading page number: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendsPageNumber = num
	return sendListOfFriendsMessage(c, &ctx)
}

func sendListOfFriendsMessage(c tg.Context, ctx *UserCtx) error {
	totalFriendsCount, err := repository.GetListOfFriendsSize(db, ctx.UserId)
	if err != nil {
		return sendError(c, ctx, fmt.Sprintf("error getting friends list size: "+err.Error()))
	}
	pages := (totalFriendsCount + 5) / 6
	if ctx.FriendsPageNumber != 0 && ctx.FriendsPageNumber >= (totalFriendsCount+5)/6 {
		ctx.FriendsPageNumber = max(pages-1, 0)
		log.Printf("error friends page number too big, ctx: %+v , page: %d\n", ctx, totalFriendsCount)
	}
	friendsList, err := repository.GetListOfFriends(db, c.Chat().ID, ctx.FriendsPageNumber)
	if err != nil {
		return sendError(c, ctx, fmt.Sprintf("error getting friends list: "+err.Error()))
	}
	keyboard := getFriendsKeyboard(friendsList, ctx.Language, ctx.FriendsPageNumber, pages)

	b := strings.Builder{}
	b.WriteString(localizer.Get(ctx.Language, "friends_list_msg") + "\n\n")
	if len(friendsList) == 0 {
		b.WriteString(localizer.Get(ctx.Language, "empty_friends_list_msg"))
	} else {
		for i, friend := range friendsList {
			b.WriteString(emojiNumbers[i])
			b.WriteString(" ")
			writeMDV2UserLinkToBuilder(&b, &friend)
			b.WriteString("\n")
		}
		addPageNumber(&b, ctx.Language, ctx.FriendsPageNumber, pages)
	}
	return myEditOrSend(c, ctx, b.String(), markdownV2, keyboard, tg.NoPreview)
}

func getFriendsKeyboard(list []repository.User, lang string, page, totalPages int64) *tg.ReplyMarkup {
	keyboard := make([][]tg.InlineButton, 0)
	for i := 0; i < (len(list)+2)/3 && i < 2; i++ {
		row := make([]tg.InlineButton, 0)
		for j := i * 3; j < i*3+3 && j < len(list); j++ {
			row = append(row, friendSelectBtn.GetInlineButton(j, list[j].ID))
		}
		keyboard = append(keyboard, row)
	}
	if totalPages > 1 {
		row := make([]tg.InlineButton, 0, 2)
		if page > 0 {
			row = append(row, anotherListOfFriendsPageBtn.GetInlineButton("<<", page-1))
		}
		if page < totalPages-1 {
			row = append(row, anotherListOfFriendsPageBtn.GetInlineButton(">>", page+1))
		}
		keyboard = append(keyboard, row)
	}
	keyboard = append(keyboard, []tg.InlineButton{addFriendBtn.GetInlineButton(lang)})
	keyboard = append(keyboard, []tg.InlineButton{backToMainMenuBtn.GetInlineButton(lang)})
	return &tg.ReplyMarkup{InlineKeyboard: keyboard}
}
