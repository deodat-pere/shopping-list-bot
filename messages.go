package main

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

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

func deleteMessage(ctx context.Context, b *bot.Bot, chatID any, msgID int) {
	_, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatID,
		MessageID: msgID,
	})

	if err != nil {
		fmt.Println(err)
	}
}

func reactMessage(ctx context.Context, b *bot.Bot, chatID any, msgID int, emoji string) {
	react := models.ReactionType{
		Type: models.ReactionTypeTypeEmoji,
		ReactionTypeEmoji: &models.ReactionTypeEmoji{
			Type:  models.ReactionTypeTypeEmoji,
			Emoji: emoji,
		},
	}

	var arr []models.ReactionType

	arr = append(arr, react)

	_, err := b.SetMessageReaction(ctx, &bot.SetMessageReactionParams{
		ChatID:    chatID,
		MessageID: msgID,
		Reaction:  arr,
	})

	if err != nil {
		fmt.Println(err)
	}
}
