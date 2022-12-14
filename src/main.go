package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

const (
	emojiStatusChilling = int64(5206513124630341906)
	emojiStatusEating   = int64(5206270085315961515)
	emojiStatusSleeping = int64(5445211704541586076)
	emojiStatusWorking  = int64(5443132326189996902)
)

const (
	location = "Europe/Moscow"
)

func main() {
	client, err := telegram.ClientFromEnvironment(telegram.Options{})
	if err != nil {
		panic(fmt.Errorf("failed to create client: %w", err))
	}

	flow := auth.NewFlow(terminalAuth{}, auth.SendCodeOptions{})
	if err := client.Run(context.Background(), func(ctx context.Context) error {
		if err := client.Auth().IfNecessary(ctx, flow); err != nil {
			return err
		}
		api := client.API()
		return RunBackground(ctx, api)
	}); err != nil {
		panic(err)
	}
}

func RunBackground(ctx context.Context, api *tg.Client) error {
	prevEmojiStatus := int64(0)
	loc, _ := time.LoadLocation(location)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		case <-time.After(time.Second):
		}

		now := time.Now().In(loc)
		hour := now.Hour()
		emojiStatus := getEmojiStatusByHour(hour)
		if emojiStatus == prevEmojiStatus {
			continue
		}
		fmt.Println("need to update status")
		prevEmojiStatus = emojiStatus

		_, err := api.AccountUpdateEmojiStatus(ctx, &tg.EmojiStatus{
			DocumentID: emojiStatus,
		})
		if err != nil {
			return err
		}
		fmt.Println("status updated")
	}
}

func getEmojiStatusByHour(hour int) int64 {
	if 0 <= hour && hour < 8 {
		return emojiStatusSleeping
	} else if 8 <= hour && hour < 10 {
		return emojiStatusEating
	} else if 10 <= hour && hour < 14 {
		return emojiStatusWorking
	} else if hour == 14 {
		return emojiStatusEating
	} else if 15 <= hour && hour < 18 {
		return emojiStatusWorking
	} else if hour == 18 {
		return emojiStatusEating
	} else {
		return emojiStatusChilling
	}
}
