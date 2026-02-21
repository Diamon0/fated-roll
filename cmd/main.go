package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Diamon0/fated-roll/internal"
	"github.com/fluxergo/fluxergo"
	"github.com/fluxergo/fluxergo/bot"
)

var (
	token string
	logger *slog.Logger
)

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	errorFound := false

	if foundToken := os.Getenv("BOT_TOKEN"); foundToken != "" {
		token = foundToken
	} else {
		logger.Error("BOT_TOKEN is not defined")
		errorFound = true
	}

	if errorFound {
		os.Exit(1)
	}
}

func main() {
	client, err := fluxergo.New(token, 
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(internal.MessageHandler),
	)
	if err != nil {
		logger.Error("Error while creating bot client", slog.Any("error", err))
		os.Exit(1)
	}

	bg := context.Background()
	ctx, stop := signal.NotifyContext(bg, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
	defer stop()

	if err = client.OpenGateway(ctx); err != nil {
		logger.Error("Error while connecting to Fluxer", slog.Any("error", err))
		stop()
		os.Exit(1)
	}

	slog.Info("Fated Roll Bot is running; CTRL-C to exit.")
	<-ctx.Done()
}
