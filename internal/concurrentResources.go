package internal

import (
	"bufio"
	"github.com/cockroachdb/pebble"
	"io"
	"os"
	"strings"
	"sync"
)

type ResourceCache struct {
	dbMutex sync.RWMutex
}

type DBCache struct {
	rc    ResourceCache
	store map[string]*pebble.DB
}

func (db *DBCache) Get(key string) *pebble.DB {
	db.rc.dbMutex.RLock()
	defer db.rc.dbMutex.RUnlock()
	return db.store[key]
}

func (db *DBCache) Set(key string, value *pebble.DB) *pebble.DB {
	db.rc.dbMutex.Lock()
	defer db.rc.dbMutex.Unlock()
	db.store[key] = value
	return value
}

type EventWriter struct {
	wg            sync.WaitGroup
	mu            sync.Mutex
	handle        *os.File
	currentOffset uint64
}

func (ew *EventWriter) Write(b []byte) {
	ew.wg.Add(1)
	go func([]byte) {
		defer ew.wg.Done()
		ew.mu.Lock()
		ew.handle.Write(b)
		ew.handle.Write([]byte("\n"))
		ew.handle.Sync()
		ew.mu.Unlock()
	}(b)
}

func (ew *EventWriter) Read(l uint64) ([]string, error) {
	reader := bufio.NewReader(ew.handle)
	ew.handle.Seek(int64(ew.currentOffset), io.SeekStart)
	var lines []string
	var offset uint64 = 0
	for i := uint64(0); i < l; i++ {
		line, err := reader.ReadString('\n')
		offset += uint64(len([]byte(line)))
		if err != nil {
			ew.currentOffset += offset
			return lines, err
		}
		lines = append(lines, strings.TrimSpace(line))
	}
	ew.currentOffset += offset
	return lines, nil
}

type EventStore struct {
	rc    ResourceCache
	store map[string]*EventWriter
}

func (es *EventStore) Get(key string) *EventWriter {
	es.rc.dbMutex.Lock()
	defer es.rc.dbMutex.Unlock()
	return es.store[key]
}

func (es *EventStore) Set(key string, value *os.File) *EventWriter {
	es.rc.dbMutex.Lock()
	defer es.rc.dbMutex.Unlock()
	var writer = &EventWriter{handle: value}
	es.store[key] = writer
	return writer
}
