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

	It("MatchingFilters any event when no filter defined", func() {
		filters := domain.FilterList{}

		matches, filter, _ := filters.MatchingFilters(history, someEvent)
		Expect(matches).To(BeTrue())
		Expect(filter).To(HaveLen(0))
	})

	Context("Event filter", func() {
		It("MatchingFilters event", func() {
			filters := domain.FilterList{{
				Event: someEvent.Type,
				Rules: nil,
			}}

			matches, filter, _ := filters.MatchingFilters(history, someEvent)
			Expect(matches).To(BeTrue())
			Expect(filter[0]).To(Equal(filters[0]))
		})

		It("matching multiple filters", func() {
			filters := domain.FilterList{{
				Event: someEvent.Type,
				Rules: domain.RuleList{{
					Comparator: "eq",
					Field:      "someKey",
					Value:      "someValue",
				}},
			}, {
				Event: someEvent.Type,
				Rules: domain.RuleList{{
					Comparator: "eq",
					Field:      "someKey",
					Value:      "someValue",
				}},
			}}

			matches, filter, _ := filters.MatchingFilters(history, someEvent)
			Expect(matches).To(BeTrue())
			Expect(filter).To(HaveLen(2))
		})

		It("does not match event", func() {
			filters := domain.FilterList{{
				Event: domain.EventUserJoin,
				Rules: nil,
			}}

			matches, filter, _ := filters.MatchingFilters(history, someEvent)
			Expect(matches).To(BeFalse())
			Expect(filter).To(HaveLen(0))
		})
	})

	Context("Event filter rules", func() {
		Context("EQ comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "someValue",
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("GT comparator", func() {
			It("MatchingFilters float event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("LT comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("contains comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "meVal",
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("startsWith comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("endsWith comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})
		})
		Context("oneOf comparator", func() {
			It("MatchingFilters event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "oneOf",
						Field:      "someKey",
						Value:      []string{"anotherValue", "someValue"},
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "oneOf",
						Field:      "someKey",
						Value:      []string{"anotherValue", "yetAnotherValue"},
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
			})

			It("behaves like eq when no array of strings given", func() {
				filters := domain.FilterList{{
					Event: someEvent.Type,
					Rules: domain.RuleList{{
						Comparator: "oneOf",
						Field:      "someKey",
						Value:      "someValue",
					}},
				}}

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
			})
		})
	})

	Context("virtual fields", func() {
		Context("event_count", func() {
			It("MatchingFilters when event_count is greater than", func() {
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeTrue())
				Expect(filter[0]).To(Equal(filters[0]))
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

				matches, filter, _ := filters.MatchingFilters(history, someEvent)
				Expect(matches).To(BeFalse())
				Expect(filter).To(HaveLen(0))
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
