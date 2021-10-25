package adapter

import (
	"cftools-relay/internal/domain"
	"encoding/json"
	"os"
	"sync"
	"time"
)

type repository struct {
	dataDir string
	lock    *sync.RWMutex
}

func NewEventRepository(dataDir string) (*repository, error) {
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		if err := os.Mkdir(dataDir, 0644); err != nil {
			return nil, err
		}
	}
	return &repository{
		dataDir: dataDir,
		lock:    &sync.RWMutex{},
	}, nil
}

func (r repository) Save(e domain.Event) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	id := e.CFToolsId()
	if id == nil {
		return domain.ErrCFToolsIdMissing
	}
	dataFile, err := ensureFile(r.dataDir + "/" + *id + ".json")
	if err != nil {
		return err
	}

	record, err := readRecords(dataFile)
	if err != nil {
		return err
	}

	record.Events = append(record.Events, e)
	c, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, c, 0655)
}

func readRecords(path string) (events, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return events{}, nil
		}
		return events{}, err
	}
	if err != nil {
		return events{}, err
	}
	var record events
	err = json.Unmarshal(c, &record)
	if err != nil {
		return events{}, err
	}
	return record, nil
}

func ensureFile(path string) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		temp := events{Events: []domain.Event{}}
		c, err := json.Marshal(temp)
		if err != nil {
			return path, err
		}
		err = os.WriteFile(path, c, 0644)
		if err != nil {
			return path, err
		}
	}
	return path, nil
}

func (r repository) FindWithin(eventType, cftoolsId string, within time.Duration) ([]domain.Event, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	dataFile := r.dataDir + "/" + cftoolsId + ".json"
	record, err := readRecords(dataFile)
	if err != nil {
		return []domain.Event{}, err
	}
	res := []domain.Event{}
	latest := time.Now().Add(-within)
	for _, event := range record.Events {
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

type events struct {
	Events []domain.Event
}
