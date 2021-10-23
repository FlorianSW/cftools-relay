package adapter

import (
	"bytes"
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"encoding/json"
	"errors"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	ColorDarkBlue = 2123412
)

type discordTarget struct {
	webhookUrl string
	logger     lager.Logger
}

func NewDiscordTarget(webhookUrl string, logger lager.Logger) *discordTarget {
	return &discordTarget{
		webhookUrl: webhookUrl,
		logger:     logger,
	}
}

func (t *discordTarget) Relay(msg string, e domain.WebhookEvent) error {
	l := t.logger.Session("relay", lager.Data{"message": msg, "event": e})
	params := discordgo.WebhookParams{
		Username: "CFTools-Discord-Relay",
		Embeds: []*discordgo.MessageEmbed{
			{
				Color: ColorDarkBlue,
				Footer: &discordgo.MessageEmbedFooter{
					Text: "CFTools Relay by FlorianSW",
				},
				Provider: &discordgo.MessageEmbedProvider{
					URL:  "https://github.com",
					Name: "CFTools Relay",
				},
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Message",
						Value:  msg,
						Inline: false,
					},
					{
						Name:   "Event",
						Value:  e.Event,
						Inline: true,
					},
				},
			},
		},
	}

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", t.webhookUrl, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	defer func() {
		err := res.Body.Close()
		if err != nil {
			return
		}
	}()
	if res.StatusCode >= 200 || res.StatusCode <= 299 {
		return nil
	}
	err = errors.New("expected status code 2xx, got " + strconv.Itoa(res.StatusCode))
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	l.Error("discord", err, lager.Data{"body": resBody})
	return err
}
