package main

import (
	"context"
	"os"
	"strconv"

	"github.com/hunternsk/ninjabot/plot"
	"github.com/hunternsk/ninjabot/plot/indicator"

	"github.com/hunternsk/ninjabot"
	"github.com/hunternsk/ninjabot/examples/strategies"
	"github.com/hunternsk/ninjabot/exchange"
	"github.com/hunternsk/ninjabot/storage"
	"github.com/hunternsk/ninjabot/tools/log"
)

// This example shows how to use NinjaBot with a simulation with a fake exchange
// A peperwallet is a wallet that is not connected to any exchange, it is a simulation with live data (realtime)
func main() {
	var (
		ctx             = context.Background()
		telegramToken   = os.Getenv("TELEGRAM_TOKEN")
		telegramUser, _ = strconv.Atoi(os.Getenv("TELEGRAM_USER"))
	)

	settings := ninjabot.Settings{
		Pairs: []string{
			"BTCUSDT",
			"ETHUSDT",
			"BNBUSDT",
			"LTCUSDT",
		},
		Telegram: ninjabot.TelegramSettings{
			Enabled: telegramToken != "" && telegramUser != 0,
			Token:   telegramToken,
			Users:   []int{telegramUser},
		},
	}

	// Use binance for realtime data feed
	binance, err := exchange.NewBinance(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// creating a storage to save trades
	storage, err := storage.FromMemory()
	if err != nil {
		log.Fatal(err)
	}

	// creating a paper wallet to simulate an exchange waller for fake operataions
	// paper wallet is simulation of a real exchange wallet
	paperWallet := exchange.NewPaperWallet(
		ctx,
		"USDT",
		exchange.WithPaperFee(0.001, 0.001),
		exchange.WithPaperAsset("USDT", 10000),
		exchange.WithDataFeed(binance),
	)

	// initializing my strategy
	strategy := new(strategies.CrossEMA)

	chart, err := plot.NewChart(
		plot.WithCustomIndicators(
			indicator.EMA(8, "red"),
			indicator.SMA(21, "blue"),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	// initializer ninjabot
	bot, err := ninjabot.NewBot(
		ctx,
		settings,
		paperWallet,
		strategy,
		ninjabot.WithStorage(storage),
		ninjabot.WithPaperWallet(paperWallet),
		ninjabot.WithCandleSubscription(chart),
		ninjabot.WithOrderSubscription(chart),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err := chart.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = bot.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
