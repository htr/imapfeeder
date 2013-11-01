package main

import (
	"encoding/json"
	"os"
	"time"
)

type Log map[string]time.Time

type AppendLog struct {
	log      Log
	filename string
}

func loadAppendLog(filename string) (*AppendLog, error) {
	appendLog := &AppendLog{
		log:      Log{},
		filename: filename,
	}

	file, err := os.Open(filename)
	if err != nil {
		// file doesnt exist FIXME
		return appendLog, nil
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&appendLog.log); err != nil {
		return nil, err
	} else {
		return appendLog, nil
	}

}

func (l *AppendLog) Add(key string) {
	l.log[key] = time.Now()
}

func (l *AppendLog) Exists(key string) bool {
	_, ok := l.log[key]
	return ok
}

func (l *AppendLog) Save() error {
	file, err := os.Create(l.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(&l.log); err != nil {
		return err
	}

	return nil
}
