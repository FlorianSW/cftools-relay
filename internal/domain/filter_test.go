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

	It("MatchingFilter any event when no filter defined", func() {
		filters := domain.FilterList{}

		matches, filter, _ := filters.MatchingFilter(history, someEvent)
		Expect(matches).To(BeTrue())
		Expect(filter).To(BeNil())
	})

	Context("Event filter", func() {
		It("MatchingFilter event", func() {
			filters := domain.FilterList{{
				Event: someEvent.Type,
				Rules: nil,
			}}

			matches, filter, _ := filters.MatchingFilter(history, someEvent)
			Expect(matches).To(BeTrue())
			Expect(*filter).To(Equal(filters[0]))
		})

		It("does not match event", func() {
			filters := domain.FilterList{{
				Event: domain.EventUserJoin,
				Rules: nil,
			}}

			matches, filter, _ := filters.MatchingFilter(history, someEvent)
			Expect(matches).To(BeFalse())
			Expect(filter).To(BeNil())
		})
	})

	Context("Event filter rules", func() {
		Context("EQ comparator", func() {
			It("MatchingFilter event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "someValue",
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
		Context("GT comparator", func() {
			It("MatchingFilter float event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
		Context("LT comparator", func() {
			It("MatchingFilter event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
		Context("contains comparator", func() {
			It("MatchingFilter event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "meVal",
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
		Context("startsWith comparator", func() {
			It("MatchingFilter event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
		Context("endsWith comparator", func() {
			It("MatchingFilter event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
			})
		})
	})

	Context("virtual fields", func() {
		Context("event_count", func() {
			It("MatchingFilter when event_count is greater than", func() {
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(*filter).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilter(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(BeNil())
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
