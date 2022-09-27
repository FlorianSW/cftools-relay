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
	"text/template"
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

func formatType(f *domain.Filter) domain.FormatType {
	if f == nil || f.Format == nil {
		return domain.FormatTypeRich
	}
	return f.Format.Type
}

func richFormatParams(e domain.Event, f *domain.Filter, p *discordgo.WebhookParams) error {
	message := e.Message()
	color := domain.ColorDarkBlue
	if f != nil && f.Format != nil && f.Format.Parameters != nil {
		if m, ok := f.Format.Parameters["message"]; ok && m != "" {
			message = m.(string)
		}
		if c, ok := f.Format.Parameters["color"]; ok {
			color = domain.Color(c.(string)).Int()
		}
	}
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Message",
			Value:  message,
			Inline: false,
		},
	}
	for _, data := range e.Metadata() {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   data.K,
			Value:  data.V,
			Inline: true,
		})
	}
	p.Embeds = []*discordgo.MessageEmbed{
		{
			Color: color,
			Footer: &discordgo.MessageEmbedFooter{
				Text: "CFTools Relay by FlorianSW",
			},
			Provider: &discordgo.MessageEmbedProvider{
				URL:  "https://github.com",
				Name: "CFTools Relay",
			},
			Fields: fields,
		},
	}
	return nil
}

func textFormatParams(e domain.Event, f *domain.Filter, p *discordgo.WebhookParams) error {
	t := ""
	if f != nil && f.Format != nil && f.Format.Parameters != nil {
		if tpl, ok := f.Format.Parameters["template"]; ok {
			t = tpl.(string)
		}
	}
	if t == "" {
		for k, _ := range e.Values {
			t += " {{." + k + "}}"
		}
	}
	tpl, err := template.New("").Parse(t)
	if err != nil {
		p.Content = "Error in text template: " + err.Error()
		return err
	}
	var content bytes.Buffer
	err = tpl.Execute(&content, e.Values)
	if err != nil {
		return err
	}
	p.Content = content.String()
	return nil
}

func (t *discordTarget) Relay(e domain.Event, f *domain.Filter) error {
	l := t.logger.Session("relay", lager.Data{"event": e})
	params := discordgo.WebhookParams{
		Username: "CFTools-Relay",
	}

	var err error
	switch formatType(f) {
	case domain.FormatTypeRich:
		err = richFormatParams(e, f, &params)
	case domain.FormatTypeText:
		err = textFormatParams(e, f, &params)
	}
	if err != nil {
		return err
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
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	httpErr := errors.New("expected status code 2xx, got " + strconv.Itoa(res.StatusCode))
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	l.Error("discord", httpErr, lager.Data{"body": resBody})
	return err
}
