package main

import (
	"log"
	"testing"
)

func _TestStart(t *testing.T) {
	m := Master{}
	got, err := m.start()
	if err != nil {
		t.Errorf("want no error, got %#v", err)
	}
	log.Println(string(got))
}
