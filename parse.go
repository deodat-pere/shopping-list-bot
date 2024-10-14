package main

import (
	"errors"
	"strings"

	"github.com/go-telegram/bot/models"
)

const (
	NewList = iota
	CloseList
	NewItem
	ModifyName
	ModifyQuantity
	DeleteItem
)

type Action struct {
	action int
	arg1   string
	arg2   string
}

func ParseCommand(update *models.Update, botName *models.BotName) (Action, error) {
	text := update.Message.Text
	words := strings.Split(text, "/")

	if len(words) < 2 {
		return Action{}, errors.New("no command provided")
	}

	switch strings.TrimSpace(words[1]) {
	case "new", "new@" + botName.Name:
		return Action{
			action: NewList,
			arg1:   "",
			arg2:   "",
		}, nil
	case "close", "close@" + botName.Name:
		return Action{
			action: CloseList,
			arg1:   "",
			arg2:   "",
		}, nil
	case "add", "add@" + botName.Name:
		if len(words) < 4 {
			return Action{}, errors.New("not enough arguments")
		}
		return Action{
			action: NewItem,
			arg1:   strings.TrimSpace(words[2]),
			arg2:   strings.TrimSpace(words[3]),
		}, nil
	case "modifyname", "modifyname@" + botName.Name:
		if len(words) < 4 {
			return Action{}, errors.New("not enough arguments")
		}
		return Action{
			action: ModifyName,
			arg1:   strings.TrimSpace(words[2]),
			arg2:   strings.TrimSpace(words[3]),
		}, nil
	case "modifyquantity", "modifyquantity@" + botName.Name:
		if len(words) < 4 {
			return Action{}, errors.New("not enough arguments")
		}
		return Action{
			action: ModifyQuantity,
			arg1:   strings.TrimSpace(words[2]),
			arg2:   strings.TrimSpace(words[3]),
		}, nil
	case "delete", "delete@" + botName.Name:
		if len(words) < 3 {
			return Action{}, errors.New("not enough arguments")
		}
		return Action{
			action: DeleteItem,
			arg1:   strings.TrimSpace(words[2]),
			arg2:   "",
		}, nil
	default:
		return Action{}, errors.New("invalid command")
	}
}
