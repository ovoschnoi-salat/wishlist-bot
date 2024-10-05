package main

import (
	"strconv"
)

type State uint8

const (
	DefaultState State = iota
	ReadNameState
	ReadNewListTitleState
	ReadListNewTitleState
	ReadNewFriendUsernameState
	ReadNewWishTitleState
	ReadWishNewTitleState
	ReadWishNewDescriptionState
	ReadWishNewUrlState
	ReadWishNewPriceState
	StartState
)

type UserCtx struct {
	State             State
	LastMessageId     int
	ErrorMessageId    int
	Name              string
	Language          string
	UserId            int64
	FriendsPageNumber int64
	FriendId          int64
	ListId            int64
	ListPageNumber    int64
	WishId            int64
}

func (ctx *UserCtx) MessageSig() (messageID string, chatID int64) {
	return strconv.Itoa(ctx.LastMessageId), ctx.UserId
}

type ErrorCtx UserCtx

func (s *ErrorCtx) MessageSig() (messageID string, chatID int64) {
	return strconv.Itoa(s.ErrorMessageId), s.UserId
}
