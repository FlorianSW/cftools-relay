package domain_test

import (
	"cftools-relay/internal/domain"
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Filter", func() {
	var history domain.EventHistory

	someEvent := domain.Event{
		Type:      domain.EventPlayerKill,
		Timestamp: time.Now(),
		Values: map[string]interface{}{
			domain.FieldCfToolsId: "AN_ID",
			"someKey":             "someValue",
			"numberKey":           100.12,
		},
	}

	BeforeEach(func() {
		history = NewInMemoryEventHistoryRepository()
	})

	It("Matches any event when no filter defined", func() {
		filters := domain.FilterList{}

		Expect(filters.Matches(history, someEvent)).To(BeTrue())
	})

	Context("Event filter", func() {
		It("Matches event", func() {
			filters := domain.FilterList{{
				Event: someEvent.Type,
				Rules: nil,
			}}

			Expect(filters.Matches(history, someEvent)).To(BeTrue())
		})

		It("does not match event", func() {
			filters := domain.FilterList{{
				Event: domain.EventUserJoin,
				Rules: nil,
			}}

			Expect(filters.Matches(history, someEvent)).To(BeFalse())
		})
	})

	Context("Event filter rules", func() {
		Context("EQ comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "someValue",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "anotherValue",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})

			It("all rules must match to match", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "someValue",
					}, {
						Comparator: "eq",
						Field:      "someKey",
						Value:      "anotherValue",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
		Context("GT comparator", func() {
			It("Matches float event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event (JSON Number)", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      json.Number("100.13"),
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
		Context("LT comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
		Context("contains comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "meVal",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "noVal",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
		Context("startsWith comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
		Context("endsWith comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
	})

	Context("virtual fields", func() {
		Context("event_count", func() {
			It("Matches when event_count is greater than", func() {
				err := history.Save(domain.Event{
					Type:      someEvent.Type,
					Timestamp: time.Now().Add(-40 * time.Minute),
					Values: map[string]interface{}{
						domain.FieldCfToolsId: *someEvent.CFToolsId(),
					},
				})
				Expect(err).ToNot(HaveOccurred())
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      domain.VirtualFieldEventCount,
						Value:      1,
						Since:      "1h",
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeTrue())
			})
			It("does not match when less than events", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      domain.VirtualFieldEventCount,
						Value:      15,
					}},
				}}

				Expect(filters.Matches(history, someEvent)).To(BeFalse())
			})
		})
	})
})

type inMemoryRepository struct {
	data map[string][]domain.Event
}

func NewInMemoryEventHistoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		data: map[string][]domain.Event{},
	}
}

func (r inMemoryRepository) Save(e domain.Event) error {
	r.data[*e.CFToolsId()] = append(r.data[*e.CFToolsId()], e)
	return nil
}

func (r inMemoryRepository) FindWithin(eventType, cftoolsId string, within time.Duration) ([]domain.Event, error) {
	d, ok := r.data[cftoolsId]
	if !ok {
		return []domain.Event{}, nil
	}
	res := []domain.Event{}
	latest := time.Now().Add(-within)
	for _, event := range d {
		if event.Type != eventType {
			continue
		}
		if event.Timestamp.Before(latest) {
			continue
		}
		res = append(res, event)
	}
	return res, nil
}
