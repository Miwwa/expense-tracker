package main

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type FakeStorage struct {
	readError error
	saveError error
	records   []TrackerRecord
}

func (f *FakeStorage) ReadAll() ([]TrackerRecord, error) {
	if f.readError != nil {
		return nil, f.readError
	}
	return f.records, nil
}

func (f *FakeStorage) Save(records []TrackerRecord) error {
	f.records = records
	return f.saveError
}

func TestNewTracker(t *testing.T) {
	tests := []struct {
		name        string
		storageErr  error
		storageData []TrackerRecord
		expectedErr error
	}{
		{
			name: "SuccessWithNoDataInStorage",
		},
		{
			name: "SuccessWithDataInStorage",
			storageData: []TrackerRecord{
				{
					Id:          0,
					Description: "TestDescription",
					Amount:      50,
					CreatedAt:   time.Now(),
				},
			},
		},
		{
			name:        "FailureDueToStorageError",
			storageErr:  errors.New("test storage error"),
			expectedErr: errors.New("test storage error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := &FakeStorage{readError: test.storageErr, records: test.storageData}

			tracker, err := NewTracker(storage)

			if err != nil {
				if err.Error() != test.expectedErr.Error() {
					t.Errorf("Got error %v, expected %v", err, test.expectedErr)
				}
			} else if !reflect.DeepEqual(tracker.records, test.storageData) {
				t.Errorf("Got tracker data %v, expected %v", tracker.records, test.storageData)
			}
		})
	}
}

func TestTrackerAdd(t *testing.T) {
	tests := []struct {
		name        string
		storageErr  error
		description string
		amount      uint
		expectedErr error
	}{
		{
			name:        "SuccessNoPreviousRecords",
			description: "TestDescription1",
			amount:      50,
		},
		{
			name:        "SuccessWithPreviousRecords",
			description: "TestDescription2",
			amount:      100,
		},
		{
			name:        "StorageError",
			description: "TestDescription3",
			amount:      150,
			storageErr:  errors.New("test storage error"),
			expectedErr: errors.New("test storage error"),
		},
		{
			name:        "EmptyDescription",
			expectedErr: errors.New("description cannot be empty"),
		},
		{
			name:        "ZeroAmount",
			description: "TestDescription5",
			amount:      0,
			expectedErr: errors.New("amount cannot be zero"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := &FakeStorage{saveError: test.storageErr}
			tracker, _ := NewTracker(storage)

			// if previous records exist, add one
			if test.name == "SuccessWithPreviousRecords" {
				existingRecord := TrackerRecord{
					Id:          0,
					Description: "ExistingDescription",
					Amount:      50,
					CreatedAt:   time.Now(),
				}
				tracker.records = append(tracker.records, existingRecord)
			}

			record, err := tracker.Add(test.description, test.amount)

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("Got error %v, expected %v", err, test.expectedErr)
			}

			if err == nil {
				if record.Description != test.description || record.Amount != test.amount {
					t.Errorf("Got record %v, expected description '%s' and amount %d", record, test.description, test.amount)
				}
			}
		})
	}
}

func TestTrackerDelete(t *testing.T) {
	tests := []struct {
		name        string
		storageErr  error
		id          RecordId
		setupData   []TrackerRecord
		expectedErr error
		expectedRes []TrackerRecord
	}{
		{
			name:        "Success",
			id:          1,
			setupData:   []TrackerRecord{{Id: 1}, {Id: 2}},
			expectedRes: []TrackerRecord{{Id: 2}},
		},
		{
			name:        "RecordNotFound",
			id:          3,
			setupData:   []TrackerRecord{{Id: 1}, {Id: 2}},
			expectedRes: []TrackerRecord{{Id: 1}, {Id: 2}},
		},
		{
			name:        "StorageError",
			id:          1,
			setupData:   []TrackerRecord{{Id: 1}, {Id: 2}},
			storageErr:  errors.New("test storage error"),
			expectedErr: errors.New("test storage error"),
			expectedRes: []TrackerRecord{{Id: 1}, {Id: 2}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := &FakeStorage{saveError: test.storageErr}
			tracker, _ := NewTracker(storage)

			// if setup data exist, add them
			if len(test.setupData) > 0 {
				tracker.records = append(tracker.records, test.setupData...)
			}

			err := tracker.Delete(test.id)

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("Got error %v, expected %v", err, test.expectedErr)
			}

			if err == nil {
				if !reflect.DeepEqual(tracker.records, test.expectedRes) {
					t.Errorf("Got tracker data %v, expected %v", tracker.records, test.expectedRes)
				}
			}
		})
	}
}

