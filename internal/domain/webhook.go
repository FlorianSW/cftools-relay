package domain

import (
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
)

type EventFlavor = string

type UserJoin struct {
	CFToolsId   string `json:"cftools_id"`
	Name        string `json:"player_name"`
	IP          string `json:"player_ipv4"`
	BEGUID      string `json:"player_guid"`
	Steam64     int    `json:"player_steam64"`
	Country     string `json:"player_country"`
	CountryCode string `json:"player_country_code"`
	Vpn         string `json:"player_vpn"`
}

type UserLeave struct {
	CFToolsId string `json:"cftools_id"`
	Name      string `json:"player_name"`
	IP        string `json:"player_ipv4"`
	BEGUID    string `json:"player_guid"`
	Playtime  string `json:"player_playtime"`
}

type Death struct {
	CFToolsId      string `json:"victim_id"`
	Victim         string `json:"victim"`
	VictimPosition string `json:"victim_position"`
}

type Kill struct {
	VictimCFToolsId   string  `json:"victim_id"`
	Victim            string  `json:"victim"`
	MurdererCFToolsId string  `json:"murderer_id"`
	Murderer          string  `json:"murderer"`
	Weapon            string  `json:"weapon"`
	Distance          float32 `json:"distance"`
}

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
	err = json.Unmarshal(p, &parsed)
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
	default:
		return "Unknown event"
	}
}

func (e WebhookEvent) Metadata() (Metadata, error) {
	switch e.Event {
	case EventUserJoin:
		var payload UserJoin
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{
			{K: "Name", V: payload.Name},
			{K: "Steam ID", V: strconv.Itoa(payload.Steam64)},
			{K: "CFTools ID", V: payload.CFToolsId},
		}, nil
	case EventUserLeave:
		var payload UserLeave
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{
			{K: "Name", V: payload.Name},
			{K: "CFTools ID", V: payload.CFToolsId},
			{K: "Playtime", V: payload.Playtime},
		}, nil
	case EventPlayerKill:
		var payload Kill
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{
			{K: "Victim", V: payload.Victim},
			{K: "Victim CFTools ID", V: payload.VictimCFToolsId},
			{K: "Murderer", V: payload.Murderer},
			{K: "Murderer CFTools ID", V: payload.MurdererCFToolsId},
			{K: "Weapon", V: payload.Weapon},
			{K: "Distance in meter", V: fmt.Sprint(payload.Distance)},
		}, nil
	case EventPlayerDeathEnvironment:
		var payload Death
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{
			{K: "Name", V: payload.Victim},
			{K: "CFTools ID", V: payload.CFToolsId},
			{K: "Position", V: payload.VictimPosition},
		}, nil
	case EventPlayerDeathStarvation:
		var payload Death
		err := json.Unmarshal([]byte(e.Payload), &payload)
		if err != nil {
			return Metadata{}, err
		}
		return Metadata{
			{K: "Name", V: payload.Victim},
			{K: "CFTools ID", V: payload.CFToolsId},
			{K: "Position", V: payload.VictimPosition},
		}, nil
	default:
		return Metadata{}, nil
	}
}
