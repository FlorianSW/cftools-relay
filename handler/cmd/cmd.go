package main

import (
	"cftools-relay/handler"
	"cftools-relay/internal"
	"cftools-relay/internal/adapter"
	"code.cloudfoundry.org/lager"
	"net/http"
	"os"
	"strconv"
)

func main() {
	logger := lager.NewLogger("cftools-relay")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.INFO))

	c, err := internal.NewConfig("./config.json", logger)
	if err != nil {
		logger.Fatal("config", err)
	}

	target := adapter.NewDiscordTarget(c.Discord.WebhookUrl, logger)
	history, err := adapter.NewEventRepository(c.History.StoragePath)
	if err != nil {
		logger.Fatal("event-history", err)
	}
	h := handler.NewWebhookHandler(target, c.Servers, c.Filter, history, logger)

	logger.Info("start-listener", lager.Data{"port": c.Port})
	err = http.ListenAndServe(":"+strconv.Itoa(c.Port), h)
	if err != nil {
		logger.Fatal("start-listener", err)
	}
}
