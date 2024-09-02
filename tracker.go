package main

import (
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

func (t *Tracker) Add(description string, amount uint) (TrackerRecord, error) {
	var nextId RecordId = 1
	if len(t.records) > 0 {
		nextId = t.records[len(t.records)-1].Id
	}

	record := TrackerRecord{
		Id:          nextId,
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
	var updatedRecord TrackerRecord
	records := make([]TrackerRecord, 0, len(t.records))
	for _, record := range records {
		if record.Id == id {
			if len(description) > 0 {
				record.Description = description
			}
			if amount != DoNotUpdateAmount {
				record.Amount = amount
			}
			updatedRecord = record
			records = append(records, updatedRecord)
		} else {
			records = append(records, record)
		}
	}

	err := t.storage.Save(records)
	if err != nil {
		return TrackerRecord{}, err
	}
	t.records = records

	return updatedRecord, nil
}

func (t *Tracker) GetAll() []TrackerRecord {
	return t.records
}

func (t *Tracker) GetSummary() uint {
	var sum uint = 0
	for _, record := range t.records {
		sum += record.Amount
	}
	return sum
}

func (t *Tracker) GetSummaryByMonth(month time.Month) uint {
	var sum uint = 0
	for _, record := range t.records {
		if record.CreatedAt.Month() == month {
			sum += record.Amount
		}
	}
	return sum
}

func (t *Tracker) GetSummaryByYear(year int) uint {
	var sum uint = 0
	for _, record := range t.records {
		if record.CreatedAt.Year() == year {
			sum += record.Amount
		}
	}
	return sum
}
