package adapter_test

import (
	"cftools-relay/internal/adapter"
	"cftools-relay/internal/domain"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"time"
)

var _ = Describe("EventHistory", func() {
	var (
		tmpPath string
		r       domain.EventHistory
	)

	BeforeEach(func() {
		path, err := os.MkdirTemp("", "test-data")
		if err != nil {
			panic(err)
		}
		tmpPath = path
		repo, err := adapter.NewEventRepository(path)
		if err != nil {
			panic(err)
		}
		r = repo
	})

	AfterEach(func() {
		err := os.RemoveAll(tmpPath)
		if err != nil {
			panic(err)
		}
	})

	It("persists event", func() {
		e := makeEventWithType(domain.EventUserJoin)

		err := r.Save(e)
		Expect(err).ToNot(HaveOccurred())

		saved, err := r.FindWithin(e.Type, *e.CFToolsId(), 1*time.Hour)
		Expect(err).ToNot(HaveOccurred())
		Expect(saved).To(HaveLen(1))
		Expect(saved[0].Type).To(Equal(e.Type))
		Expect(saved[0].Timestamp).To(BeTemporally("==", e.Timestamp))
		Expect(saved[0].Values).To(Equal(e.Values))
	})

	It("filters events by type", func() {
		e := mustSave(r, makeEventWithType(domain.EventUserJoin))
		mustSave(r, makeEventWithType(domain.EventPlayerDamage))
		e3 := mustSave(r, makeEventWithType(domain.EventUserJoin))

		events, err := r.FindWithin(domain.EventUserJoin, *e.CFToolsId(), 1*time.Hour)

		Expect(err).ToNot(HaveOccurred())
		Expect(events).To(HaveLen(2))
		Expect(events[0].Type).To(Equal(domain.EventUserJoin))
		Expect(events[0].Timestamp).To(BeTemporally("==", e.Timestamp))
		Expect(events[1].Type).To(Equal(domain.EventUserJoin))
		Expect(events[1].Timestamp).To(BeTemporally("==", e3.Timestamp))
	})

	It("does not save more than 100 events per id", func() {
		var e domain.Event
		for i := 0; i < 101; i++ {
			e = mustSave(r, makeEventWithType(domain.EventUserJoin))
		}

		events, err := r.FindWithin(domain.EventUserJoin, *e.CFToolsId(), 1*time.Hour)

		Expect(err).ToNot(HaveOccurred())
		Expect(events).To(HaveLen(100))
		Expect(events[len(events)-1].Timestamp).To(BeTemporally("==", e.Timestamp))
	})

	It("returns events for specified timeframe only", func() {
		e := mustSave(r, makeEventWithType(domain.EventUserJoin, time.Now().Add(-58*time.Minute)))
		mustSave(r, makeEventWithType(domain.EventUserJoin, time.Now().Add(-61*time.Minute)))
		e3 := mustSave(r, makeEventWithType(domain.EventUserJoin, time.Now().Add(-59*time.Minute)))

		events, err := r.FindWithin(domain.EventUserJoin, *e.CFToolsId(), 1*time.Hour)

		Expect(err).ToNot(HaveOccurred())
		Expect(events).To(HaveLen(2))
		Expect(events[0].Type).To(Equal(domain.EventUserJoin))
		Expect(events[0].Timestamp).To(BeTemporally("==", e.Timestamp))
		Expect(events[1].Type).To(Equal(domain.EventUserJoin))
		Expect(events[1].Timestamp).To(BeTemporally("==", e3.Timestamp))
	})
})

func makeEventWithType(t string, ts ...time.Time) domain.Event {
	var timestamp time.Time
	if len(ts) != 0 {
		timestamp = ts[0]
	} else {
		timestamp = time.Now()
	}
	return domain.Event{
		Type:      t,
		Timestamp: timestamp,
		Values: map[string]interface{}{
			domain.FieldCfToolsId: "AN_ID",
		},
	}
}

func mustSave(r domain.EventHistory, e domain.Event) domain.Event {
	err := r.Save(e)
	Expect(err).ToNot(HaveOccurred())

	return e
}
