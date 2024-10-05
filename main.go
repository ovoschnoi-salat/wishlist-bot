package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	f "wishlist_bot/repository"

	tg "gopkg.in/telebot.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"wishlist_bot/locales"
	"wishlist_bot/storage"
)

var db *gorm.DB
var ctxStorage = storage.NewStorage[UserCtx]()
var localizer *locales.I18n

func main() {
	dsn := os.Getenv("PG")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&f.Wish{}, &f.List{}, &f.User{})
	if err != nil {
		log.Fatalln(err)
	}

	err = ctxStorage.Load()
	if err != nil {
		log.Fatalln("Can't load users states:", err)
	}
	localizer, err = locales.NewLocalizer()
	if err != nil {
		log.Fatalln(err)
	}
	languagesKeyboard = make([][]tg.InlineButton, 0, len(locales.AvailableLocales))
	for _, lang := range locales.AvailableLocales {
		languagesKeyboard = append(languagesKeyboard, []tg.InlineButton{
			langSelectBtn.GetInlineButton(lang.Name, lang.Key),
		})
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		log.Fatalln("token is empty")
	}
	conf := tg.Settings{
		Token:       token,
		Poller:      &tg.LongPoller{Timeout: 10 * time.Second},
		Synchronous: false,
		Verbose:     false,
		OnError:     nil,
	}
	bot, err := tg.NewBot(conf)
	if err != nil {
		log.Fatalln(err)
	}

	registerHandlers(bot)

	signal.Ignore(syscall.SIGHUP)
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
	go func(bot *tg.Bot) {
		bot.Start()
	}(bot)
	<-c
	log.Println("signal received, stopping bot")
	bot.Stop()
	log.Println("bot stopped")
	if err := ctxStorage.Save(); err != nil {
		log.Println("error saving user states: " + err.Error())
	} else {
		log.Println("user states saved")
	}
}

func registerHandlers(bot *tg.Bot) {
	bot.Handle("/start", startHandler)
	bot.Handle("/help", func(c tg.Context) error {
		// TODO
		panic("not implemented")
	})
	bot.Handle(tg.OnText, textHandler)
	registerAllHandlers(bot)
}

var textStateHandlers = make(map[State]func(tg.Context) error)

func textHandler(c tg.Context) error {
	err := c.Delete()
	if err != nil {
		sendAlert(c, "error deleting message: "+err.Error())
	}
	ctx := GetUserState(c.Chat().ID)
	if fn, ok := textStateHandlers[ctx.State]; ok {
		return fn(c)
	}
	return sendMainMenu(c, &ctx)
}
