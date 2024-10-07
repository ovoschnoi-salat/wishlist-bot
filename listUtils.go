package main

import (
	tg "gopkg.in/telebot.v3"
	"strings"
	"wishlist_bot/repository"
)

func writeWishesToBuilder(sb *strings.Builder, ctx *UserCtx, wishes []repository.Wish, pages int64) {
	if len(wishes) == 0 {
		sb.WriteString(localizer.Get(ctx.Language, "no_wishes_msg"))
	} else {
		for i, wish := range wishes {
			sb.WriteString(emojiNumbers[i])
			sb.WriteByte(' ')
			if wish.Url != "" {
				writeMDV2LinkToBuilder(sb, wish.Title, wish.Url)
			} else {
				sb.WriteString(EscapeMarkdown(wish.Title))
			}
			sb.WriteByte('\n')
		}
		addPageNumber(sb, ctx.Language, ctx.ListPageNumber, pages)
	}
}

func getWishesSelectors(wishes []repository.Wish, wishBtn MySelectorBtn, pageBtn MyPageNavBtn, page, totalPages int64) [][]tg.InlineButton {
	keyboard := make([][]tg.InlineButton, 0, (len(wishes)+2)/3+2)
	if len(wishes) == 4 {
		keyboard = append(keyboard,
			[]tg.InlineButton{wishBtn.GetInlineButton(0, wishes[0].ID), wishBtn.GetInlineButton(1, wishes[1].ID)},
			[]tg.InlineButton{wishBtn.GetInlineButton(2, wishes[2].ID), wishBtn.GetInlineButton(3, wishes[3].ID)},
		)
	} else {
		if len(wishes) > 0 {
			row := make([]tg.InlineButton, len(wishes))
			for i := 0; i < len(wishes) && i < 3; i++ {
				row[i] = wishBtn.GetInlineButton(i, wishes[i].ID)
			}
			keyboard = append(keyboard, row)
		}
		if len(wishes) > 4 {
			row := make([]tg.InlineButton, len(wishes)-3)
			for i := 3; i < len(wishes); i++ {
				row[i-3] = wishBtn.GetInlineButton(i, wishes[i].ID)
			}
			keyboard = append(keyboard, row)
		}
	}
	if totalPages > 1 {
		row := make([]tg.InlineButton, 0, 2)
		if page > 0 {
			row = append(row, pageBtn.GetInlineButton("<<", page-1))
		}
		if page < totalPages-1 {
			row = append(row, anotherMyListPageBtn.GetInlineButton(">>", page+1))
		}
		keyboard = append(keyboard, row)
	}
	return keyboard
}
