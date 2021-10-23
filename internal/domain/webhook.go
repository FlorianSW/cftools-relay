package domain

import (
	"bytes"
	"cftools-relay/internal/stringutil"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

const (
	FlavorCftools = "WebHookFlavor.CFTOOLS"
	FlavorDiscord = "WebHookFlavor.DISCORD"

	EventVerification           = "verification"
	EventUserJoin               = "user.join"
	EventUserLeave              = "user.leave"
	EventPlayerPlace            = "player.place"
	EventPlayerDeathStarvation  = "player.death_starvation"
	EventPlayerDeathEnvironment = "player.death_environment"
	EventPlayerKill             = "player.kill"
	EventPlayerDamage           = "player.damage"
)

type EventFlavor = string

type WebhookEvent struct {
	ShardId       int
	Flavor        EventFlavor
	Event         string
	Id            string
	Signature     string
	Payload       string
	ParsedPayload map[string]interface{}
}

type Metadata []Data

type Data struct {
	K, V string
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
	var parsed map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(p))
	d.UseNumber()
	err = d.Decode(&parsed)
	if err != nil {
		return WebhookEvent{}, err
	}
	return WebhookEvent{
		ShardId:       shardId,
		Flavor:        r.Header.Get("X-Hephaistos-Flavor"),
		Event:         r.Header.Get("X-Hephaistos-Event"),
		Id:            r.Header.Get("X-Hephaistos-Delivery"),
		Signature:     r.Header.Get("X-Hephaistos-Signature"),
		Payload:       string(p),
		ParsedPayload: parsed,
	}, nil
}

func (e WebhookEvent) IsValidSignature(secret string) bool {
	if e.Event == EventVerification {
		return true
	}
	a := sha256.New()
	a.Write([]byte(e.Id))
	a.Write([]byte(secret))
	r := a.Sum(nil)

	return hex.EncodeToString(r) == e.Signature
}

func (e WebhookEvent) Message() string {
	switch e.Event {
	case EventUserJoin:
		return "Player connected."
	case EventUserLeave:
		return "Player disconnected."
	case EventPlayerKill:
		return "Player was killed."
	case EventPlayerDeathEnvironment:
		return "Player was killed by the environment."
	case EventPlayerDeathStarvation:
		return "Player died from starvation."
	case EventPlayerDamage:
		return "Player injured another player."
	case EventPlayerPlace:
		return "Player played an item."
	default:
		return fmt.Sprintf("Event: %s", e.Event)
	}
}

func (e WebhookEvent) Metadata() Metadata {
	var m Metadata

	if v, ok := e.ParsedPayload["player_name"]; ok {
		m = append(m, Data{
			K: "Name",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["player_steam64"]; ok {
		m = append(m, Data{
			K: "Steam ID",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["cftools_id"]; ok {
		m = append(m, Data{
			K: "CFTools ID",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["player_playtime"]; ok {
		m = append(m, Data{
			K: "CFTools ID",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["victim"]; ok {
		m = append(m, Data{
			K: "Victim",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["victim_position"]; ok {
		m = append(m, Data{
			K: "Victim Potision",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["victim_id"]; ok {
		m = append(m, Data{
			K: "Victim CFTools ID",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["murderer"]; ok {
		m = append(m, Data{
			K: "Murderer",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["murderer_id"]; ok {
		m = append(m, Data{
			K: "Murderer CFTools ID",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["weapon"]; ok {
		m = append(m, Data{
			K: "Weapon",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["damage"]; ok {
		m = append(m, Data{
			K: "Damage points",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["distance"]; ok {
		m = append(m, Data{
			K: "Distance in meter",
			V: stringutil.Itos(v),
		})
	}
	if v, ok := e.ParsedPayload["item"]; ok {
		m = append(m, Data{
			K: "Item",
			V: stringutil.Itos(v),
		})
	}
	return m
}
