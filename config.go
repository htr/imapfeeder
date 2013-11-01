package main

import (
	"encoding/json"
	"os"
)

type Context struct {
	Folders      map[string][]string
	ImapServer   string
	Labels       []string
	Password     string
	From         string
	To           string
	Username     string
	Template     string
	FolderPrefix string
	Cleanup      bool
	Jobs         int
	Imap         *ImapSession `json:"-"`
}

func loadConfig(filename string) (*Context, error) {
	ctx := &Context{}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(ctx)

	if err != nil {
		return nil, err
	}

	return ctx, nil
}
