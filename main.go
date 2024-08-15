package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"gopkg.in/yaml.v3"
)

var LMAP = ListMap{
	lists: make(map[any]List),
}

// Send any text message to the bot after the bot has been started

func main() {

	f, err := os.Open("config.yml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	opts := []bot.Option{
		bot.WithDefaultHandler(handler),
	}

	bot, err := bot.New(cfg.Token, opts...)
	if err != nil {
		panic(err)
	}

	bot.Start(ctx)
}

func handler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	action, err := ParseCommand(update)

	if err != nil {
		return
	}

	chatID := update.Message.Chat.ID

	switch action.action {
	case NewList:
		err := LMAP.newList(chatID, b, ctx)
		if err != nil {
			fmt.Println(err)
		}
	case CloseList:
		err := LMAP.closeList(chatID, b, ctx)
		if err != nil {
			fmt.Println(err)
		}
	case NewItem:
		msgID, err := LMAP.addItem(chatID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}

	case ModifyName:
		msgID, err := LMAP.modifyName(chatID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	case ModifyQuantity:
		msgID, err := LMAP.modifyQuantity(chatID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	case DeleteItem:
		msgID, err := LMAP.deleteItem(chatID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	}
}

func editMessage(ctx context.Context, b *bot.Bot, chatID any, msgID int) {
	list := LMAP.lists[chatID]
	newtext := "Liste de course:\n"
	for name, quantity := range list.items {
		newtext += "- `" + name + "`: " + quantity + "\n"
	}

	_, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatID,
		MessageID: msgID,
		Text:      newtext,
		ParseMode: "Markdown",
	})

	if err != nil {
		fmt.Println(err)
	}
}

type Config struct {
	Token    string `yaml:"token"`
	Username string `yaml:"username"`
}
