package benchmark

import (
	"log"
	"testing"
)

// TODO: integration test
func TestStart(t *testing.T) {
	master, err := NewMaster("localhost", 8080, "data/param.json", "test")
	if err != nil {
		t.Fatalf("want no error, got %#v", err)
	}

	got, err := master.start()
	if err != nil {
		t.Fatalf("want no error, got %#v", err)
	}
	log.Println(string(got))
}
