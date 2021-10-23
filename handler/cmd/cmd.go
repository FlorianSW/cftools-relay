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

	h := handler.NewWebhookHandler(adapter.NewDiscordTarget(c.Discord.WebhookUrl, logger), logger)

	logger.Info("start-listener", lager.Data{"port": c.Port})
	err = http.ListenAndServe(":"+strconv.Itoa(c.Port), h)
	if err != nil {
		logger.Fatal("start-listener", err)
	}
}
