package server

import (
	"fmt"
	"sync"
)

// ILog is an interface that defines the methods for a log.
type ILog interface {
	Append(Record) (uint64, error)
	Read(uint64) (Record, error)
}

// Log is a struct that implements the ILog interface.
type Log struct {
	mu      sync.Mutex
	records []Record
}

// NewLog is a constructor for the Log struct.
func NewLog() *Log {
	return &Log{}
}

// Append is a method that appends a record to the log and returns the offset of the record.
func (l *Log) Append(record Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)
	return record.Offset, nil
}

// Read is a method that reads a record from the log at a given offset.
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if offset >= uint64(len(l.records)) {
		return Record{}, ErrOffsetNotFound
	}
	return l.records[offset], nil
}

// ErrOffsetNotFound is an error that is returned when a record is not found at a given offset.
var ErrOffsetNotFound = fmt.Errorf("offset not found")
