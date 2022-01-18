package domain

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Event", func() {
	Context("CFToolsId", func() {
		It("returns primary player CFTools, if present", func() {
			e := Event{
				Type:      EventUserJoin,
				Timestamp: time.Now(),
				Values: map[string]interface{}{
					FieldCfToolsId:         "AN_ID",
					FieldMurdererCfToolsId: "ANOTHER_ID",
				},
			}

			Expect(*e.CFToolsId()).To(Equal("AN_ID"))
		})

		It("returns murder CFTools, if present", func() {
			e := Event{
				Type:      EventUserJoin,
				Timestamp: time.Now(),
				Values: map[string]interface{}{
					FieldVictimCfToolsId:   "ANOTHER_ID",
					FieldMurdererCfToolsId: "AN_ID",
				},
			}

			Expect(*e.CFToolsId()).To(Equal("AN_ID"))
		})

		It("returns player ID", func() {
			e := Event{
				Type:      EventPlayerPlace,
				Timestamp: time.Now(),
				Values: map[string]interface{}{
					FieldPlayerId: "AN_ID",
				},
			}

			Expect(*e.CFToolsId()).To(Equal("AN_ID"))
		})

		Context("victim_id", func() {
			It("uses victim_id when no other values present", func() {
				e := Event{
					Type:      EventPlayerDeathStarvation,
					Timestamp: time.Now(),
					Values: map[string]interface{}{
						FieldVictimCfToolsId: "AN_ID",
					},
				}

				Expect(*e.CFToolsId()).To(Equal("AN_ID"))
			})

			It("ignores victim_id when murderer_id is present", func() {
				e := Event{
					Type:      EventPlayerKill,
					Timestamp: time.Now(),
					Values: map[string]interface{}{
						FieldMurdererCfToolsId: "AN_ID",
						FieldVictimCfToolsId:   "ANOTHER_ID",
					},
				}

				Expect(*e.CFToolsId()).To(Equal("AN_ID"))
			})
		})

		It("returns nil if no id present", func() {
			e := Event{
				Type:      EventUserJoin,
				Timestamp: time.Now(),
				Values:    map[string]interface{}{},
			}

			Expect(e.CFToolsId()).To(BeNil())
		})
	})
})
