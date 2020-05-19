package server

import (
	"bytes"
	"testing"

	"github.com/laiboonh/proglog/internal/server"
)

func TestAppend(t *testing.T) {
	log := server.NewLog()
	offset, err := log.Append(server.Record{[]byte("hello"), 0})
	if err != nil {
		t.Error(err)
	}
	record, err := log.Read(offset)
	if err != nil {
		t.Error(err)
	}
	if bytes.Compare(record.Value, []byte("hello")) != 0 {
		t.Error("record read different from appended")
	}
}
