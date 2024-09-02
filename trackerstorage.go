package main

type TrackerStorage interface {
	ReadAll() []TrackerRecord
	Save(records []TrackerRecord) error
}
