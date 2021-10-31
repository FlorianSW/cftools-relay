package handler

import (
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"errors"
	"net/http"
	"strings"
)

const cfToolsWebhookPrefix = "/cftools-webhook"

type webhookHandler struct {
	target  domain.Target
	servers map[string]domain.Server
	filter  domain.FilterList
	history domain.EventHistory
	logger  lager.Logger
}

func NewWebhookHandler(t domain.Target, s map[string]domain.Server, filter domain.FilterList, h domain.EventHistory, logger lager.Logger) *webhookHandler {
	return &webhookHandler{
		target:  t,
		servers: s,
		filter:  filter,
		history: h,
		logger:  logger,
	}
}

func (h webhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if err = h.onEvent(e); err != nil {
		l.Error("handle-event", err)
		w.WriteHeader(500)
	}
}

func (h webhookHandler) onEvent(e domain.WebhookEvent) error {
	if err := h.history.Save(e.Event); err != nil {
		return err
	}

	m, f, err := h.filter.MatchingFilter(h.history, e.Event)
	if err != nil {
		return err
	}
	if m {
		return h.target.Relay(e.Event, f)
	}
	return nil
}

func (h webhookHandler) serverFromRequest(r *http.Request) (domain.Server, error) {
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
