package main

type TrackerStorage interface {
	ReadAll() ([]TrackerRecord, error)
	Save(records []TrackerRecord) error
}
