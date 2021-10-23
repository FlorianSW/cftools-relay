package domain_test

import (
	"cftools-relay/internal/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	someEvent := domain.WebhookEvent{
		ShardId:   0,
		Flavor:    domain.FlavorCftools,
		Event:     domain.EventPlayerKill,
		Id:        "2b00830c-bbda-4387-8d8e-917dc5591177",
		Signature: "SOME_SIGNATURE",
		Payload:   "{\"someKey\": \"someValue\"}",
		ParsedPayload: map[string]interface{}{
			"someKey": "someValue",
			"numberKey": 100.12,
		},
	}

	It("Matches any event when no filter defined", func() {
		filters := domain.FilterList{}

		Expect(filters.Matches(someEvent)).To(BeTrue())
	})

	Context("Event filter", func() {
		It("Matches event", func() {
			filters := domain.FilterList{{
				Event: someEvent.Event,
				Rules: nil,
			}}

			Expect(filters.Matches(someEvent)).To(BeTrue())
		})

		It("does not match event", func() {
			filters := domain.FilterList{{
				Event: domain.EventUserJoin,
				Rules: nil,
			}}

			Expect(filters.Matches(someEvent)).To(BeFalse())
		})
	})

	Context("Event filter rules", func() {
		Context("EQ comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "someValue",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "eq",
						Field:      "someKey",
						Value:      "anotherValue",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})

			It("all rules must match to match", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
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

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
		Context("GT comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "gt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
		Context("LT comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.13,
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "lt",
						Field:      "numberKey",
						Value:      100.11,
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
		Context("contains comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "meVal",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "contains",
						Field:      "someKey",
						Value:      "noVal",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
		Context("startsWith comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "startsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
		Context("endsWith comparator", func() {
			It("Matches event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "Value",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeTrue())
			})

			It("does not match event", func() {
				filters := domain.FilterList{{
					Event: someEvent.Event,
					Rules: domain.RuleList{{
						Comparator: "endsWith",
						Field:      "someKey",
						Value:      "some",
					}},
				}}

				Expect(filters.Matches(someEvent)).To(BeFalse())
			})
		})
	})
})
