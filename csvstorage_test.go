package main

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestCsvTrackerStorage_ReadAll(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []TrackerRecord
		wantErr bool
	}{
		{
			name:    "EmptyFile",
			content: "",
			want:    []TrackerRecord{},
			wantErr: false,
		},
		{
			name:    "HeaderOnly",
			content: "Id,CreatedAt,Amount,Description\n",
			want:    []TrackerRecord{},
			wantErr: false,
		},
		{
			name: "SingleRecord",
			content: "Id,CreatedAt,Amount,Description\n" +
				"1,2024-01-01T01:01:01Z,100,record1\n",
			want: []TrackerRecord{
				{
					Id:          1,
					CreatedAt:   time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
					Amount:      100,
					Description: "record1",
				},
			},
			wantErr: false,
		},
		{
			name: "MultipleRecords",
			content: "Id,CreatedAt,Amount,Description\n" +
				"1,2024-01-01T01:01:01Z,100,record1\n" +
				"2,2024-01-02T02:02:02Z,200,record2\n",
			want: []TrackerRecord{
				{
					Id:          1,
					CreatedAt:   time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
					Amount:      100,
					Description: "record1",
				},
				{
					Id:          2,
					CreatedAt:   time.Date(2024, 1, 2, 2, 2, 2, 0, time.UTC),
					Amount:      200,
					Description: "record2",
				},
			},
			wantErr: false,
		},
		{
			name:    "InvalidFormat",
			content: "Invalid content",
			want:    nil,
			wantErr: true,
		},
		{
			name: "InvalidId",
			content: "Id,CreatedAt,Amount,Description\n" +
				"-1,2024-01-01T01:01:01Z,100,record1\n",
			want:    nil,
			wantErr: true,
		},
		{
			name: "InvalidDate",
			content: "Id,CreatedAt,Amount,Description\n" +
				"1,2024-01-01111T01:01:01Z,100,record1\n",
			want:    nil,
			wantErr: true,
		},
		{
			name: "InvalidAmount",
			content: "Id,CreatedAt,Amount,Description\n" +
				"1,2024-01-01T01:01:01Z,-125,record1\n",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &CsvTrackerStorage{
				filename: "trackerstorage_test.csv",
			}
			if err := os.WriteFile(s.filename, []byte(tt.content), 0666); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
			defer os.Remove(s.filename)

			got, err := s.ReadAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("CsvTrackerStorage.ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CsvTrackerStorage.ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCsvTrackerStorage_Save(t *testing.T) {
	tests := []struct {
		name     string
		records  []TrackerRecord
		expected string
		wantErr  bool
	}{
		{
			name:     "EmptyRecords",
			records:  []TrackerRecord{},
			wantErr:  false,
			expected: "Id,CreatedAt,Amount,Description\n",
		},
		{
			name: "SingleRecord",
			records: []TrackerRecord{
				{
					Id:          1,
					CreatedAt:   time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
					Amount:      100,
					Description: "record1",
				},
			},
			expected: "Id,CreatedAt,Amount,Description\n1,2024-01-01T01:01:01Z,100,record1\n",
			wantErr:  false,
		},
		{
			name: "MultipleRecords",
			records: []TrackerRecord{
				{
					Id:          1,
					CreatedAt:   time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
					Amount:      100,
					Description: "record1",
				},
				{
					Id:          2,
					CreatedAt:   time.Date(2024, 1, 2, 2, 2, 2, 0, time.UTC),
					Amount:      200,
					Description: "record2",
				},
			},
			expected: "Id,CreatedAt,Amount,Description\n1,2024-01-01T01:01:01Z,100,record1\n2,2024-01-02T02:02:02Z,200,record2\n",
			wantErr:  false,
		},
		{
			name: "CommaInDescription",
			records: []TrackerRecord{
				{
					Id:          1,
					CreatedAt:   time.Date(2024, 1, 1, 1, 1, 1, 0, time.UTC),
					Amount:      100,
					Description: "long, lorem ipsum",
				},
			},
			expected: "Id,CreatedAt,Amount,Description\n1,2024-01-01T01:01:01Z,100,\"long, lorem ipsum\"\n",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStorageFromFile("trackerstorage_test.csv")
			// Clean up the file after each test.
			defer os.Remove(s.filename)

			if err := s.Save(tt.records); (err != nil) != tt.wantErr {
				t.Errorf("CsvTrackerStorage.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			bytes, err := os.ReadFile(s.filename)
			if err != nil {
				t.Errorf("CsvTrackerStorage.Save() file read error = %v", err)
			}
			str := string(bytes)
			if str != tt.expected {
				t.Errorf("CsvTrackerStorage.Save() = %v, want %v", str, tt.expected)
			}
		})
	}
}
