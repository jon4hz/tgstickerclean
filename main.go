package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
)

var (
	defaultTimeout = time.Second * 5
	timeout        time.Duration
)

func init() {
	t := os.Getenv("TIMEOUT")
	if t == "" {
		timeout = defaultTimeout
	}
	var err error
	timeout, err = time.ParseDuration(t)
	if err != nil {
		panic(err)
	}
}

func main() {
	b, err := gotgbot.NewBot(os.Getenv("BOT_TOKEN"), &gotgbot.BotOpts{
		Client:      http.Client{},
		GetTimeout:  gotgbot.DefaultGetTimeout,
		PostTimeout: gotgbot.DefaultPostTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	u := ext.NewUpdater(nil)
	u.Dispatcher.AddHandler(handlers.NewMessage(nil, stickerHandler))
	err = u.StartPolling(b, &ext.PollingOpts{
		DropPendingUpdates: true,
		GetUpdatesOpts: gotgbot.GetUpdatesOpts{
			AllowedUpdates: []string{
				"message",
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s has been started...", b.User.Username)
	u.Idle()
}

func stickerHandler(b *gotgbot.Bot, ctx *ext.Context) error {
	if ctx.EffectiveMessage.Sticker != nil {
		go stickerRemover(b, ctx.EffectiveChat.Id, ctx.EffectiveMessage.MessageId)
	}
	return nil
}

func stickerRemover(b *gotgbot.Bot, chatID, msgID int64) {
	time.Sleep(timeout)
	_, err := b.DeleteMessage(chatID, msgID)
	if err != nil {
		log.Printf("error deleting message: %s", err)
	}
}
