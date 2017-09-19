package models

import (
	"database/sql"
	"reflect"
	"sort"
	"testing"
)

func TestScoreGet(t *testing.T) {
	setupFixture(t, "fixture/scores.yaml")

	cases := []struct {
		scoreID   int
		teamID    int
		expectErr error
	}{
		{1, 1, nil},
		{1, 2, sql.ErrNoRows},
	}
	for i, c := range cases {
		_, err := NewScore().Get(c.scoreID, c.teamID)
		if err != c.expectErr {
			t.Errorf("#%d: want error %v, got %v", i, c.expectErr, err)
		}
	}
}

func TestHistory(t *testing.T) {
	setupFixture(t, "fixture/scores.yaml")

	cases := []struct {
		teamID    int
		expectIDs []int
	}{
		{1, []int{1, 2, 4}},
		{2, []int{3}},
	}
	for i, c := range cases {
		scores, err := NewScore().History(c.teamID)
		if err != nil {
			t.Fatalf("#%d: want non error, got %v", i, err)
		}
		ids := []int{}
		for _, s := range scores {
			ids = append(ids, s.ID)
		}
		sort.Ints(ids)
		if !reflect.DeepEqual(c.expectIDs, ids) {
			t.Errorf("#%d: want %v, got %v", i, c.expectIDs, ids)
		}
	}
}
