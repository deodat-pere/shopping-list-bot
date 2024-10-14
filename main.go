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

	args := os.Args[1]

	f, err := os.Open(args)
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

	botName, err := b.GetMyName(ctx, &bot.GetMyNameParams{})

	if err != nil {
		println("Couldn't get name")
		return
	}

	action, err := ParseCommand(update, &botName)

	if err != nil {
		return
	}

	chatID := update.Message.Chat.ID
	msgID := update.Message.ID

	switch action.action {
	case NewList:
		err := LMAP.newList(chatID, msgID, b, ctx)
		if err != nil {
			fmt.Println(err)
		}
	case CloseList:
		err := LMAP.closeList(chatID, msgID, b, ctx)
		if err != nil {
			fmt.Println(err)
		}
	case NewItem:
		msgID, err := LMAP.addItem(b, ctx, chatID, msgID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}

	case ModifyName:
		msgID, err := LMAP.modifyName(b, ctx, chatID, msgID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	case ModifyQuantity:
		msgID, err := LMAP.modifyQuantity(b, ctx, chatID, msgID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	case DeleteItem:
		msgID, err := LMAP.deleteItem(b, ctx, chatID, msgID, action)
		if err == nil {
			editMessage(ctx, b, chatID, msgID)
		}
	}
}

type Config struct {
	Token string `yaml:"token"`
}
