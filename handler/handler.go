package handler

import (
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"errors"
	"golang.org/x/sync/singleflight"
	"net/http"
	"strings"
	"time"
)

const cfToolsWebhookPrefix = "/cftools-webhook"

type webhookHandler struct {
	target         domain.Target
	servers        map[string]domain.Server
	filter         domain.FilterList
	history        domain.EventHistory
	logger         lager.Logger
	eventGroup     singleflight.Group
	executedEvents map[string]time.Time
}

func NewWebhookHandler(t domain.Target, s map[string]domain.Server, filter domain.FilterList, h domain.EventHistory, logger lager.Logger) *webhookHandler {
	handler := &webhookHandler{
		target:         t,
		servers:        s,
		filter:         filter,
		history:        h,
		logger:         logger,
		executedEvents: map[string]time.Time{},
	}

	go handler.invalidator(1 * time.Minute)

	return handler
}

func (h *webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	l := h.logger.Session("serve", lager.Data{"url": r.URL.String()})
	defer func() {
		err := r.Body.Close()
		if err != nil {
			l.Error("close-body", err)
		}
	}()

	s, err := h.serverFromRequest(r)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	e, err := domain.WebhookFromRequest(r)
	if err != nil {
		l.Error("parse-event", err)
		w.WriteHeader(500)
		return
	}

	l.Info("event", lager.Data{"event": e})
	if e.Event.Type == domain.EventVerification {
		w.WriteHeader(204)
		return
	}

	if !e.IsValidSignature(s.Secret) {
		l.Error("handle-event", errors.New("signature mismatch"))
		w.WriteHeader(403)
		return
	}

	_, err, _ = h.eventGroup.Do(e.Id, func() (interface{}, error) {
		if _, ok := h.executedEvents[e.Id]; ok {
			l.Info("de-duplicated-event", lager.Data{"id": e.Id})
			return nil, nil
		}

		err := h.onEvent(e, s.Name)
		h.executedEvents[e.Id] = time.Now()

		return nil, err
	})

	if err != nil {
		l.Error("handle-event", err)
		w.WriteHeader(500)
	}
}

func (h *webhookHandler) onEvent(e domain.WebhookEvent, serverName *string) error {
	if err := h.history.Save(e.Event); err != nil {
		return err
	}

	m, f, err := h.filter.MatchingFilters(h.history, e.Event)
	if err != nil {
		return err
	}
	if m && len(f) == 0 {
		return h.target.Relay(e.Event, nil, serverName)
	} else if m {
		for _, filter := range f {
			err = h.target.Relay(e.Event, &filter, serverName)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *webhookHandler) serverFromRequest(r *http.Request) (domain.Server, error) {
	l := h.logger.Session("server-from-request", lager.Data{"url": r.URL.String()})

	if !strings.HasPrefix(r.URL.String(), cfToolsWebhookPrefix) {
		l.Debug("not-found", lager.Data{"path": r.URL.String()})
		return domain.Server{}, errors.New("not a webhook event")
	}

	parts := strings.Split(strings.TrimPrefix(r.URL.String(), cfToolsWebhookPrefix), "/")
	sn := parts[len(parts)-1]
	s, ok := h.servers[sn]
	if !ok {
		l.Debug("not-found", lager.Data{"path": r.URL.String(), "serverName": sn})
		return domain.Server{}, errors.New("not a known server")
	}
	return s, nil
}

func (h *webhookHandler) invalidator(ttl time.Duration) {
	t := time.NewTicker(ttl)
	defer t.Stop()

	for range t.C {
		func() {
			defer func() {
				if err := recover(); err != nil {
					h.logger.Error("invalidator-recover", nil, lager.Data{"panic": err})
				}
			}()

			count := 0
			expired := time.Now().Add(2 * time.Minute)
			for k, v := range h.executedEvents {
				if v.Before(expired) {
					delete(h.executedEvents, k)
					count++
				}
			}
			if count != 0 {
				h.logger.Info("expire-executed-events", lager.Data{"invalidated": count, "remaining": len(h.executedEvents)})
			}
		}()
	}
}
