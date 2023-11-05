package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"golang.org/x/crypto/sha3"
	tg "gopkg.in/telebot.v3"
	"regexp"
	"strconv"
	"strings"
)

var (
	emojiNumbers = []string{
		"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣",
	}
	b64    = base64.StdEncoding
	b64url = base64.RawURLEncoding
	prev   = "<<"
	next   = ">>"
	rxURL  = regexp.MustCompile(URL)
	URL    = "^(https?:\\/\\/)?" +
		"(([a-zA-Z0-9]([a-zA-Z0-9-]{0,253}[a-zA-Z0-9])?\\.)+[a-zA-Z0-9]([a-zA-Z0-9-]{0,253}[a-zA-Z0-9])?|" +
		"([а-яА-Я0-9]([а-яА-Я0-9-]{0,253}[а-яА-Я0-9])?\\.)+(рф|РФ)|" +
		"((25[0-5]|(2[0-4]|1\\d|[1-9]|)\\d)\\.?\\b){4}" +
		"(:(6(553[0-5]|55[0-2][0-9]|5[0-4][0-9]{2}|[0-4][0-9]{3})|[1-5][0-9]{4}|[1-9][0-9]{0,3}))?)" +
		"(\\/([a-zA-Z0-9-._~:@!$&'()*+,;=]|%[a-fA-f0-9]{2})*)*" +
		"(\\?([a-zA-Z0-9-._~:@!$&'()*+,;=/\\?]|%[a-fA-f0-9]{2})*)?" +
		"(#([a-zA-Z0-9-._~:@!$&'()*+,;=/\\?]|%[a-fA-f0-9]{2})*)?$"
)

func addPageNumber(str *strings.Builder, index, all int64) {
	str.WriteString("\n страница ")
	str.WriteString(strconv.FormatInt(index+1, 10))
	str.WriteRune('/')
	str.WriteString(strconv.FormatInt(all, 10))
}

func getNewBtn(text, unique, data string) tg.InlineButton {
	return tg.InlineButton{Text: text, Unique: unique, Data: data}
}

func getBtnWithData(btn tg.InlineButton, data string) tg.InlineButton {
	return getNewBtn(btn.Text, btn.Unique, data)
}

func getBtnWithId(btn tg.InlineButton, id int64) tg.InlineButton {
	return getBtnWithData(btn, strconv.FormatInt(id, 16))
}

func getBtnWithIdAndData(btn tg.InlineButton, id, data int64) tg.InlineButton {
	return getBtnWithData(btn, strconv.FormatInt(id, 16)+":"+strconv.FormatInt(data, 16))
}

func getNewBtnWithId(text, unique string, id int64) tg.InlineButton {
	return getNewBtn(text, unique, strconv.FormatInt(id, 16))
}

func getNewBtnWithIdAndData(text, unique string, id, data int64) tg.InlineButton {
	return getNewBtn(text, unique, strconv.FormatInt(id, 16)+":"+strconv.FormatInt(data, 16))
}

func getId(c tg.Context) (int64, error) {
	return strconv.ParseInt(c.Data(), 16, 0)
}

func getIdAndData(c tg.Context) (userId, wishId int64, err error) {
	array := strings.Split(c.Data(), ":")
	if len(array) != 2 {
		return 0, 0, errors.New("wrong data passed:" + c.Data())
	}
	userId, err = strconv.ParseInt(array[0], 16, 0)
	if err != nil {
		return
	}
	wishId, err = strconv.ParseInt(array[1], 16, 0)
	return
}

func calcHash(array []byte, userId int64, username string) []byte {
	h := sha3.New224()
	id := make([]byte, 8)
	binary.LittleEndian.PutUint64(id, uint64(userId))
	h.Write(id)
	h.Write([]byte(username))
	binary.BigEndian.PutUint64(id, uint64(userId))
	h.Write(id)
	return h.Sum(array)
}

func getEndpointFromUnique(unique string) string {
	return "\f" + unique
}
