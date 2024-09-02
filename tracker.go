package main

import (
	"errors"
	"slices"
	"time"
)

type RecordId uint

type TrackerRecord struct {
	Id          RecordId
	Description string
	Amount      uint
	CreatedAt   time.Time
}

type Tracker struct {
	storage TrackerStorage
	records []TrackerRecord
}

type TrackerStorage interface {
	Save(records []TrackerRecord) error
}

func (t *Tracker) Add(description string, amount uint) (TrackerRecord, error) {
	record := TrackerRecord{
		Id:          t.records[len(t.records)-1].Id,
		Description: description,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}
	records := append(t.records, record)

	err := t.storage.Save(records)
	if err != nil {
		return TrackerRecord{}, err
	}
	t.records = records
	return record, nil
}

func (t *Tracker) Delete(id RecordId) error {
	records := slices.DeleteFunc(t.records, func(record TrackerRecord) bool {
		return record.Id == id
	})

	err := t.storage.Save(records)
	if err != nil {
		return err
	}
	t.records = records

	return nil
}

const DoNotUpdateAmount = 0

func (t *Tracker) Update(id RecordId, description string, amount uint) (TrackerRecord, error) {
	// todo
	return TrackerRecord{}, errors.New("not implemented yet")
}

func (t *Tracker) GetAll() []TrackerRecord {
	return t.records
}

func (t *Tracker) GetSummary() uint {
	// todo
	return 0
}

func (t *Tracker) GetSummaryByMonth(month time.Month) uint {
	// todo
	return 0
}

func (t *Tracker) GetSummaryByYear(year int) uint {
	// todo
	return 0
}
