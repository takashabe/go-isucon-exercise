package models

import (
	"database/sql"
	"reflect"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	setupFixture(t, "fixture/teams.yaml")

	cases := []struct {
		email      string
		password   string
		expectTeam *Team
		expectErr  error
	}{
		{
			"foo",
			"foo",
			&Team{ID: 1, Name: "team1", Instance: "localhost:8080"},
			nil,
		},
		{
			"",
			"",
			nil,
			sql.ErrNoRows,
		},
	}
	for i, c := range cases {
		team, err := NewTeam().Authentication(c.email, c.password)
		if err != c.expectErr {
			t.Fatalf("#%d: want %v, got %v", i, c.expectErr, err)
		}
		if !reflect.DeepEqual(c.expectTeam, team) {
			t.Fatalf("#%d: want %v, got %v", i, c.expectTeam, team)
		}
	}
}

func TestGet(t *testing.T) {
	setupFixture(t, "fixture/teams.yaml")

	cases := []struct {
		id         int
		expectTeam *Team
		expectErr  error
	}{
		{
			1,
			&Team{ID: 1, Name: "team1", Instance: "localhost:8080"},
			nil,
		},
		{
			0,
			nil,
			sql.ErrNoRows,
		},
	}
	for i, c := range cases {
		team, err := NewTeam().Get(c.id)
		if err != c.expectErr {
			t.Fatalf("#%d: want %v, got %v", i, c.expectErr, err)
		}
		if !reflect.DeepEqual(c.expectTeam, team) {
			t.Fatalf("#%d: want %v, got %v", i, c.expectTeam, team)
		}
	}
}
