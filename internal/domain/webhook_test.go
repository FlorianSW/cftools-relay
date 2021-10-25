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