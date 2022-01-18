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
	"time"
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

	FieldCfToolsId         = "cftools_id"
	FieldPlayerId          = "player_id"
	FieldVictimCfToolsId   = "victim_id"
	FieldMurdererCfToolsId = "murderer_id"
)

var possibleMetadata = map[string]string{
	"player_name":          "Name",
	"player_steam64":       "Steam ID",
	FieldCfToolsId:         "CFTools ID",
	"player_playtime":      "Playtime",
	"victim":               "Victim",
	"victim_position":      "Victim Position",
	FieldVictimCfToolsId:   "Victim CFTools ID",
	"murderer":             "Murderer",
	FieldMurdererCfToolsId: "Murderer CFTools ID",
	"weapon":               "Weapon",
	"damage":               "Damage points",
	"distance":             "Distance in meter",
	"item":                 "Item",
}

type Server struct {
	Secret string `json:"secret"`
}

type EventFlavor = string

type WebhookEvent struct {
	ShardId   int
	Flavor    EventFlavor
	Id        string
	Signature string
	Payload   string
	Event     Event
}

type Event struct {
	Type      string
	Timestamp time.Time
	Values    map[string]interface{}
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
		ShardId:   shardId,
		Flavor:    r.Header.Get("X-Hephaistos-Flavor"),
		Id:        r.Header.Get("X-Hephaistos-Delivery"),
		Signature: r.Header.Get("X-Hephaistos-Signature"),
		Payload:   string(p),
		Event: Event{
			Type:      r.Header.Get("X-Hephaistos-Event"),
			Timestamp: time.Now(),
			Values:    parsed,
		},
	}, nil
}

func (e WebhookEvent) IsValidSignature(secret string) bool {
	if e.Event.Type == EventVerification {
		return true
	}
	a := sha256.New()
	a.Write([]byte(e.Id))
	a.Write([]byte(secret))
	r := a.Sum(nil)

	return hex.EncodeToString(r) == e.Signature
}

func (e Event) Message() string {
	switch e.Type {
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
		return fmt.Sprintf("Event: %s", e.Type)
	}
}

func (e Event) Metadata() Metadata {
	var m Metadata

	for key, label := range possibleMetadata {
		if v, ok := e.Values[key]; ok {
			m = append(m, Data{
				K: label,
				V: stringutil.Itos(v),
			})
		}
	}
	return m
}

func (e Event) CFToolsId() *string {
	for _, f := range []string{FieldCfToolsId, FieldMurdererCfToolsId, FieldPlayerId, FieldVictimCfToolsId} {
		possibleId, ok := e.Values[f]
		if id, isString := possibleId.(string); ok && isString {
			return &id
		}
	}
	return nil
}
