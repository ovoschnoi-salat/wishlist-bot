package main

import (
	"crypto/cipher"
	"encoding/base64"
	tg "gopkg.in/telebot.v3"
	"regexp"
	"strconv"
	"strings"
	f "wishlist_bot/repository"
)

var (
	markdownV2   = &tg.SendOptions{ParseMode: tg.ModeMarkdownV2}
	emojiNumbers = []string{
		"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣",
	}
	aesCipher cipher.Block
	b64       = base64.StdEncoding
	b64url    = base64.RawURLEncoding
	rxURL     = regexp.MustCompile(URL)
	URL       = "^(https?:\\/\\/)?" +
		"(([a-zA-Z0-9]([a-zA-Z0-9-]{0,253}[a-zA-Z0-9])?\\.)+[a-zA-Z0-9]([a-zA-Z0-9-]{0,253}[a-zA-Z0-9])?|" +
		"([а-яА-Я0-9]([а-яА-Я0-9-]{0,253}[а-яА-Я0-9])?\\.)+(рф|РФ)|" +
		"((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d)\\.?\\b){4}" +
		"(:(6(553[0-5]|55[0-2][0-9]|5[0-4][0-9]{2}|[0-4][0-9]{3})|[1-5][0-9]{4}|[1-9][0-9]{0,3}))?)" +
		"(\\/([a-zA-Z0-9-._~:@!$&'()*+,;=]|%[a-fA-f0-9]{2})*)*" +
		"(\\?([a-zA-Z0-9-._~:@!$&'()*+,;=/\\?]|%[a-fA-f0-9]{2})*)?" +
		"(#([a-zA-Z0-9-._~:@!$&'()*+,;=/\\?]|%[a-fA-f0-9]{2})*)?$"
)

func addPageNumber(str *strings.Builder, lang string, index, all int64) {
	str.WriteString("\n " + localizer.Get(lang, "page_number") + " ")
	str.WriteString(strconv.FormatInt(index+1, 10))
	str.WriteRune('/')
	str.WriteString(strconv.FormatInt(all, 10))
}

//func calcHash(array []byte, userId int64, username string) []byte {
//	h := sha3.New224()
//	id := make([]byte, 8)
//	binary.LittleEndian.PutUint64(id, uint64(userId))
//	h.Write(id)
//	h.Write([]byte(username))
//	binary.BigEndian.PutUint64(id, uint64(userId))
//	h.Write(id)
//	return h.Sum(array)
//}

func createMDV2Link(title, url string) string {
	return "[" + EscapeMarkdown(title) + "](" + EscapeMarkdownLink(url) + ")"
}

func buildListsMsg(sb *strings.Builder, lists []f.List) {
	for i := range lists {
		sb.WriteByte('\n')
		sb.WriteString(emojiNumbers[i])
		sb.WriteByte(' ')
		sb.WriteString(lists[i].Title)
	}
}

func writeMDV2LinkToBuilder(sb *strings.Builder, title, url string) {
	sb.WriteString("[")
	sb.WriteString(EscapeMarkdown(title))
	sb.WriteString("](")
	sb.WriteString(EscapeMarkdownLink(url))
	sb.WriteString(")")
}
func writeMDV2UserLinkToBuilder(sb *strings.Builder, user *f.User) {
	if user.Name != "" {
		writeMDV2LinkToBuilder(sb, user.Name, "t.me/"+user.Username)
	} else {
		writeMDV2LinkToBuilder(sb, "@"+user.Username, "t.me/"+user.Username)
	}
}

// EscapeMarkdown escapes special symbols for Telegram MarkdownV2 syntax
func EscapeMarkdown(s string) string {
	var result []rune
	for _, r := range s {
		if strings.ContainsRune("_*[]()~`>#+-=|{}.!\\", r) {
			result = append(result, '\\')
		}
		result = append(result, r)
	}
	return string(result)
}

func EscapeMarkdownLink(url string) string {
	var result []rune
	for _, r := range url {
		if strings.ContainsRune(")\\", r) {
			result = append(result, '\\')
		}
		result = append(result, r)
	}
	return string(result)
}

func GetUserState(id int64) UserCtx {
	ctx, ok := ctxStorage.Get(id)
	if ok {
		return ctx
	}
	user, err := f.GetUserById(db, id)
	if err != nil {
		return UserCtx{UserId: id}
	}
	res := UserCtx{
		UserId:   id,
		Name:     user.Name,
		Language: user.Language,
	}
	return res
}

func getListsSelectors(lists []f.List, btn MySelectorBtn) [][]tg.InlineButton {
	keyboard := make([][]tg.InlineButton, 0, (len(lists)+2)/3+2)
	if len(lists) == 4 {
		keyboard = append(keyboard,
			[]tg.InlineButton{btn.GetInlineButton(0, lists[0].ID), btn.GetInlineButton(1, lists[1].ID)},
			[]tg.InlineButton{btn.GetInlineButton(2, lists[2].ID), btn.GetInlineButton(3, lists[3].ID)},
		)
	} else {
		if len(lists) > 0 {
			row := make([]tg.InlineButton, 0, 3)
			for i := 0; i < len(lists) && i < 3; i++ {
				row = append(row, btn.GetInlineButton(i, lists[i].ID))
			}
			keyboard = append(keyboard, row)
		}
		if len(lists) > 4 {
			row := make([]tg.InlineButton, 0, 3)
			for i := 3; i < len(lists); i++ {
				row = append(row, btn.GetInlineButton(i, lists[i].ID))
			}
			keyboard = append(keyboard, row)
		}
	}
	return keyboard
}
