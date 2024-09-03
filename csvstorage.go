package main

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	invalidCsvLine = errors.New("invalid csv line")
)

type CsvTrackerStorage struct {
	filename string
}

func NewStorageFromFile(filename string) *CsvTrackerStorage {
	return &CsvTrackerStorage{filename: filename}
}

func (s *CsvTrackerStorage) ReadAll() ([]TrackerRecord, error) {
	var records = make([]TrackerRecord, 0)
	file, err := os.OpenFile(s.filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return records, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 4

	for {
		parts, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return records, nil
			}
			return nil, errors.Join(invalidCsvLine, err)
		}
		if parts[0] == "Id" {
			continue
		}

		record, err := fromCsv(parts)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
}

func (s *CsvTrackerStorage) Save(records []TrackerRecord) error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	writer := csv.NewWriter(file)

	defer func() {
		writer.Flush()
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()

	headers := []string{"Id", "CreatedAt", "Amount", "Description"}
	err = writer.Write(headers)
	if err != nil {
		return err
	}
	for _, record := range records {
		err := writer.Write(toCsv(record))
		if err != nil {
			return err
		}
	}
	return nil
}

func fromCsv(parts []string) (TrackerRecord, error) {
	if len(parts) != 4 {
		return TrackerRecord{}, invalidCsvLine
	}
	id, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return TrackerRecord{}, errors.Join(invalidCsvLine, err)
	}

	createdAt, err := time.Parse(time.RFC3339, parts[1])
	if err != nil {
		return TrackerRecord{}, errors.Join(invalidCsvLine, err)
	}

	amount, err := strconv.ParseUint(parts[2], 10, 32)
	if err != nil {
		return TrackerRecord{}, errors.Join(invalidCsvLine, err)
	}

	return TrackerRecord{
		Id:          RecordId(id),
		Description: parts[3],
		Amount:      uint(amount),
		CreatedAt:   createdAt,
	}, nil
}

func toCsv(record TrackerRecord) []string {
	return []string{
		strconv.FormatUint(uint64(record.Id), 10),
		record.CreatedAt.Format(time.RFC3339),
		strconv.FormatUint(uint64(record.Amount), 10),
		record.Description,
	}
}
