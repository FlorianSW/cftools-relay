package domain

import (
	"strconv"
)

const (
	ComparatorEquals      = "eq"
	ComparatorGreaterThan = "gt"
	ComparatorLessThan    = "lt"
)

type FilterList []Filter

func (l FilterList) Matches(e WebhookEvent) bool {
	if len(l) == 0 {
		return true
	}

	for _, filter := range l {
		if filter.Matches(e) {
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
}

func (f Filter) Matches(e WebhookEvent) bool {
	if e.Event != f.Event {
		return false
	}
	if len(f.Rules) == 0 {
		return true
	}
	for _, rule := range f.Rules {
		v, ok := e.ParsedPayload[rule.Field]
		if !ok {
			return false
		}
		switch rule.Comparator {
		case ComparatorEquals:
			if v != rule.Value {
				return false
			}
		case ComparatorGreaterThan:
			if itof(v) < itof(rule.Value) {
				return false
			}
		case ComparatorLessThan:
			if itof(v) > itof(rule.Value) {
				return false
			}
		default:
			return false
		}
	}
	return true
}

func itof(v interface{}) float32 {
	switch value := v.(type) {
	case float32:
		return value
	case int:
	case float64:
		return float32(value)
	case string:
		r, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return -1
		}
		return float32(r)
	}
	return -1
}
