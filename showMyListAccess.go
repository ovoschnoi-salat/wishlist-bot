package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	tg "gopkg.in/telebot.v3"
	"wishlist_bot/repository"
)

var (
	grantFriendAccessBtn       = MySelectorBtn{unique: "grant_friend_access"}
	anotherMyListAccessPageBtn = MyPageNavBtn{unique: "another_my_list_access_page"}
)

func registerMyListAccessHandlers(b *tg.Bot) {
	b.Handle(&grantFriendAccessBtn, grantFriendAccessHandler)
	b.Handle(&anotherMyListAccessPageBtn, showAnotherListAccessPageHandler)
}

func showMyListAccessHandler(c tg.Context) error {
	ctx := GetUserState(c.Chat().ID)
	return sendMyListAccessMessage(c, &ctx)
}

func showAnotherListAccessPageHandler(c tg.Context) error {
	num, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error reading page number: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	ctx.FriendsPageNumber = num
	return sendMyListAccessMessage(c, &ctx)
}

func sendMyListAccessMessage(c tg.Context, ctx *UserCtx) error {
	totalFriendsCount, err := repository.GetListOfFriendsSize(db, ctx.UserId)
	if err != nil {
		return sendError(c, ctx, fmt.Sprintf("error getting friends list size: "+err.Error()))
	}
	pages := (totalFriendsCount + 5) / 6
	if ctx.FriendsPageNumber != 0 && ctx.FriendsPageNumber >= (totalFriendsCount+5)/6 {
		ctx.FriendsPageNumber = max(pages-1, 0)
		log.Printf("error friends page number too big, ctx: %+v , page: %d\n", ctx, totalFriendsCount)
	}
	friendsList, err := repository.GetFriendsAccessListForList(db, c.Chat().ID, ctx.ListId, ctx.FriendsPageNumber)
	if err != nil {
		return sendError(c, ctx, fmt.Sprintf("error getting friends list: "+err.Error()))
	}
	keyboard := getMyListAccessKeyboard(friendsList, ctx.Language, ctx.FriendsPageNumber, pages)

	b := strings.Builder{}
	b.WriteString(localizer.Get(ctx.Language, "friends_list_msg") + "\n\n")
	if len(friendsList) == 0 {
		b.WriteString(localizer.Get(ctx.Language, "empty_friends_list_msg"))
	} else {
		for i, friend := range friendsList {
			b.WriteString(emojiNumbers[i])
			b.WriteString(" ")
			if friend.HasAccess {

				b.WriteRune('✅')
			} else {

				b.WriteRune('🚫')
			}
			b.WriteString(" ")
			writeMDV2UserLinkToBuilder(&b, &friend.User)
			b.WriteString("\n")
		}
		addPageNumber(&b, ctx.Language, ctx.FriendsPageNumber, pages)
	}
	return myEditOrSend(c, ctx, b.String(), markdownV2, keyboard, tg.NoPreview)

}

func getMyListAccessKeyboard(list []repository.UserAccess, lang string, page, totalPages int64) *tg.ReplyMarkup {
	keyboard := make([][]tg.InlineButton, 0)
	for i := 0; i < (len(list)+2)/3 && i < 2; i++ {
		row := make([]tg.InlineButton, 0)
		for j := i * 3; j < i*3+3 && j < len(list); j++ {
			row = append(row, grantFriendAccessBtn.GetInlineButton(j, list[j].ID))
		}
		keyboard = append(keyboard, row)
	}
	if totalPages > 1 {
		row := make([]tg.InlineButton, 0, 2)
		if page > 0 {
			row = append(row, anotherMyListAccessPageBtn.GetInlineButton("<<", page-1))
		}
		if page < totalPages-1 {
			row = append(row, anotherMyListAccessPageBtn.GetInlineButton(">>", page+1))
		}
		keyboard = append(keyboard, row)
	}
	keyboard = append(keyboard, []tg.InlineButton{backToListSettingsBtn.GetInlineButton(lang)})
	return &tg.ReplyMarkup{InlineKeyboard: keyboard}
}

func grantFriendAccessHandler(c tg.Context) error {
	if c.Data() == "" {
		sendAlert(c, "error reading friend is: empty data")
		return nil
	}
	num, err := strconv.ParseInt(c.Data(), 10, 64)
	if err != nil {
		sendAlert(c, "error parsing friend id: "+err.Error())
		return nil
	}
	ctx := GetUserState(c.Chat().ID)
	err = repository.GrantAccess(db, ctx.ListId, num)
	if err != nil {
		sendAlert(c, "error granting friend access: "+err.Error())
		return nil
	}
	return sendMyListAccessMessage(c, &ctx)
}
