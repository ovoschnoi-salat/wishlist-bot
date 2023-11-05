package main

import (
	"bytes"
	"errors"
	"github.com/jackc/pgx/v5"
	tg "gopkg.in/telebot.v3"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"wishlist_bot/pg"
)

var (
	mainMenu = &tg.ReplyMarkup{ResizeKeyboard: true,
		ReplyKeyboard: [][]tg.ReplyButton{{addWishBtn, shareBtn}, {showWishlistBtn}, {showFriendsListBtn}}}
	addWishBtn         = tg.ReplyButton{Text: "🆕 Add wish"}
	shareBtn           = tg.ReplyButton{Text: "🔗 Share"}
	showWishlistBtn    = tg.ReplyButton{Text: "📜 Show my wishlist"}
	showFriendsListBtn = tg.ReplyButton{Text: "🎁 Show friends list"}

	cancelKeyboard = &tg.ReplyMarkup{InlineKeyboard: [][]tg.InlineButton{{btnCancel}}, RemoveKeyboard: true}
	btnCancel      = tg.InlineButton{Unique: "cancel", Text: "Отмена"}
)

var pgStorage pg.Storage

func main() {
	var err error
	pgStorage, err = pg.NewPgStorage(os.Getenv("PG"))
	if err != nil {
		log.Fatalln(err)
	}
	conf := tg.Settings{
		Token:       os.Getenv("TOKEN"),
		Poller:      &tg.LongPoller{Timeout: 10 * time.Second},
		Synchronous: false,
		Verbose:     false,
		OnError:     nil,
	}
	bot, err := tg.NewBot(conf)
	if err != nil {
		log.Fatalln(err)
	}

	bot.Handle("/start", func(c tg.Context) error {
		conn, err := pgStorage.Acquire()
		if err != nil {
			return sendError(c, err)
		}
		defer conn.Release()
		_, err = conn.GetUserByUsername(c.Chat().Username)
		if errors.Is(err, pgx.ErrNoRows) {
			err = conn.AddUser(c.Chat().ID, c.Chat().Username)
			if err != nil {
				return sendError(c, err)
			}
		} else if err != nil {
			return sendError(c, err)
		}
		if c.Message().Payload != "" {
			decodeString, err := b64url.DecodeString(c.Message().Payload)
			if err != nil {
				return c.Send("неправильный payload: " + c.Message().Payload)
			}
			payload := strings.Split(string(decodeString), "+")
			if len(payload) != 2 {
				return c.Send("неправильный payload: " + c.Message().Payload)
			}
			friend, err := conn.GetUserByUsername(payload[0])
			if err != nil {
				return c.Send("cant find user: " + err.Error())
			}
			hash := calcHash(nil, friend.Id, friend.Username)
			if !bytes.Equal([]byte(payload[1]), hash) {
				return c.Send("user not found: " + err.Error())
			}
			if c.Chat().ID == friend.Id {
				return c.Send("Отправь эту ссылку своим друзьяи, чтобы они смогли увидеть твой список желаний.")
			}
			err = conn.AddFriend(c.Chat().ID, friend.Id)
			if err != nil {
				return c.Send("error adding to friends: " + err.Error())
			}
			// todo show friends wishlist
			return c.Send("user added to friends: "+payload[0], mainMenu)
		}
		return c.Send("Привет, это бот списка желаний\n"+
			"Просто начни добавлять свои желания!", mainMenu)
		// todo add language choice
	})

	bot.Handle("/help", func(c tg.Context) error {
		return c.Send("todo", mainMenu)
	})

	bot.Handle(&btnCancel, cancel)

	registerAddWishHandlers(bot)
	registerShowListHandlers(bot)
	registerShowWishHandlers(bot)
	registerShowFriendsListHandlers(bot)
	registerShareLinkHandlers(bot)
	registerFriendsRequestHandlers(bot)

	bot.Handle(tg.OnText, func(c tg.Context) error {
		state, _ := userStates.GetUserState(c.Chat().ID)
		switch state.State {
		case DefaultState:
			return unknownCommand(c)
		case NewWishState:
			return readWish(c)
		case ReadUrlState:
			return readUrl(c)
		case ReadNewTitleState:
			return readNewTitle(c)
		case NewFriendState:
			return unknownCommand(c)
		default:
			return readUnexpectedWish(c)
		}
	})
	signal.Ignore(syscall.SIGHUP)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func(bot *tg.Bot) {
		// todo maybe add
		<-c
		bot.Stop()
		pgStorage.Close()
	}(bot)
	bot.Start()
}

func cancel(c tg.Context) error {
	_ = userStates.DeleteUserState(c.Chat().ID)
	return c.Edit("Отменено", &tg.ReplyMarkup{})
}

func unknownCommand(c tg.Context) error {
	return c.Send("Простите, бот немного запутался и не смог понять, что вы сейчас сделали.\nПопробуйте заново.",
		mainMenu)
}

func sendError(c tg.Context, err error) error {
	return c.Send("Возникла ошибка: "+err.Error(), mainMenu)
}
