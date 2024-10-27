package main

import (
	"context"
	"errors"
	"time"

	"github.com/go-telegram/bot"
)

// A shopping list is a map between item names and their quantity
type List struct {
	items map[string]string
	msgID int
}

func (list List) addItem(name string, quantity string) {
	_, prs := list.items[name]

	if !prs {
		list.items[name] = quantity
	}
}

func (list List) modifyName(name string, rename string) {
	quantity, prs := list.items[name]

	if prs {
		delete(list.items, name)
		list.items[rename] = quantity
	}
}

func (list List) modifyQuantity(name string, quantity string) {
	_, prs := list.items[name]

	if prs {
		list.items[name] = quantity
	}
}

func (list List) deleteItem(name string) {
	delete(list.items, name)
}

type ListMap struct {
	lists map[any]List
}

func (lmap ListMap) closeList(chatID any, callingMsgID int, b *bot.Bot, ctx context.Context) error {
	list, prs := lmap.lists[chatID]
	if !prs {
		reactMessage(ctx, b, chatID, callingMsgID, "🤔")
		return errors.New("no list to close")
	}

	msgID := list.msgID

	_, err := b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
		ChatID:    chatID,
		MessageID: msgID,
	})

	if err != nil {
		return err
	}

	delete(lmap.lists, chatID)

	reactMessage(ctx, b, chatID, callingMsgID, "👍")

	time.Sleep(2 * 1000 * 1000 * 1000) // Sleep for 2 seconds

	deleteMessage(ctx, b, chatID, callingMsgID)

	return nil
}

func (lmap ListMap) newList(chatID any, callingMsgID int, b *bot.Bot, ctx context.Context) error {
	list, prs := lmap.lists[chatID]
	if prs {
		msgID := list.msgID

		_, err := b.UnpinChatMessage(ctx, &bot.UnpinChatMessageParams{
			ChatID:    chatID,
			MessageID: msgID,
		})

		if err != nil {
			return err
		}
	}

	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   "Liste de course:",
	})

	lmap.lists[chatID] = List{
		items: make(map[string]string),
		msgID: msg.ID,
	}

	if err != nil {
		return err
	}

	_, err = b.PinChatMessage(ctx, &bot.PinChatMessageParams{
		ChatID:              chatID,
		MessageID:           msg.ID,
		DisableNotification: true,
	})

	deleteMessage(ctx, b, chatID, callingMsgID)

	return err
}

func (lmap ListMap) addItem(b *bot.Bot, ctx context.Context, chatID any, msgID int, action Action) (int, error) {
	list, prs := lmap.lists[chatID]

	if !prs {
		reactMessage(ctx, b, chatID, msgID, "🤔")
		return 0, errors.New("no list to edit")
	}

	name, quantity := action.arg1, action.arg2
	list.addItem(name, quantity)
	deleteMessage(ctx, b, chatID, msgID)

	return list.msgID, nil
}

func (lmap ListMap) modifyName(b *bot.Bot, ctx context.Context, chatID any, msgID int, action Action) (int, error) {
	list, prs := lmap.lists[chatID]

	if !prs {
		reactMessage(ctx, b, chatID, msgID, "🤔")
		return 0, errors.New("no list to edit")
	}

	name, newName := action.arg1, action.arg2
	list.modifyName(name, newName)
	deleteMessage(ctx, b, chatID, msgID)

	return list.msgID, nil
}

func (lmap ListMap) modifyQuantity(b *bot.Bot, ctx context.Context, chatID any, msgID int, action Action) (int, error) {
	list, prs := lmap.lists[chatID]

	if !prs {
		reactMessage(ctx, b, chatID, msgID, "🤔")
		return 0, errors.New("no list to edit")
	}

	name, quantity := action.arg1, action.arg2
	list.modifyQuantity(name, quantity)
	deleteMessage(ctx, b, chatID, msgID)

	return list.msgID, nil
}

func (lmap ListMap) deleteItem(b *bot.Bot, ctx context.Context, chatID any, msgID int, action Action) (int, error) {
	list, prs := lmap.lists[chatID]

	if !prs {
		reactMessage(ctx, b, chatID, msgID, "🤔")
		return 0, errors.New("no list to edit")
	}

	name := action.arg1
	list.deleteItem(name)
	deleteMessage(ctx, b, chatID, msgID)

	return list.msgID, nil
}