func TestTrackerUpdate(t *testing.T) {
	tests := []struct {
		name            string
		storageErr      error
		id              RecordId
		setupData       []TrackerRecord
		updateDesc      string
		updateAmount    uint
		expectedErr     error
		expectedUpdated TrackerRecord
		expectedRes     []TrackerRecord
	}{
		{
			name:            "SuccessUpdateDescription",
			id:              1,
			updateDesc:      "UpdatedDescription",
			setupData:       []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
			expectedUpdated: TrackerRecord{Id: 1, Description: "UpdatedDescription", Amount: 100},
			expectedRes:     []TrackerRecord{{Id: 1, Description: "UpdatedDescription", Amount: 100}, {Id: 2}},
		},
		{
			name:            "SuccessUpdateAmount",
			id:              1,
			updateAmount:    200,
			setupData:       []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
			expectedUpdated: TrackerRecord{Id: 1, Description: "InitialDescription", Amount: 200},
			expectedRes:     []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 200}, {Id: 2}},
		},
		{
			name:            "SuccessUpdateBoth",
			id:              1,
			updateDesc:      "UpdatedDescription",
			updateAmount:    200,
			setupData:       []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
			expectedUpdated: TrackerRecord{Id: 1, Description: "UpdatedDescription", Amount: 200},
			expectedRes:     []TrackerRecord{{Id: 1, Description: "UpdatedDescription", Amount: 200}, {Id: 2}},
		},
		{
			name:        "RecordNotFound",
			id:          3,
			updateDesc:  "UpdatedDescription",
			setupData:   []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
			expectedErr: errors.New("record not found"),
			expectedRes: []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
		},
		{
			name:        "StorageError",
			id:          1,
			storageErr:  errors.New("test storage error"),
			setupData:   []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
			expectedErr: errors.New("test storage error"),
			expectedRes: []TrackerRecord{{Id: 1, Description: "InitialDescription", Amount: 100}, {Id: 2}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := &FakeStorage{saveError: test.storageErr}
			tracker, _ := NewTracker(storage)

			// if setup data exist, add them
			if len(test.setupData) > 0 {
				tracker.records = append(tracker.records, test.setupData...)
			}

			updatedRecord, err := tracker.Update(test.id, test.updateDesc, test.updateAmount)

			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("Got error %v, expected %v", err, test.expectedErr)
			}

			if err == nil {
				if !reflect.DeepEqual(updatedRecord, test.expectedUpdated) {
					t.Errorf("Got updated record %v, expected %v", updatedRecord, test.expectedUpdated)
				}
				if !reflect.DeepEqual(tracker.records, test.expectedRes) {
					t.Errorf("Got tracker data %v, expected %v", tracker.records, test.expectedRes)
				}
			}
		})
	}
}

func TestTrackerGetAll(t *testing.T) {
	tests := []struct {
		name      string
		setupData []TrackerRecord
		expected  []TrackerRecord
	}{
		{
			name:     "NoData",
			expected: []TrackerRecord{},
		},
		{
			name:      "WithData",
			setupData: []TrackerRecord{{Id: 1}, {Id: 2}},
			expected:  []TrackerRecord{{Id: 1}, {Id: 2}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			storage := &FakeStorage{records: []TrackerRecord{}}
			tracker, _ := NewTracker(storage)

			// if setup data exist, add them
			if len(test.setupData) > 0 {
				tracker.records = append(tracker.records, test.setupData...)
			}

			result := tracker.GetAll()

			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Got tracker data %v, expected %v", result, test.expected)
			}
		})
	}
}

func TestTrackerGetSummary(t *testing.T) {
	tests := []struct {
		name string
		data []TrackerRecord
		want uint
	}{
		{name: "NoData", want: 0},
		{
			name: "SingleRecord",
			data: []TrackerRecord{{Amount: 100}},
			want: 100,
		},
		{
			name: "MultipleRecords",
			data: []TrackerRecord{{Amount: 100}, {Amount: 200}, {Amount: 300}},
			want: 600,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &FakeStorage{records: tt.data}
			tracker, _ := NewTracker(storage)
			if got := tracker.GetSummary(); got != tt.want {
				t.Errorf("Tracker.GetSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTrackerGetSummaryByMonth(t *testing.T) {
	tests := []struct {
		name  string
		data  []TrackerRecord
		month time.Month
		want  uint
	}{
		{name: "NoData", month: time.January, want: 0},
		{
			name:  "SingleRecord",
			month: time.January,
			data:  []TrackerRecord{{Amount: 100, CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)}},
			want:  100,
		},
		{
			name:  "MultipleRecords",
			month: time.January,
			data: []TrackerRecord{
				{Amount: 100, CreatedAt: time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)},
				{Amount: 200, CreatedAt: time.Date(2024, time.January, 2, 0, 0, 0, 0, time.UTC)},
				{Amount: 300, CreatedAt: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)}},
			want: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &FakeStorage{records: tt.data}
			tracker, _ := NewTracker(storage)
			if got := tracker.GetSummaryByMonth(tt.month); got != tt.want {
				t.Errorf("Tracker.GetSummaryByMonth() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestTrackerGetSummaryByYear(t *testing.T) {
	tests := []struct {
		name string
		year int
		data []TrackerRecord
		want uint
	}{
		{
			name: "NoData",
			year: 2024,
			want: 0,
		},
		{
			name: "SingleRecord",
			year: 2023,
			data: []TrackerRecord{
				{Amount: 100, CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want: 100,
		},
		{
			name: "MultipleRecordsSameYear",
			year: 2022,
			data: []TrackerRecord{
				{Amount: 100, CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Amount: 200, CreatedAt: time.Date(2022, 2, 1, 0, 0, 0, 0, time.UTC)},
				{Amount: 300, CreatedAt: time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC)},
			},
			want: 600,
		},
		{
			name: "MultipleRecordsDifferentYears",
			year: 2021,
			data: []TrackerRecord{
				{Amount: 100, CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)},
				{Amount: 200, CreatedAt: time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)},
				{Amount: 300, CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &FakeStorage{records: tt.data}
			tracker, _ := NewTracker(storage)
			if got := tracker.GetSummaryByYear(tt.year); got != tt.want {
				t.Errorf("Tracker.GetSummaryByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
