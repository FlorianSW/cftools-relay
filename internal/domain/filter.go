package domain

import (
	"cftools-relay/internal/stringutil"
	"strings"
	"time"
)

const (
	ComparatorEquals      = "eq"
	ComparatorGreaterThan = "gt"
	ComparatorLessThan    = "lt"
	ComparatorContains    = "contains"
	ComparatorStartsWith  = "startsWith"
	ComparatorEndsWith    = "endsWith"

	VirtualFieldEventCount = "vf_event_count"
)

type FilterList []Filter

func (l FilterList) Matches(h EventHistory, e Event) bool {
	if len(l) == 0 {
		return true
	}

	for _, filter := range l {
		if filter.Matches(h, e) {
			return true
		}
	}
	return false
}

type Filter struct {
	Event string   `json:"event"`
	Rules RuleList `json:"rules"`
}

type RuleList []Rule

type Rule struct {
	Comparator string      `json:"comparator"`
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Since      string      `json:"since,omitempty"`
}

func (f Filter) Matches(h EventHistory, e Event) bool {
	if e.Type != f.Event {
		return false
	}
	if len(f.Rules) == 0 {
		return true
	}
	for _, rule := range f.Rules {
		err := populateVirtualField(h, e, rule)
		if err != nil {
			return false
		}
		v, ok := e.Values[rule.Field]
		if !ok {
			return false
		}
		switch rule.Comparator {
		case ComparatorEquals:
			if v != rule.Value {
				return false
			}
		case ComparatorContains:
			if !strings.Contains(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false
			}
		case ComparatorStartsWith:
			if !strings.HasPrefix(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false
			}
		case ComparatorEndsWith:
			if !strings.HasSuffix(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false
			}
		case ComparatorGreaterThan:
			if stringutil.Itof(v) < stringutil.Itof(rule.Value) {
				return false
			}
		case ComparatorLessThan:
			if stringutil.Itof(v) > stringutil.Itof(rule.Value) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func populateVirtualField(h EventHistory, e Event, rule Rule) error {
	if rule.Field == VirtualFieldEventCount {
		var d time.Duration
		if rule.Since != "" {
			parsed, err := time.ParseDuration(rule.Since)
			if err != nil {
				return err
			}
			d = parsed
		} else {
			d = 1 * time.Hour
		}
		events, err := h.FindWithin(e.Type, *e.CFToolsId(), d)
		if err != nil {
			return err
		}
		e.Values[VirtualFieldEventCount] = len(events)
	}
	return nil
}
