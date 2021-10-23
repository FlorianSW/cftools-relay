package domain

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	FlavorCftools = "WebHookFlavor.CFTOOLS"
	FlavorDiscord = "WebHookFlavor.DISCORD"

	EventVerification = "verification"
)

type EventFlavor = string

type WebhookEvent struct {
	ShardId   int
	Flavor    EventFlavor
	Event     string
	Id        string
	Signature string
	Payload   string
}

func WebhookFromRequest(r *http.Request) (WebhookEvent, error) {
	p, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return WebhookEvent{}, err
	}
	shardId, err := strconv.Atoi(r.Header.Get("X-Hephaistos-Shard"))
	if err != nil {
		return WebhookEvent{}, err
	}
	return WebhookEvent{
		ShardId:   shardId,
		Flavor:    r.Header.Get("X-Hephaistos-Flavor"),
		Event:     r.Header.Get("X-Hephaistos-Event"),
		Id:        r.Header.Get("X-Hephaistos-Delivery"),
		Signature: r.Header.Get("X-Hephaistos-Signature"),
		Payload:   string(p),
	}, nil
}
