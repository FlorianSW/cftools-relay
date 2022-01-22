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
	ComparatorOneOf       = "oneOf"

	VirtualFieldEventCount = "vf_event_count"

	ColorAqua            = 1752220
	ColorDarkAqua        = 1146986
	ColorGreen           = 3066993
	ColorDarkGreen       = 2067276
	ColorBlue            = 3447003
	ColorDarkBlue        = 2123412
	ColorPurple          = 10181046
	ColorDarkPurple      = 7419530
	ColorGold            = 15844367
	ColorDarkGold        = 12745742
	ColorOrange          = 15105570
	ColorDarkOrange      = 11027200
	ColorRed             = 15158332
	ColorDarkRed         = 10038562
	ColorGrey            = 9807270
	ColorDarkGrey        = 9936031
	ColorDarkerGrey      = 8359053
	ColorLightGrey       = 12370112
	ColorNavy            = 3426654
	ColorDarkNavy        = 2899536
	ColorYellow          = 16776960
	ColorWhite           = 16777215
	ColorBlurple         = 5793266
	ColorGreyple         = 10070709
	ColorDarkButNotBlack = 2895667
	ColorNotQuiteBlack   = 2303786
	ColorFuschia         = 15418782
	ColorBlack           = 2303786
)

var colorMapping = map[string]int{
	"AQUA":               ColorAqua,
	"GREEN":              ColorGreen,
	"DARK_GREEN":         ColorDarkGreen,
	"BLUE":               ColorBlue,
	"DARK_BLUE":          ColorDarkBlue,
	"PURPLE":             ColorPurple,
	"DARK_PURPLE":        ColorDarkPurple,
	"GOLD":               ColorGold,
	"DARK_GOLD":          ColorDarkGold,
	"ORANGE":             ColorOrange,
	"DARK_ORANGE":        ColorDarkOrange,
	"RED":                ColorRed,
	"DARK_RED":           ColorDarkRed,
	"GREY":               ColorGrey,
	"DARK_GREY":          ColorDarkGrey,
	"DARKER_GREY":        ColorDarkerGrey,
	"LIGHT_GREY":         ColorLightGrey,
	"NAVY":               ColorNavy,
	"DARK_NAVY":          ColorDarkNavy,
	"YELLOW":             ColorYellow,
	"WHITE":              ColorWhite,
	"BLURPLE":            ColorBlurple,
	"GREYPLE":            ColorGreyple,
	"DARK_BUT_NOT_BLACK": ColorDarkButNotBlack,
	"NOT_QUITE_BLACK":    ColorNotQuiteBlack,
	"FUSCHIA":            ColorFuschia,
	"BLACK":              ColorBlack,
}

type FilterList []Filter

func (l FilterList) MatchingFilter(h EventHistory, e Event) (bool, *Filter, error) {
	if len(l) == 0 {
		return true, nil, nil
	}

	for _, filter := range l {
		m, err := filter.Matches(h, e)
		if err != nil {
			return false, nil, err
		}
		if m {
			return true, &filter, nil
		}
	}
	return false, nil, nil
}

type Color string

func (c Color) Int() int {
	v, ok := colorMapping[string(c)]
	if ok {
		return v
	}
	return ColorDarkBlue
}

type Filter struct {
	Event   string   `json:"event"`
	Rules   RuleList `json:"rules"`
	Message string   `json:"message,omitempty"`
	Color   Color    `json:"color,omitempty"`
}

type RuleList []Rule

type Rule struct {
	Comparator string      `json:"comparator"`
	Field      string      `json:"field"`
	Value      interface{} `json:"value"`
	Since      string      `json:"since,omitempty"`
}

func (f Filter) Matches(h EventHistory, e Event) (bool, error) {
	if e.Type != f.Event {
		return false, nil
	}
	if len(f.Rules) == 0 {
		return true, nil
	}
	for _, rule := range f.Rules {
		err := populateVirtualField(h, e, rule)
		if err != nil {
			return false, err
		}
		v, ok := e.Values[rule.Field]
		if !ok {
			return false, nil
		}
		switch rule.Comparator {
		case ComparatorEquals:
			if v != rule.Value {
				return false, nil
			}
		case ComparatorContains:
			if !strings.Contains(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false, nil
			}
		case ComparatorStartsWith:
			if !strings.HasPrefix(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false, nil
			}
		case ComparatorEndsWith:
			if !strings.HasSuffix(stringutil.Itos(v), stringutil.Itos(rule.Value)) {
				return false, nil
			}
		case ComparatorGreaterThan:
			if stringutil.Itof(v) < stringutil.Itof(rule.Value) {
				return false, nil
			}
		case ComparatorLessThan:
			if stringutil.Itof(v) > stringutil.Itof(rule.Value) {
				return false, nil
			}
		case ComparatorOneOf:
			if !containsValue(v, rule.Value) {
				return false, nil
			}
		default:
			return false, nil
		}
	}
	return true, nil
}

func containsValue(v interface{}, values interface{}) bool {
	switch x := values.(type) {
	case []string:
		for _, i := range x {
			if v == i {
				return true
			}
		}
	case string:
		return v == x
	}
	return false
}

func populateVirtualField(h EventHistory, e Event, rule Rule) error {
	if _, ok := e.Values[VirtualFieldEventCount]; rule.Field == VirtualFieldEventCount && !ok {
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
