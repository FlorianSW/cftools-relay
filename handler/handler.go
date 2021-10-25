package handler

import (
	"cftools-relay/internal/domain"
	"code.cloudfoundry.org/lager"
	"errors"
	"net/http"
)

type webhookHandler struct {
	target domain.Target
	secret string
	filter  domain.FilterList
	history domain.EventHistory
	logger  lager.Logger
}

func NewWebhookHandler(t domain.Target, s string, filter domain.FilterList, h domain.EventHistory, logger lager.Logger) *webhookHandler {
	return &webhookHandler{
		target:  t,
		secret:  s,
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

	if r.URL.String() != "/cftools-webhook" {
		l.Debug("not-found", lager.Data{"path": r.URL.String()})
		w.WriteHeader(404)
	} else {
		e, err := domain.WebhookFromRequest(r)
		if err != nil {
			l.Error("parse-event", err)
			w.WriteHeader(500)
			return
		}
		l.Info("event", lager.Data{"event": e})
		if e.Event.Type != domain.EventVerification {
			err = h.onEvent(e)
			if err != nil {
				l.Error("handle-event", err)
				w.WriteHeader(500)
				return
			}
		}
		w.WriteHeader(204)
	}
}

func (h webhookHandler) onEvent(e domain.WebhookEvent) error {
	if !e.IsValidSignature(h.secret) {
		return errors.New("signature-mismatch")
	}
	err := h.history.Save(e.Event)
	if err != nil {
		return err
	}
	if h.filter.Matches(h.history, e.Event) {
		return h.target.Relay(e.Event)
	}
	return nil
}
